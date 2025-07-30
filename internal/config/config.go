package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Listen    string            `mapstructure:"listen"`
	JWTSecret string            `mapstructure:"jwt_secret"`
	LogLevel  string            `mapstructure:"log_level"`
	Policies  map[string]string `mapstructure:"policies"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/mini-mcp/")
	if path != "" {
		v.SetConfigFile(path)
	}
	v.AutomaticEnv()
	v.SetEnvPrefix("MINIMCP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	if err := validate(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func validate(c *Config) error {
	if c.Listen == "" {
		return fmt.Errorf("listen is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("jwt_secret is required")
	}
	if len(c.Policies) == 0 {
		return fmt.Errorf("at least one policy is required")
	}
	for team, backend := range c.Policies {
		if team == "" || backend == "" {
			return fmt.Errorf("invalid policy: team and backend required")
		}
	}
	return nil
}
