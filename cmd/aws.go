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
	"github.com/spf13/viper"
)

type AWSBill struct {
	Service       string
	Region        string
	UnblendedCost string
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

func GetAWSBill(month string) ([]*AWSBill, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(viper.GetString("cloud.aws.access_key_id"), viper.GetString("cloud.aws.access_key_secret"), "")),
	)

	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	svc := costexplorer.NewFromConfig(cfg)

	// 设置查询时间范围
	start, end := GetAWSTimePeriod(month)

	var resourceSummarySet []*AWSBill
	resourceSummarySet = make([]*AWSBill, 0)

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
					})
					//fmt.Printf("key: %s, metric: %s, Amount: %s, Unit: %s\n", group.Keys, metrics_name, *metric.Amount, *metric.Unit)
				}
			}
		}
		if result.NextPageToken == nil {
			break
		}
		input.NextPageToken = result.NextPageToken
	}

	fmt.Printf("Total: %d\n", len(resourceSummarySet))

	return resourceSummarySet, nil
}

func SaveAWSBillToDB(billMonth string, resourceSummarySet []*AWSBill) {
	db, err := leopards.OpenOptions{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Database: viper.GetString("database.dbname"),
		Debug:    false, // 是否开启调试，开启调试会输出SQL到标准输出
		Dialect:  leopards.MySQL,
	}.Open()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	batchInsert := db.Insert().Table("aws_bill_resource_summary").
		Columns(
			"bill_month",
			"service",
			"region",
			"unblended_cost",
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
		)
	}

	_, err2 := batchInsert.Save(ctx)
	if err2 != nil {
		panic(err2)
	}
}
func SyncAWSBillToDB(month string) {
	resourceSummarySet, err := GetAWSBill(month)
	if err != nil {
		panic(err)
	}
	SaveAWSBillToDB(month, resourceSummarySet)
}
