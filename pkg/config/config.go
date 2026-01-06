package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HideHeader      bool `toml:"hide_header"`
	HideFooter      bool `toml:"hide_footer"`
	ShowFileTree    bool `toml:"show_file_tree"`
	FileTreeWidth   int  `toml:"file_tree_width"`
	SearchTreeWidth int  `toml:"search_tree_width"`
}

func DefaultConfig() Config {
	return Config{
		HideHeader:      false,
		HideFooter:      false,
		ShowFileTree:    true,
		FileTreeWidth:   26,
		SearchTreeWidth: 50,
	}
}

func Load() Config {
	cfg := DefaultConfig()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return cfg
	}

	configPath := filepath.Join(configDir, "diffnav", "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg
	}

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return DefaultConfig()
	}

	return cfg
}
