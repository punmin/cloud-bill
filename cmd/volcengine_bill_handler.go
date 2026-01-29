package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/volcengine/volcengine-go-sdk/service/billing"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

const (
	//Volcengine月账单表名
	VolcengineBillTableName = "volcengine_bill_resource_summary"
	//Volcengine月账单归属账号的字段名字
	VolcengineMainAccountIDFieldName = "owner_id"
)

type VolcengineBillHandler struct{}

func CreateVolcengineClient(access_key_id string, access_key_secret string) (*billing.BILLING, error) {
	sess, err := session.NewSession(&volcengine.Config{
		Credentials: credentials.NewStaticCredentials(
			access_key_id,
			access_key_secret,
			"",
		),
		Region: volcengine.String("cn-beijing"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	client := billing.New(sess)
	return client, nil
}

func (handler *VolcengineBillHandler) GetBill(billMonth string, account CloudAccount) ([]*billing.ListForListBillDetailOutput, error) {
	client, _err := CreateVolcengineClient(account.AccessKeyID, account.AccessKeySecret)
	if _err != nil {
		return nil, fmt.Errorf("client error has returned: %w", _err)
	}

	listBillDetailInput := &billing.ListBillDetailInput{
		BillPeriod:    volcengine.String(billMonth), // Format: YYYY-MM
		Limit:         volcengine.Int32(300),        //官方限制每次最多返回300
		Offset:        volcengine.Int32(0),
		GroupPeriod:   volcengine.Int32(0), //使用账期作为统计周期
		GroupTerm:     volcengine.Int32(0), //按计费项作为统计项，因为按计费项带实例规格信息，按实例不带实例规格信息
		IgnoreZero:    volcengine.Int32(1), //忽略折后价为0的数据
		NeedRecordNum: volcengine.Int32(1), //返回总记录数
	}

	resourceSummarySet := []*billing.ListForListBillDetailOutput{}

	for {
		response, err := client.ListBillDetail(listBillDetailInput)

		if err != nil {
			return nil, fmt.Errorf("failed to get bill detail: %w", err)
		}

		if response.List != nil {
			resourceSummarySet = append(resourceSummarySet, response.List...)
		}

		listBillDetailInput.Offset = volcengine.Int32(*response.Offset + *listBillDetailInput.Limit)
		if *listBillDetailInput.Offset >= *response.Total {
			break
		}

		//防止请求超过云平台的限制
		sleepForFraction(account.FetchPerSecond)

	}

	fmt.Printf("%s %s Volcengine Total: %d\n", billMonth, account.MainAccountID, len(resourceSummarySet))

	return resourceSummarySet, nil
}

func (handler *VolcengineBillHandler) SaveBill(billMonth string, account CloudAccount, resourceSummarySet []*billing.ListForListBillDetailOutput) {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table(VolcengineBillTableName).
		Columns(
			"bill_month",
			"bill_category",
			"bill_detail_id",
			"bill_id",
			"bill_period",
			"billing_function",
			"billing_method_code",
			"billing_mode",
			"busi_period",
			"business_mode",
			"config_name",
			"configuration_code",
			"count",
			"country_region",
			"coupon_amount",
			"credit_carried_amount",
			"currency",
			"currency_settlement",
			"deduction_count",
			"deduction_use_duration",
			"discount_bill_amount",
			"discount_biz_billing_function",
			"discount_biz_measure_interval",
			"discount_biz_unit_price",
			"discount_biz_unit_price_interval",
			"discount_info",
			"effective_factor",
			"element",
			"element_code",
			"exchange_rate",
			"expand_field",
			"expense_begin_time",
			"expense_date",
			"expense_end_time",
			"factor",
			"factor_code",
			"formula",
			"instance_name",
			"instance_no",
			"main_contract_number",
			"market_price",
			"measure_interval",
			"original_bill_amount",
			"original_order_no",
			"owner_customer_name",
			"owner_id",
			"owner_user_name",
			"paid_amount",
			"payable_amount",
			"payer_customer_name",
			"payer_id",
			"payer_user_name",
			"posttax_amount",
			"pre_tax_payable_amount",
			"preferential_bill_amount",
			"pretax_amount",
			"pretax_real_value",
			"price",
			"price_interval",
			"price_unit",
			"product",
			"product_zh",
			"project",
			"project_display_name",
			"real_value",
			"region",
			"region_code",
			"reservation_instance",
			"round_amount",
			"saving_plan_deduction_discount_amount",
			"saving_plan_deduction_sp_id",
			"saving_plan_original_amount",
			"seller_customer_name",
			"seller_id",
			"seller_user_name",
			"selling_mode",
			"settle_payable_amount",
			"settle_posttax_amount",
			"settle_pre_tax_payable_amount",
			"settle_pretax_amount",
			"settle_pretax_real_value",
			"settle_real_value",
			"settle_tax",
			"settlement_type",
			"solution_zh",
			"subject_name",
			"tag",
			"tax",
			"tax_rate",
			"trade_time",
			"unit",
			"unpaid_amount",
			"use_duration",
			"use_duration_unit",
			"zone",
			"zone_code",
		)

	_billMonth, _err := time.Parse("2006-01", billMonth)
	if _err != nil {
		panic(_err)
	}
	_billMonthString := _billMonth.Format("2006-01-02 00:00:00")

	for _, rs := range resourceSummarySet {
		batchInsert.Values(
			_billMonthString,
			rs.BillCategory,
			rs.BillDetailId,
			rs.BillID,
			rs.BillPeriod,
			rs.BillingFunction,
			rs.BillingMethodCode,
			rs.BillingMode,
			rs.BusiPeriod,
			rs.BusinessMode,
			rs.ConfigName,
			rs.ConfigurationCode,
			rs.Count,
			rs.CountryRegion,
			rs.CouponAmount,
			rs.CreditCarriedAmount,
			rs.Currency,
			rs.CurrencySettlement,
			rs.DeductionCount,
			rs.DeductionUseDuration,
			rs.DiscountBillAmount,
			rs.DiscountBizBillingFunction,
			rs.DiscountBizMeasureInterval,
			rs.DiscountBizUnitPrice,
			rs.DiscountBizUnitPriceInterval,
			rs.DiscountInfo,
			rs.EffectiveFactor,
			rs.Element,
			rs.ElementCode,
			rs.ExchangeRate,
			rs.ExpandField,
			rs.ExpenseBeginTime,
			rs.ExpenseDate,
			rs.ExpenseEndTime,
			rs.Factor,
			rs.FactorCode,
			rs.Formula,
			rs.InstanceName,
			rs.InstanceNo,
			rs.MainContractNumber,
			rs.MarketPrice,
			rs.MeasureInterval,
			rs.OriginalBillAmount,
			rs.OriginalOrderNo,
			rs.OwnerCustomerName,
			rs.OwnerID,
			rs.OwnerUserName,
			rs.PaidAmount,
			rs.PayableAmount,
			rs.PayerCustomerName,
			rs.PayerID,
			rs.PayerUserName,
			rs.PosttaxAmount,
			rs.PreTaxPayableAmount,
			rs.PreferentialBillAmount,
			rs.PretaxAmount,
			rs.PretaxRealValue,
			rs.Price,
			rs.PriceInterval,
			rs.PriceUnit,
			rs.Product,
			rs.ProductZh,
			rs.Project,
			rs.ProjectDisplayName,
			rs.RealValue,
			rs.Region,
			rs.RegionCode,
			rs.ReservationInstance,
			rs.RoundAmount,
			rs.SavingPlanDeductionDiscountAmount,
			rs.SavingPlanDeductionSpID,
			rs.SavingPlanOriginalAmount,
			rs.SellerCustomerName,
			rs.SellerID,
			rs.SellerUserName,
			rs.SellingMode,
			rs.SettlePayableAmount,
			rs.SettlePosttaxAmount,
			rs.SettlePreTaxPayableAmount,
			rs.SettlePretaxAmount,
			rs.SettlePretaxRealValue,
			rs.SettleRealValue,
			rs.SettleTax,
			rs.SettlementType,
			rs.SolutionZh,
			rs.SubjectName,
			rs.Tag,
			rs.Tax,
			rs.TaxRate,
			rs.TradeTime,
			rs.Unit,
			rs.UnpaidAmount,
			rs.UseDuration,
			rs.UseDurationUnit,
			rs.Zone,
			rs.ZoneCode,
		)
	}

	_, err2 := batchInsert.Save(ctx)
	if err2 != nil {
		panic(err2)
	}
}

func (handler *VolcengineBillHandler) HasBill(billMonth string, account CloudAccount) bool {
	return HasBill(billMonth, account, VolcengineBillTableName, VolcengineMainAccountIDFieldName)
}

func SyncVolcengineBillToDB(billMonth string, account CloudAccount) {
	executor := CloudBillExecutor[*billing.ListForListBillDetailOutput]{
		handler: &VolcengineBillHandler{},
	}
	executor.SyncBill(billMonth, account)
}
