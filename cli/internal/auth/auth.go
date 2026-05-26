package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/power721/115driver/pkg/driver"
	"github.com/spf13/viper"
)

const (
	DefaultConfigDir  = ".115driver"
	DefaultConfigFile = "config.toml"
	DefaultProfile    = "main"
	EnvCookie         = "DRIVER115_COOKIE"
	EnvConfig         = "DRIVER115_CONFIG"
	EnvProfile        = "DRIVER115_PROFILE"
	EnvDebug          = "DRIVER115_DEBUG"
)

type Config struct {
	DefaultProfile string             `mapstructure:"default_profile"`
	Profiles       map[string]Profile `mapstructure:"profiles"`
}

type Profile struct {
	Cookie string `mapstructure:"cookie"`
}

func ResolveCredential(cookieFlag, configPath, profile string) (*driver.Credential, error) {
	if cookieFlag != "" {
		cr := &driver.Credential{}
		if err := cr.FromCookie(cookieFlag); err != nil {
			return nil, fmt.Errorf("invalid cookie: %w", err)
		}
		return cr, nil
	}

	if envCookie := os.Getenv(EnvCookie); envCookie != "" {
		cr := &driver.Credential{}
		if err := cr.FromCookie(envCookie); err != nil {
			return nil, fmt.Errorf("invalid cookie from env: %w", err)
		}
		return cr, nil
	}

	path := configPath
	if path == "" {
		if envPath := os.Getenv(EnvConfig); envPath != "" {
			path = envPath
		} else {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, DefaultConfigDir, DefaultConfigFile)
		}
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no credential found. Run '115driver login' to authenticate")
		}
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("no credential found. Run '115driver login' to authenticate")
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if profile == "" {
		if envProfile := os.Getenv(EnvProfile); envProfile != "" {
			profile = envProfile
		}
	}
	if profile == "" {
		profile = v.GetString("default_profile")
	}
	if profile == "" {
		profile = DefaultProfile
	}

	cookieStr := v.GetString("profiles." + profile + ".cookie")
	if cookieStr == "" {
		return nil, fmt.Errorf("no credential found for profile '%s'. Run '115driver login'", profile)
	}

	cr := &driver.Credential{}
	if err := cr.FromCookie(cookieStr); err != nil {
		return nil, fmt.Errorf("invalid cookie in config: %w", err)
	}
	return cr, nil
}

func ResolveProfile(profile string) string {
	if profile != "" {
		return profile
	}
	if envProfile := os.Getenv(EnvProfile); envProfile != "" {
		return envProfile
	}
	return DefaultProfile
}

func SaveCredential(configPath, profile, cookie string) error {
	path := configPath
	if path == "" {
		if envPath := os.Getenv(EnvConfig); envPath != "" {
			path = envPath
		} else {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, DefaultConfigDir, DefaultConfigFile)
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	profile = ResolveProfile(profile)

	v.Set("default_profile", profile)
	v.Set("profiles."+profile+".cookie", cookie)

	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return os.Chmod(path, 0600)
}
