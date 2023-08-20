package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestViper(T *testing.T) {
	viper.SetConfigFile("./test.toml")
	err := viper.ReadInConfig()
	if err != nil {
		T.Fatal(err)
	}

	viper.SetConfigFile("./test_override.toml")
	err = viper.MergeInConfig()
	if err != nil {
		T.Fatal(err)
	}

	databaseType := viper.GetString("database.type")
	databaseUser := viper.GetString("database.user")
	println(databaseType, databaseUser)
	if len(databaseType) == 0 {
		T.Fatal("databaseType not loaded")
	}
}
