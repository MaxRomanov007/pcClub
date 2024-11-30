package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type SQLServerConfig struct {
	Username               string `yaml:"username"`
	Password               string `yaml:"password"`
	Hostname               string `yaml:"hostname"`
	Port                   int    `yaml:"port"`
	Database               string `yaml:"database"`
	Encrypt                bool   `yaml:"encrypt"`
	TrustServerCertificate bool   `yaml:"trust_server_certificate"`
	AppName                string `yaml:"app_name"`
}

type RedisConfig struct {
	Address    string        `yaml:"address"`
	Password   string        `yaml:"password"`
	Database   int           `yaml:"database"`
	DefaultTTL time.Duration `yaml:"default_ttl"`
}

type DatabaseConfig struct {
	SQLServer *SQLServerConfig `yaml:"sql_server"`
	Redis     *RedisConfig     `yaml:"redis"`
}

type HTTPSServerConfig struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	SSLCert     string        `yaml:"ssl_cert"`
	SSLKey      string        `yaml:"ssl_key"`
}

type PCSImagesConfig struct {
	Path string `yaml:"path"`
}

type DishesImagesConfig struct {
	Path string `yaml:"path"`
}

type ImagesConfig struct {
	PCS    *PCSImagesConfig    `yaml:"pcs"`
	Dishes *DishesImagesConfig `yaml:"dish"`
}

type AccessTokenConfig struct {
	Secret string        `yaml:"secret"`
	TTL    time.Duration `yaml:"ttl"`
}

type RefreshTokenConfig struct {
	Secret     string        `yaml:"secret"`
	TTL        time.Duration `yaml:"ttl"`
	CookieName string        `yaml:"cookie_name"`
}

type AuthConfig struct {
	UrlPath string              `yaml:"url_path"`
	Access  *AccessTokenConfig  `yaml:"access"`
	Refresh *RefreshTokenConfig `yaml:"refresh"`
}

type UserConfig struct {
	AdminRoleName string `yaml:"admin_role_name"`
}

type Config struct {
	Env         string             `yaml:"env"`
	Database    *DatabaseConfig    `yaml:"database"`
	HttpsServer *HTTPSServerConfig `yaml:"https_server"`
	Images      *ImagesConfig      `yaml:"images"`
	Auth        *AuthConfig        `yaml:"auth"`
	User        *UserConfig        `yaml:"user"`
}

func MustLoad() *Config {
	path := MustGetPath()

	return MustLoadByPath(path)
}

func MustGetPath() string {
	path := getPath()
	if path == "" {
		log.Fatal("config path not set")
	}

	return path
}

func getPath() string {
	if path := getPathByEnv(); path != "" {
		return path
	}
	return getPathByFlag()
}

func getPathByEnv() string {
	path := os.Getenv("CONFIG_PATH")
	return path
}

func getPathByFlag() string {
	var path string

	flag.StringVar(&path, "config_path", "", "path to config file")
	flag.Parse()
	return path
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("failed to read config:" + err.Error())
	}

	return &cfg
}
