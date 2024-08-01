package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

const configDir = ".mycli"

func SaveConfig(config DatabaseConfig) error {
	configPath := filepath.Join(os.Getenv("HOME"), configDir)
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	viper.Reset()
	viper.SetConfigName(config.Name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.Set("host", config.Host)
	viper.Set("port", config.Port)
	viper.Set("user", config.User)
	viper.Set("password", config.Password)
	viper.Set("dbname", config.DBName)

	configFile := filepath.Join(configPath, config.Name+".yaml")
	return viper.WriteConfigAs(configFile)
}

func LoadConfig(name string) (DatabaseConfig, error) {
	configPath := filepath.Join(os.Getenv("HOME"), configDir)

	viper.Reset()
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return DatabaseConfig{}, fmt.Errorf("failed to read config: %v", err)
	}

	return DatabaseConfig{
		Name:     name,
		Host:     viper.GetString("host"),
		Port:     viper.GetInt("port"),
		User:     viper.GetString("user"),
		Password: viper.GetString("password"),
		DBName:   viper.GetString("dbname"),
	}, nil
}

func ListConfigs() ([]string, error) {
	configPath := filepath.Join(os.Getenv("HOME"), configDir)
	var configs []string

	err := filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			configs = append(configs, filepath.Base(path[:len(path)-len(filepath.Ext(path))]))
		}
		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to list configs: %v", err)
	}

	return configs, nil
}
