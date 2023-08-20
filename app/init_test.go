package app

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestInit(T *testing.T) {
	err := LoadConfig("../config")
	if err != nil {
		T.Fatal(err)
	}

	databaseType := viper.GetString("database.type")
	databaseUser := viper.GetString("database.user")
	println(databaseType, databaseUser)
	if len(databaseType) == 0 {
		T.Fatal("databaseType not loaded")
	}

	SetupLog()

	logrus.Error("debug")
}
