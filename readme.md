# 说明
## 腾讯云
- 接口
  - 费用中心>费用账单>账单查看>资源账单 Billing.DescribeBillResourceSummary
  - 费用中心>费用账单>账单查看>明细账单（组件名称等同于阿里云的计费项） Billing.DescribeBillDetail
- 账号权限 QcloudFinanceBillReadOnlyAccess、QcloudFinanceCostExplorerReadOnlyAccess	


## 阿里云
- 接口
  - 账单详情>明细账单（统计项：实例；统计周期：账期） Billing.DescribeInstanceBill(BillingCycle，IsBillingItem=false) 
  - 账单详情>明细账单（统计项：计费项；统计周期：账期） Billing.DescribeInstanceBill(BillingCycle，IsBillingItem=true) 
- 账号权限 AliyunBSSReadOnlyAccess

## UCloud
- 接口 
  - 财务中心>交易账单>账单明细 ListUBillDetail
- 账号权限 UBillFullAccess
- [参考资料](https://docs.ucloud.cn/api/ubill-api/list_u_bill_detail)

## AWS
- 接口 
  - GetCostAndUsage(groupby 最多两个维度)
- 账号权限 AWSBillingReadOnlyAccess
- [参考资料](https://aws.github.io/aws-sdk-go-v2/docs/code-examples/)
- [参考资料](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
- [参考资料](https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_GetCostAndUsage.html#API_GetCostAndUsage_RequestParameters)


## 操作步骤
1. 使用init-db.sql 初始化数据库
2. 添加云账号，根据说明授予相关权限，创建账号的 secretId、secretKey
3. 执行同步程序，譬如同步阿里云的2024年7月份账单： cloud-bill -m 2024-07 -c aliyun

# 遇到的问题
- 阿里云的账单数据并不准确，譬如短信已经有优惠了，但是优惠金额为0
- aws 账单的统计项和cost explorer的统计项不一致
  - 账单的统计项会有Data Transfer（流量费用），而cost explorer会将此类费用分摊到相关项目上（譬如rds、kafka、ELB、EC2等）
  - 账单的统计项Elastic Compute Cloud在cost explorer上会被分成 EC2 实例 和 EC2-其他
