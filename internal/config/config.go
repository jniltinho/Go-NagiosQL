// Package config loads and validates application configuration from config.toml
// and NAGIOSQL_* environment variable overrides via Viper.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config is the top-level configuration structure.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
	Nagios   NagiosConfig   `mapstructure:"nagios"`
}

// ServerConfig controls the HTTP server behaviour.
type ServerConfig struct {
	Port int  `mapstructure:"port"`
	Dev  bool `mapstructure:"dev"`
}

// JWTConfig holds token signing parameters.
type JWTConfig struct {
	Secret         string `mapstructure:"secret"`
	AccessTTLMin   int    `mapstructure:"access_ttl_min"`
	RefreshTTLDays int    `mapstructure:"refresh_ttl_days"`
}

// DatabaseConfig holds MariaDB connection parameters.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// NagiosConfig holds all filesystem paths used when generating or verifying
// Nagios configuration files. These values are seeded into tbl_configtarget on
// first migrate and can be updated via the admin API later.
type NagiosConfig struct {
	BaseDir          string `mapstructure:"base_dir"`
	ConfigFile       string `mapstructure:"config_file"`
	CgiFile          string `mapstructure:"cgi_file"`
	ResourceFile     string `mapstructure:"resource_file"`
	ReloadTrigger    string `mapstructure:"reload_trigger"`
	Binary           string `mapstructure:"binary"`
	PidFile          string `mapstructure:"pid_file"`
	HostConfigDir    string `mapstructure:"host_config_dir"`
	ServiceConfigDir string `mapstructure:"service_config_dir"`
	BackupDir        string `mapstructure:"backup_dir"`
	ImportDir        string `mapstructure:"import_dir"`
}

// Load reads the configuration file at cfgFile and applies any NAGIOSQL_*
// environment variable overrides. All NAGIOSQL_ env vars map to nested keys
// using underscores as separators (e.g. NAGIOSQL_DATABASE_PASSWORD).
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("server.port", 8081)
	v.SetDefault("server.dev", false)
	v.SetDefault("jwt.access_ttl_min", 15)
	v.SetDefault("jwt.refresh_ttl_days", 7)
	v.SetDefault("database.host", "db")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.name", "nagiosql")
	v.SetDefault("database.user", "nagiosql")
	v.SetDefault("nagios.base_dir", "/usr/local/nagios")
	v.SetDefault("nagios.config_file", "/usr/local/nagios/etc/nagios.cfg")
	v.SetDefault("nagios.cgi_file", "/usr/local/nagios/etc/cgi.cfg")
	v.SetDefault("nagios.resource_file", "/usr/local/nagios/etc/resource.cfg")
	v.SetDefault("nagios.reload_trigger", "/usr/local/nagios/var/reload.trigger")
	v.SetDefault("nagios.binary", "/usr/local/nagios/bin/nagios")
	v.SetDefault("nagios.pid_file", "/usr/local/nagios/var/nagios.lock")
	v.SetDefault("nagios.host_config_dir", "/usr/local/nagios/etc/nagiosql/hosts/")
	v.SetDefault("nagios.service_config_dir", "/usr/local/nagios/etc/nagiosql/services/")
	v.SetDefault("nagios.backup_dir", "/usr/local/nagios/etc/nagiosql/backup/")
	v.SetDefault("nagios.import_dir", "/usr/local/nagios/etc/import/")

	// Config file
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("toml")
		v.AddConfigPath(".")
	}

	// Environment variable overrides: NAGIOSQL_DATABASE_PASSWORD → database.password
	v.SetEnvPrefix("NAGIOSQL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
		// Config file is optional; defaults and env vars are sufficient.
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("jwt.secret must be set in config.toml or NAGIOSQL_JWT_SECRET")
	}
	if len(cfg.JWT.Secret) < 32 {
		return nil, fmt.Errorf("jwt.secret must be at least 32 characters")
	}

	return &cfg, nil
}

// DSN builds a GORM-compatible MySQL DSN string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}
