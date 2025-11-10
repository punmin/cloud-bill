# 说明
## 注意事项
- account_info账户信息表需要手动插入数据，将配置文件的main_account_id,account_alias_name和关联起来
```sql
insert into account_info(bill_account_id, bill_account_alias)values('your_main_account_id', 'your_account_alias_name');
```

- 备份时不备份视图，要不然导入时，会提示权限问题
```sql
-- 导出所有表，不包括视图
mysqldump -uroot -h your_host --no-create-db --ignore-table=cloud_bill.all_bill_resource_summary   -p cloud_bill  > cloud_bill.sql

-- 导入所有表
mysql -ucloud_bill -h your_host  -p cloud_bill  < cloud_bill.sql
```


## 腾讯云
- 接口
  - 费用中心>费用账单>账单查看>资源账单 Billing.DescribeBillResourceSummary
  - 费用中心>费用账单>账单查看>明细账单（组件名称等同于阿里云的计费项） Billing.DescribeBillDetail
- 账号权限 QcloudFinanceBillReadOnlyAccess、QcloudFinanceCostExplorerReadOnlyAccess	
- [接口说明](https://console.cloud.tencent.com/api/explorer?Product=billing&Version=2018-07-09&Action=DescribeBillResourceSummary)
- 请求频率限制 5次/秒

## 阿里云
- 接口
  - 账单详情>明细账单（统计项：实例；统计周期：账期） Billing.DescribeInstanceBill(BillingCycle，IsBillingItem=false) 
  - 账单详情>明细账单（统计项：计费项；统计周期：账期） Billing.DescribeInstanceBill(BillingCycle，IsBillingItem=true) 
- 账号权限 AliyunBSSReadOnlyAccess
- [接口说明](https://next.api.aliyun.com/api/BssOpenApi/2017-12-14/DescribeInstanceBill)
- 请求频率限制 10次/秒

## UCloud
- 接口 
  - 财务中心>交易账单>账单明细 ListUBillDetail
- 账号权限 UBillFullAccess
- [参考资料](https://docs.ucloud.cn/api/ubill-api/list_u_bill_detail)
- 请求频率限制 未找到文档描述，建议设置为5次/秒

## AWS
- 接口 
  - GetCostAndUsage(groupby 最多两个维度)
- 账号权限 AWSBillingReadOnlyAccess
- [参考资料](https://aws.github.io/aws-sdk-go-v2/docs/code-examples/)
- [参考资料](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
- [参考资料](https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_GetCostAndUsage.html#API_GetCostAndUsage_RequestParameters)
- 请求频率限制 未找到文档描述，建议设置为5次/秒
- 默认支持最近14个月，超过14个月的历史数据(最多38个月)，需要做以下设置：账单与成本管理>首选项和设置>成本管理首选项>Cost Explorer > 勾选每月粒度的多年数据（最多48小时生效）

## 操作步骤
1. 使用init-db.sql 初始化数据库
2. 添加云账号，根据说明授予相关权限，创建账号的 secretId、secretKey
3. 执行同步程序

## 参数说明
- cloud-bill -s last 同步上个月份
- cloud-bill -s 2024-07 同步指定单个月份
- cloud-bill -s 2024-07 -e 2024-08 同步指定多个月份

# 遇到的问题
- 阿里云的账单数据并不准确，譬如短信已经有优惠了，但是优惠金额为0
- aws 账单的统计项和cost explorer的统计项不一致
  - 账单的统计项会有Data Transfer（流量费用），而cost explorer会将此类费用分摊到相关项目上（譬如rds、kafka、ELB、EC2等）
  - 账单的统计项Elastic Compute Cloud在cost explorer上会被分成 EC2 实例 和 EC2-其他
- 腾讯云费用账单中的对象存储的折扣是错的，而消耗账单的折扣是对的