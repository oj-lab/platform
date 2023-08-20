package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestViper(T *testing.T) {
	err := loadConfig("./")
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
