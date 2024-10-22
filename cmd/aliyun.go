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
)

const (
	//Aliyun月账单表名
	AliyunBillTableName = "aliyun_bill_resource_summary"
	//Aliyun月账单归属账号的字段名字
	AliyunMainAccountIDFieldName = "bill_account_id"
)

type AliyunCloudOperation struct{}

func CreateClient(access_key_id string, access_key_secret string) (_result *bssopenapi20171214.Client, _err error) {
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

func (cloud *AliyunCloudOperation) GetBill(billMonth string, account CloudAccount) ([]*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems, error) {
	client, _err := CreateClient(account.AccessKeyID, account.AccessKeySecret)
	if _err != nil {
		return nil, fmt.Errorf("client error has returned: %w", _err)
	}

	describeInstanceBillRequest := &bssopenapi20171214.DescribeInstanceBillRequest{
		BillingCycle:  tea.String(billMonth),
		IsBillingItem: tea.Bool(false),
		MaxResults:    tea.Int32(300),
	}
	runtime := &util.RuntimeOptions{}

	resourceSummarySet := []*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems{}

	for {
		response, err := client.DescribeInstanceBillWithOptions(describeInstanceBillRequest, runtime)
		if _, ok := err.(*tea.SDKError); ok {
			return nil, fmt.Errorf("an API error has returned: %w", err)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get bill detail: %w", err)
		}

		resourceSummarySet = append(resourceSummarySet, response.Body.Data.Items...)

		if response.Body.Data.NextToken == nil {
			break
		}
		describeInstanceBillRequest.NextToken = response.Body.Data.NextToken

		//防止请求超过云平台的限制
		sleepForFraction(account.FetchPerSecond)
	}

	fmt.Printf("%s %s Aliyun Total: %d\n", billMonth, account.MainAccountID, len(resourceSummarySet))

	return resourceSummarySet, nil

}

func (cloud *AliyunCloudOperation) SaveBill(billMonth string, account CloudAccount, resourceSummarySet []*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems) {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table(AliyunBillTableName).
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

func (cloud *AliyunCloudOperation) HasBill(billMonth string, account CloudAccount) bool {
	return HasBill(billMonth, account, AliyunBillTableName, AliyunMainAccountIDFieldName)
}

func SyncAliyunBillToDB(billMonth string, account CloudAccount) {
	operation := CommonBillOperation[*bssopenapi20171214.DescribeInstanceBillResponseBodyDataItems]{
		BillOperation: &AliyunCloudOperation{},
	}
	operation.SyncBill(billMonth, account)
}
