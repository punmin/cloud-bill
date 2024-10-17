// This file is auto-generated, don't edit it. Thanks.
package cmd

import (
	"context"
	"fmt"
	"time"

	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/liqiongfan/leopards"
)

// Description:
//
// 使用AK&SK初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient(access_key_id string, access_key_secret string) (_result *bssopenapi20171214.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		AccessKeyId:     tea.String(access_key_id),
		AccessKeySecret: tea.String(access_key_secret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/BssOpenApi
	config.Endpoint = tea.String("business.aliyuncs.com")
	_result = &bssopenapi20171214.Client{}
	_result, _err = bssopenapi20171214.NewClient(config)
	return _result, _err
}

func GetAliyunBill(month string, account CloudAccount) ([]*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems, error) {
	client, _err := CreateClient(account.AccessKeyID, account.AccessKeySecret)
	if _err != nil {
		return nil, fmt.Errorf("client error has returned: %s", _err)
	}

	describeInstanceBillRequest := &bssopenapi20171214.DescribeInstanceBillRequest{
		BillingCycle:  tea.String(month),
		IsBillingItem: tea.Bool(false),
		MaxResults:    tea.Int32(300),
	}
	runtime := &util.RuntimeOptions{}

	resourceSummarySet := []*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems{}

	for {
		response, err := client.DescribeInstanceBillWithOptions(describeInstanceBillRequest, runtime)
		if _, ok := err.(*tea.SDKError); ok {
			return nil, fmt.Errorf("An API error has returned: %s", err)
		}

		if err != nil {
			return nil, fmt.Errorf("Failed to get bill detail: %s", err)
		}

		resourceSummarySet = append(resourceSummarySet, response.Body.Data.Items...)

		if response.Body.Data.NextToken == nil {
			break
		}
		describeInstanceBillRequest.NextToken = response.Body.Data.NextToken

	}

	fmt.Printf("%s %s Aliyun Total: %d\n", month, account.MainAccountID, len(resourceSummarySet))

	return resourceSummarySet, nil

}

func SaveAliyunBillToDB(billMonth string, resourceSummarySet []*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems) {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table("aliyun_bill_resource_summary").
		Columns(
			"bill_month",
			"adjust_amount",
			"bill_account_id",
			"bill_account_name",
			"billing_date",
			"billing_item",
			"billing_item_code",
			"billing_type",
			"biz_type",
			"cash_amount",
			"commodity_code",
			"cost_unit",
			"currency",
			"deducted_by_cash_coupons",
			"deducted_by_coupons",
			"deducted_by_prepaid_card",
			"deducted_by_resource_package",
			"instance_config",
			"instance_id",
			"instance_spec",
			"internet_ip",
			"intranet_ip",
			"invoice_discount",
			"discount",
			"item",
			"item_name",
			"list_price",
			"list_price_unit",
			"nick_name",
			"outstanding_amount",
			"owner_id",
			"payment_amount",
			"pip_code",
			"pretax_amount",
			"pretax_gross_amount",
			"product_code",
			"product_detail",
			"product_name",
			"product_type",
			"region",
			"resource_group",
			"service_period",
			"service_period_unit",
			"subscription_type",
			"tag",
			"usage",
			"usage_unit",
			"zone",
		)

	_billMonth, _err := time.Parse("2006-01", billMonth)
	if _err != nil {
		panic(_err)
	}
	_billMonthString := _billMonth.Format("2006-01-02 00:00:00")

	for _, rs := range resourceSummarySet {
		billingDate := rs.BillingDate
		if *rs.BillingDate == "" {
			billingDate = nil
		}

		listPrice := rs.ListPrice
		if *rs.ListPrice == "" {
			listPrice = tea.String("0")
		}

		usage := rs.Usage
		if *rs.Usage == "" {
			usage = tea.String("0")
		}

		discount := *tea.Float32(0)

		if *rs.PretaxAmount != 0.0 {
			discount = *rs.PretaxAmount / *rs.PretaxGrossAmount
		}

		batchInsert.Values(
			_billMonthString,
			rs.AdjustAmount,
			rs.BillAccountID,
			rs.BillAccountName,
			billingDate,
			rs.BillingItem,
			rs.BillingItemCode,
			rs.BillingType,
			rs.BizType,
			rs.CashAmount,
			rs.CommodityCode,
			rs.CostUnit,
			rs.Currency,
			rs.DeductedByCashCoupons,
			rs.DeductedByCoupons,
			rs.DeductedByPrepaidCard,
			rs.DeductedByResourcePackage,
			rs.InstanceConfig,
			rs.InstanceID,
			rs.InstanceSpec,
			rs.InternetIP,
			rs.IntranetIP,
			rs.InvoiceDiscount,
			discount,
			rs.Item,
			rs.ItemName,
			listPrice,
			rs.ListPriceUnit,
			rs.NickName,
			rs.OutstandingAmount,
			rs.OwnerID,
			rs.PaymentAmount,
			rs.PipCode,
			rs.PretaxAmount,
			rs.PretaxGrossAmount,
			rs.ProductCode,
			rs.ProductDetail,
			rs.ProductName,
			rs.ProductType,
			rs.Region,
			rs.ResourceGroup,
			rs.ServicePeriod,
			rs.ServicePeriodUnit,
			rs.SubscriptionType,
			rs.Tag,
			usage,
			rs.UsageUnit,
			rs.Zone,
		)
	}

	_, err2 := batchInsert.Save(ctx)
	if err2 != nil {
		panic(err2)
	}
}

func HasAliyunBill(billMonth string, account CloudAccount) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	var result = struct {
		Count int `json:"count"`
	}{}

	err := db.Query().From("aliyun_bill_resource_summary").Select(leopards.As(leopards.Count(`id`), `count`)).Where(
		leopards.And(
			leopards.EQ("bill_month", fmt.Sprintf("%s-01 00:00:00", billMonth)),
			leopards.EQ(`bill_account_id`, account.MainAccountID),
		),
	).Scan(ctx, &result)

	if err != nil {
		panic(err)
	}

	return result.Count > 0

}

func SyncAliyunBillToDB(month string, account CloudAccount) {
	if HasAliyunBill(month, account) {
		fmt.Printf("%s bill for %s has been synced\n", account.AccountAliasName, month)
		return
	}

	resourceSummarySet, err := GetAliyunBill(month, account)
	if err != nil {
		panic(err)
	}
	SaveAliyunBillToDB(month, resourceSummarySet)
}
