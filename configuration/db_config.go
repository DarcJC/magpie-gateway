package configuration

type DatabaseConfig struct {
    DBHost string `env:"MAGPIE_DB_HOST" default:"localhost"`
    DBPort int `env:"MAGPIE_DB_PORT" default:"5432"`
    DBUser string `env:"MAGPIE_DB_USER" default:"root"`
    DBPassword string `env:"MAGPIE_DB_PASSWORD" default:"root"`
    DBName string `env:"MAGPIE_DB_NAME" default:"magpie_gateway"`
    DBSSLMode string `env:"MAGPIE_DB_SSL" default:"disable"`
    DBTimezone string `env:"MAGPIE_DB_TIMEZONE" default:"Asia/Shanghai"`
    RedisURI string `env:"MAGPIE_REDIS_URI" default:"redis://localhost:6379/"`
    DBMock bool `env:"MAGPIE_DB_MOCK" default:"false"`
}
