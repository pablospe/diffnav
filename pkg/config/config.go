package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HideHeader      bool   `toml:"hide_header"`
	HideFooter      bool   `toml:"hide_footer"`
	ShowFileTree    bool   `toml:"show_file_tree"`
	FileTreeWidth   int    `toml:"file_tree_width"`
	SearchTreeWidth int    `toml:"search_tree_width"`
	Icons           string `toml:"icons"` // "nerd-fonts", "unicode" (default), "ascii"
}

func DefaultConfig() Config {
	return Config{
		HideHeader:      false,
		HideFooter:      false,
		ShowFileTree:    true,
		FileTreeWidth:   26,
		SearchTreeWidth: 50,
		Icons:           "ascii",
	}
}

func getConfigFilePath() string {
	var configDirs []string

	// Environment variable override - useful for development or non-standard setups.
	if dir := os.Getenv("DIFFNAV_CONFIG_DIR"); dir != "" {
		if s, err := os.Stat(dir); err == nil && s.IsDir() {
			return filepath.Join(dir, "config.toml")
		}
	}

	// On macOS, check ~/.config first (common for CLI tools),
	// then XDG_CONFIG_HOME if set.
	// os.UserConfigDir() already handles this for Linux.
	if runtime.GOOS == "darwin" {
		if home := os.Getenv("HOME"); home != "" {
			configDirs = append(configDirs, filepath.Join(home, ".config"))
		}
		if xdgConfigDir := os.Getenv("XDG_CONFIG_HOME"); xdgConfigDir != "" {
			configDirs = append(configDirs, xdgConfigDir)
		}
	}

	// Standard OS-specific config directory.
	if configDir, err := os.UserConfigDir(); err == nil {
		configDirs = append(configDirs, configDir)
	}

	// Return the first config file that exists.
	for _, dir := range configDirs {
		configPath := filepath.Join(dir, "diffnav", "config.toml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// If no config file exists, return the preferred path for creation.
	if len(configDirs) > 0 {
		return filepath.Join(configDirs[0], "diffnav", "config.toml")
	}
	return ""
}

func Load() Config {
	cfg := DefaultConfig()

	configPath := getConfigFilePath()
	if configPath == "" {
		return cfg
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg
	}

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return DefaultConfig()
	}

	return cfg
}
