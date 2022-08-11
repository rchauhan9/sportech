package config

type Config struct {
	Port        string
	Database    DatabaseConfig
	Environment string
}

type DatabaseConfig struct {
	URL string
}
