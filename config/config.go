package config

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Environment string
}

type ServerConfig struct {
	GRPCAddress string `mapstructure:"grpc-address"`
	HTTPAddress string `mapstructure:"http-address"`
}

type DatabaseConfig struct {
	URL string
}
