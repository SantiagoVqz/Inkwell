package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	VaultPath string `yaml:"vault_path"`
	VaultSubfolder string `yaml:"vault_subfolder"`
	DatabasePath string `yaml:"database_path"`
	FetchConcurrency int `yaml:"fetch_concurrency"`
	Log LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	Format string `yaml:"format"`
}

func Default() Config {
	home, _ := os.UserHomeDir()

	return Config {
		VaultSubfolder: "Inbox/Inkwell",
		DatabasePath: filepath.Join(home, ".local", "share", "inkwell", "inkwell.db"),
		FetchConcurrency: 8,
		Log: LogConfig{
			Level: "info",
			Format: "auto",
		},
	}
}

func DefaultPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "inkwell", "config.yaml")
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "inkwell", "config.yaml")
}