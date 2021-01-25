package configuration

import "time"

type SecurityConfig struct {
    EncryptCost int `env:"MAGPIE_SEC_ENCRYPT_COST" default:"10"`
    SessionKeyLength int `env:"MAGPIE_SEC_SESSION_KEY_LEN" default:"128"`
    SessionExpireTime time.Duration `env:"MAGPIE_SEC_SESSION_EXPIRE_DURATION" default:"1h"`
}
