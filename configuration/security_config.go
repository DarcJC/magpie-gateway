package configuration

type SecurityConfig struct {
    EncryptCost int `env:"MAGPIE_SEC_ENCRYPT_COST" default:"10"`
    SessionKeyLength int `env:"MAGPIE_SEC_SESSION_KEY_LEN" default:"128"`
    SessionJWTSecret []byte `env:"MAGPIE_SEC_SESSION_JWT_SECRET" default:"y72GNO3IwZf5BzX&Fsl0kjydPmbJmiXs"`
    SessionJWTExpireTime string `env:"MAGPIE_SEC_SESSION_JWT_EXPIRE" default:"60m"`
}
