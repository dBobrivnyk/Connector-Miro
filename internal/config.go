package config

import "github.com/spf13/viper"

const defaultConfigPath = "."

type Config struct {
	IdStoragePath string
	SyncPeriod    int64
	BufferSize    uint64
}

func Load() (*Config, error) {
	viper.AddConfigPath(defaultConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		IdStoragePath: viper.GetString("idStoragePath"),
		SyncPeriod:    viper.GetInt64("syncPeriod"),
		BufferSize:    viper.GetUint64("bufferSize"),
	}

	return cfg, nil
}
