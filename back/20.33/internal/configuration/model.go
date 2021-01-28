package configuration

type Config struct {
	Logger   LoggerConf
	DB       DBConf
	Server   GRPCConf
	RabbitMQ RabbitMQ
}

type GRPCConf struct {
	Host string `mapstructure:"server_host"`
	Port string `mapstructure:"server_port"`
}

type DBConf struct {
	User   string `mapstructure:"db_user"`
	Pass   string `mapstructure:"db_password"`
	DBName string `mapstructure:"db_database"`
	Host   string `mapstructure:"db_host"`
	Port   int    `mapstructure:"db_port"`
}

type LoggerConf struct {
	Level   string `mapstructure:"log_level"`
	File    string `mapstructure:"log_file"`
	IsProd  bool   `mapstructure:"log_trace_on"`
	TraceOn bool   `mapstructure:"log_prod_on"`
}

type RabbitMQ struct {

}