package cmd

// DatabaseConfig 代表数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Debug    bool   `mapstructure:"debug"`
}

// CloudAccount 代表云账户配置
type CloudAccount struct {
	AccountAliasName string `mapstructure:"account_alias_name"`
	CloudProvider    string `mapstructure:"cloud_provider"`
	MainAccountID    string `mapstructure:"main_account_id"`
	AccessKeyID      string `mapstructure:"access_key_id"`
	AccessKeySecret  string `mapstructure:"access_key_secret"`
	Enabled          bool   `mapstructure:"enabled"`
}

// AppConfig 代表整个应用的配置
type AppConfig struct {
	Database             DatabaseConfig `mapstructure:"database"`
	UsdToCnyExchangeRate float64        `mapstructure:"usd_to_cny_exchange_rate"`
	Cloud                []CloudAccount `mapstructure:"cloud"`
}
