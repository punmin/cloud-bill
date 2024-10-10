package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 定义年月参数
var monthYear string
var cloudPlatform string

// 定义主命令
var RootCmd = &cobra.Command{
	Use:   "cloud-bill",
	Short: "sync bill data from cloud providers to local database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 验证年月参数格式
		if err := validateMonthYear(monthYear); err != nil {
			return err
		}

		switch cloudPlatform {
		case "aliyun":
			SyncAliyunBillToDB(monthYear)
		case "tencent":
			SyncTencentBillToDB(monthYear)
		case "ucloud":
			SyncUCloudBillToDB(monthYear)
		case "aws":
			SyncAWSBillToDB(monthYear)
		default:
			return fmt.Errorf("unsupported cloud platform: %s", cloudPlatform)
		}

		return nil
	},
}

// 验证年月参数格式
func validateMonthYear(monthYear string) error {
	parts := strings.Split(monthYear, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format for month-year: %s. Expected format: YYYY-MM", monthYear)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid year in month-year: %s. Expected format: YYYY-MM", monthYear)
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid month in month-year: %s. Expected format: YYYY-MM", monthYear)
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month in month-year: %s. Month should be between 1 and 12", monthYear)
	}

	// 验证年份是否合理
	if year < 0 || year > time.Now().Year()+100 {
		return fmt.Errorf("invalid year in month-year: %s. Year should be within a reasonable range", monthYear)
	}

	return nil
}

// 定义年月参数的标志
func init() {
	viper.SetConfigName("config") // 配置文件名（不包括扩展名）
	viper.SetConfigType("yaml")   // 如果配置文件没有扩展名，则必须设置类型
	viper.AddConfigPath(".")      // 当前目录
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	RootCmd.Flags().StringVarP(&monthYear, "bill-month", "m", "", "Specify the bill-month (e.g., 2024-08)")
	RootCmd.MarkFlagRequired("bill-month") // 标记为必填参数

	RootCmd.Flags().StringVarP(&cloudPlatform, "cloud", "c", "", "Specify the cloud platform (tencent, aliyun, ucloud, aws)")
	RootCmd.MarkFlagRequired("cloud") // 标记为必填参数

}
