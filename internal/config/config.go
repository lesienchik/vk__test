package config

type Config struct {
	Api      Api      `json:"api"`
	Postgres Postgres `json:"postgres"`
	Logic    Logic    `json:"logic"`
	Email    Email    `json:"email"`
	LogLevel string   `json:"log_level"`
}

type Api struct {
	Addr         string   `json:"addr"`
	ReadTimeout  int      `json:"read_timeout"`
	WriteTimeout int      `json:"write_timeout"`
	IdleTimeout  int      `json:"idle_timeout"`
	AllowOrigins []string `json:"allow_origins"`
}

type Logic struct {
	SecretKey string `env:"SECRET_KEY,notEmpty"`
}

type Postgres struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	Sslmode  string `json:"sslmode"`
	MaxConns int    `json:"max_conns"`
	AppName  string `json:"app_name"`
	User     string `env:"PG_USER,notEmpty"`
	Password string `env:"PG_PASSWORD,notEmpty"`
}

type Email struct {
	Addr     string `json:"addr"`
	Site     string `json:"site"`
	Password string `env:"EMAIL_PASSWORD,notEmpty"`
}
