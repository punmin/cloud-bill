package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func StringTags(tags []*billing.BillTagInfo) string {
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		tagsJSON = []byte("[]") // 返回空数组作为默认值
	}
	return string(tagsJSON)
}

func GetTencentBill(month string, account CloudAccount) ([]*billing.BillResourceSummary, error) {
	credential := common.NewCredential(account.AccessKeyID, account.AccessKeySecret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "billing.tencentcloudapi.com"
	client, err := billing.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %s", err)
	}

	request := billing.NewDescribeBillResourceSummaryRequest()
	request.Month = common.StringPtr(month)
	request.Limit = common.Uint64Ptr(1000)
	request.NeedRecordNum = common.Int64Ptr(1)

	offset := common.Uint64Ptr(0)
	total := common.Int64Ptr(0)

	resourceSummarySet := []*billing.BillResourceSummary{}

	for {
		request.Offset = offset
		response, err := client.DescribeBillResourceSummary(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return nil, fmt.Errorf("An API error has returned: %s", err)
		}

		if err != nil {
			return nil, fmt.Errorf("Failed to get bill detail: %s", err)
		}

		resourceSummarySet = append(resourceSummarySet, response.Response.ResourceSummarySet...)
		total = response.Response.Total

		//更改偏移量
		*offset = *offset + *request.Limit

		//如果总条数为0或者偏移大于总条数，退出循环
		if *total == 0 || *offset >= uint64(*total) {
			break
		}
	}

	fmt.Printf("%s Tencent Total: %d\n", month, len(resourceSummarySet))

	return resourceSummarySet, nil
}

func SaveTencentBillToDB(resourceSummarySet []*billing.BillResourceSummary) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table("tencent_bill_resource_summary").
		Columns(
			"bill_month",
			"tags",
			"action_type_name",
			"business_code",
			"business_code_name",
			"cash_pay_amount",
			"config_desc",
			"discount",
			"extend_field1",
			"extend_field2",
			"extend_field3",
			"extend_field4",
			"extend_field5",
			"fee_begin_time",
			"fee_end_time",
			"incentive_pay_amount",
			"instance_type",
			"operate_uin",
			"order_id",
			"original_cost_with_ri",
			"original_cost_with_sp",
			"owner_uin",
			"pay_mode_name",
			"pay_time",
			"payer_uin",
			"product_code",
			"product_code_name",
			"project_name",
			"real_total_cost",
			"reduce_type",
			"region_id",
			"region_name",
			"resource_id",
			"resource_name",
			"total_cost",
			"transfer_pay_amount",
			"voucher_pay_amount",
			"zone_name")

	for _, rs := range resourceSummarySet {
		payTime := rs.PayTime
		if *rs.PayTime == "0000-00-00 00:00:00" {
			payTime = nil
		}

		batchInsert.Values(
			rs.BillMonth,
			StringTags(rs.Tags),
			rs.ActionTypeName,
			rs.BusinessCode,
			rs.BusinessCodeName,
			rs.CashPayAmount,
			rs.ConfigDesc,
			rs.Discount,
			rs.ExtendField1,
			rs.ExtendField2,
			rs.ExtendField3,
			rs.ExtendField4,
			rs.ExtendField5,
			rs.FeeBeginTime,
			rs.FeeEndTime,
			rs.IncentivePayAmount,
			rs.InstanceType,
			rs.OperateUin,
			rs.OrderId,
			rs.OriginalCostWithRI,
			rs.OriginalCostWithSP,
			rs.OwnerUin,
			rs.PayModeName,
			payTime,
			rs.PayerUin,
			rs.ProductCode,
			rs.ProductCodeName,
			rs.ProjectName,
			rs.RealTotalCost,
			rs.ReduceType,
			rs.RegionId,
			rs.RegionName,
			rs.ResourceId,
			rs.ResourceName,
			rs.TotalCost,
			rs.TransferPayAmount,
			rs.VoucherPayAmount,
			rs.ZoneName)
	}

	_, err2 := batchInsert.Save(ctx)
	if err2 != nil {
		panic(err2)
	}
}

func SyncTencentBillToDB(month string, account CloudAccount) {
	resourceSummarySet, err := GetTencentBill(month, account)
	if err != nil {
		panic(err)
	}
	SaveTencentBillToDB(resourceSummarySet)
}
