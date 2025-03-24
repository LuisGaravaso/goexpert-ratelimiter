package configs

import (
	"fmt"
	"log"
	"path/filepath"
	"ratelim/internal/api/web/middleware/ratelimiter/entity"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*entity.Config, error) {
	v := viper.New()

	// Detect if user passed a full path or just a filename
	if strings.Contains(configPath, string(filepath.Separator)) {
		dir := filepath.Dir(configPath)
		file := filepath.Base(configPath)
		ext := filepath.Ext(file)
		name := strings.TrimSuffix(file, ext)

		v.SetConfigName(name)
		v.SetConfigType(strings.TrimPrefix(ext, "."))
		v.AddConfigPath(dir)
	} else {
		_, currentFile, _, _ := runtime.Caller(0)
		baseDir := filepath.Dir(currentFile)

		ext := filepath.Ext(configPath)
		name := strings.TrimSuffix(configPath, ext)

		v.SetConfigName(name)
		v.SetConfigType(strings.TrimPrefix(ext, "."))

		v.AddConfigPath(baseDir)
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg entity.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate the config
	errs := cfg.Validate()
	if len(cfg.Services) == 0 {
		return nil, fmt.Errorf("error validating config: %w", errs[0])
	}
	for _, err := range errs {
		log.Printf("%s", err.Error())
	}

	return &cfg, nil
}
