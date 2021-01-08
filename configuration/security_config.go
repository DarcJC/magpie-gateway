package configuration

type SecurityConfig struct {
	EncryptCost int `env:"MAGPIE_SEC_ENCRYPT_COST" default:"10"`
}
