package configuration

type SecurityConfig struct {
    EncryptCost int `env:"MAGPIE_SEC_ENCRYPT_COST" default:"10"`
    SessionKeyLength int `env:"MAGPIE_SEC_SESSION_KEY_LEN" default:"128"`
}
