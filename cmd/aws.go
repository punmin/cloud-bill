package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/liqiongfan/leopards"
)

const (
	//AWS月账单表名
	AWSBillTableName = "aws_bill_resource_summary"
	//AWS月账单归属账号的字段名字
	AWSMainAccountIDFieldName = "bill_account_id"
)

type AWSCloudOperation struct {
	BillTableName          string
	MainAccountIDFieldName string
}

type AWSBill struct {
	Service       string
	Region        string
	UnblendedCost string
	ExchangeRate  string
}

func GetAWSTimePeriod(month string) (string, string) {
	// 解析输入的年月字符串
	t, err := time.Parse("2006-01", month)
	if err != nil {
		fmt.Println("Invalid input format. Please use YYYY-MM.")
		return "", ""
	}

	// 计算下一个月
	next := t.AddDate(0, 1, 0)

	//输入2024-08，输出2024-08-01 and 2024-09-01
	return month + "-01", next.Format("2006-01") + "-01"
}

func (cloud *AWSCloudOperation) GetBill(billMonth string, account CloudAccount) ([]*AWSBill, error) {
	exchange_rate := appConfig.UsdToCnyExchangeRate
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(account.AccessKeyID, account.AccessKeySecret, "")),
	)

	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	svc := costexplorer.NewFromConfig(cfg)

	// 设置查询时间范围
	start, end := GetAWSTimePeriod(billMonth)
	resourceSummarySet := make([]*AWSBill, 0)

	// 创建 GetCostAndUsageInput 请求
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
		Granularity: types.GranularityMonthly,
		Metrics:     []string{"UnblendedCost"},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("REGION"),
			},
		},
	}

	// 发送请求
	for {
		result, err := svc.GetCostAndUsage(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, resultByTime := range result.ResultsByTime {
			for _, group := range resultByTime.Groups {
				for _, metric := range group.Metrics {
					resourceSummarySet = append(resourceSummarySet, &AWSBill{
						Service:       group.Keys[0],
						Region:        group.Keys[1],
						UnblendedCost: *metric.Amount,
						ExchangeRate:  fmt.Sprintf("%.4f", exchange_rate),
					})
					//fmt.Printf("key: %s, metric: %s, Amount: %s, Unit: %s\n", group.Keys, metrics_name, *metric.Amount, *metric.Unit)
				}
			}
		}
		if result.NextPageToken == nil {
			break
		}
		input.NextPageToken = result.NextPageToken

		//防止请求超过云平台的限制
		sleepForFraction(account.FetchPerSecond)
	}

	fmt.Printf("%s %s AWS Total: %d\n", billMonth, account.MainAccountID, len(resourceSummarySet))

	return resourceSummarySet, nil
}

func (cloud *AWSCloudOperation) SaveBill(billMonth string, account CloudAccount, resourceSummarySet []*AWSBill) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table(cloud.BillTableName).
		Columns(
			"bill_month",
			"service",
			"region",
			"unblended_cost",
			"exchange_rate",
			"bill_account_id",
		)

	_billMonth, _err := time.Parse("2006-01", billMonth)
	if _err != nil {
		panic(_err)
	}
	_billMonthString := _billMonth.Format("2006-01-02 00:00:00")

	for _, rs := range resourceSummarySet {
		batchInsert.Values(
			_billMonthString,
			rs.Service,
			rs.Region,
			rs.UnblendedCost,
			rs.ExchangeRate,
			account.MainAccountID,
		)
	}

	_, err2 := batchInsert.Save(ctx)
	if err2 != nil {
		panic(err2)
	}
}

func (cloud *AWSCloudOperation) HasBill(billMonth string, account CloudAccount) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	var result = struct {
		Count int `json:"count"`
	}{}

	err := db.Query().From(cloud.BillTableName).Select(leopards.As(leopards.Count(`id`), `count`)).Where(
		leopards.And(
			leopards.EQ("bill_month", fmt.Sprintf("%s-01 00:00:00", billMonth)),
			leopards.EQ(cloud.MainAccountIDFieldName, account.MainAccountID),
		),
	).Scan(ctx, &result)

	if err != nil {
		panic(err)
	}

	return result.Count > 0

}
func SyncAWSBillToDB(billMonth string, account CloudAccount) {
	operation := CommonBillOperation[*AWSBill]{
		BillOperation: &AWSCloudOperation{
			BillTableName:          AWSBillTableName,
			MainAccountIDFieldName: AWSMainAccountIDFieldName,
		},
	}
	operation.SyncBill(billMonth, account)
}
