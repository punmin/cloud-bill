package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ucloud/ucloud-sdk-go/services/ubill"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

func GetUCloudBill(month string, account CloudAccount) ([]ubill.BillDetailItem, error) {
	cfg := ucloud.NewConfig()
	cfg.Region = "cn-gd"

	credential := auth.NewCredential()
	credential.PublicKey = account.AccessKeyID
	credential.PrivateKey = account.AccessKeySecret

	ubillClient := ubill.NewClient(&cfg, &credential)

	req := ubillClient.NewListUBillDetailRequest()
	// 设置请求参数
	req.BillingCycle = ucloud.String(month)
	req.Limit = ucloud.Int(100)

	offset := ucloud.Int(0)
	total := 0

	resourceSummarySet := []ubill.BillDetailItem{}

	for {
		req.Offset = offset
		response, err := ubillClient.ListUBillDetail(req)

		if err != nil {
			return nil, fmt.Errorf("failed to get bill detail: %s", err)
		}

		resourceSummarySet = append(resourceSummarySet, response.Items...)
		total = response.TotalCount

		//更改偏移量
		*offset = *offset + *req.Limit

		//如果总条数为0或者偏移大于总条数，退出循环
		if total == 0 || *offset >= *ucloud.Int(total) {
			break
		}
	}

	fmt.Printf("%s Ucloud Total: %d\n", month, len(resourceSummarySet))

	return resourceSummarySet, nil

}

func SaveUCloudBillToDB(billMonth string, resourceSummarySet []ubill.BillDetailItem) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table("ucloud_bill_resource_summary").
		Columns(
			"bill_month",
			"admin",
			"amount",
			"amount_coupon",
			"amount_free",
			"amount_real",
			"az_group_c_name",
			"charge_type",
			"create_time",
			"item_details",
			"order_no",
			"order_type",
			"project_name",
			"resource_extend_info",
			"resource_id",
			"resource_type",
			"resource_type_code",
			"show_hover",
			"start_time",
			"user_display_name",
			"user_email",
			"user_name")

	_billMonth, _err := time.Parse("2006-01", billMonth)
	if _err != nil {
		panic(_err)
	}
	_billMonthString := _billMonth.Format("2006-01-02 00:00:00")

	for _, rs := range resourceSummarySet {
		batchInsert.Values(
			_billMonthString,
			rs.Admin,
			rs.Amount,
			rs.AmountCoupon,
			rs.AmountFree,
			rs.AmountReal,
			rs.AzGroupCName,
			rs.ChargeType,
			time.Unix(int64(rs.CreateTime), 0),
			fmt.Sprintf("%+v", rs.ItemDetails),
			rs.OrderNo,
			rs.OrderType,
			rs.ProjectName,
			fmt.Sprintf("%+v", rs.ResourceExtendInfo),
			rs.ResourceId,
			rs.ResourceType,
			rs.ResourceTypeCode,
			rs.ShowHover,
			time.Unix(int64(rs.StartTime), 0),
			rs.UserDisplayName,
			rs.UserEmail,
			rs.UserName)
	}

	_, err2 := batchInsert.Save(ctx)
	if err2 != nil {
		panic(err2)
	}
}

func SyncUCloudBillToDB(month string, account CloudAccount) {
	resourceSummarySet, err := GetUCloudBill(month, account)
	if err != nil {
		panic(err)
	}
	SaveUCloudBillToDB(month, resourceSummarySet)
}
