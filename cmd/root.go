package cmd

import (
	"fmt"

	"github.com/liqiongfan/leopards"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 定义年月参数
var monthYear string
var appConfig AppConfig
var db *leopards.DB

// 定义主命令
var RootCmd = &cobra.Command{
	Use:   "cloud-bill",
	Short: "sync bill data from cloud providers to local database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 验证年月参数格式
		if err := validateMonthYear(monthYear); err != nil {
			return err
		}

		for _, account := range appConfig.Cloud {
			if account.Enabled {
				switch account.CloudProvider {
				case "aliyun":
					SyncAliyunBillToDB(monthYear, account)
				case "tencent":
					SyncTencentBillToDB(monthYear, account)
				case "ucloud":
					SyncUCloudBillToDB(monthYear, account)
				case "aws":
					SyncAWSBillToDB(monthYear, account)
				default:
					return fmt.Errorf("unsupported cloud platform: %s", account.CloudProvider)
				}
			}
		}

		return nil
	},
}

func init() {
	viper.SetConfigName("config") // 配置文件名（不包括扩展名）
	viper.SetConfigType("yaml")   // 如果配置文件没有扩展名，则必须设置类型
	viper.AddConfigPath(".")      // 当前目录
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 使用Unmarshal方法将整个配置文件映射到结构体
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	db, err = leopards.OpenOptions{
		User:     appConfig.Database.User,
		Password: appConfig.Database.Password,
		Host:     appConfig.Database.Host,
		Port:     appConfig.Database.Port,
		Database: appConfig.Database.DBName,
		Debug:    false, // 是否开启调试，开启调试会输出SQL到标准输出
		Dialect:  leopards.MySQL,
	}.Open()
	if err != nil {
		panic(err)
	}

	RootCmd.Flags().StringVarP(&monthYear, "bill-month", "m", "", "Specify the bill-month (e.g., 2024-08)")
	RootCmd.MarkFlagRequired("bill-month") // 标记为必填参数

}
