package configuration

type Config struct {
	Logger LoggerConf
	DB     DBConf
	Server GRPCConf
	Rabbit RabbitMQ
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
	User       string `mapstructure:"rabbit_user"`
	Pass       string `mapstructure:"rabbit_user"`
	Host       string `mapstructure:"rabbit_host"`
	Port       int    `mapstructure:"rabbit_port"`
	Durable    bool   `mapstructure:"rabbit_durable"`
	AutoDelete bool   `mapstructure:"rabbit_autodelete"`
	NoWait     bool   `mapstructure:"rabbit_no_wait"`
	Internal   bool   `mapstructure:"rabbit_internal"`
	Name       string `mapstructure:"rabbit_name"`
	Kind       string `mapstructure:"rabbit_kind"`
	Key        string `mapstructure:"rabbit_key"`
}
