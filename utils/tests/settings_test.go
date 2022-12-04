package tests

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/utils"
)

func TestIniBasicUsage(t *testing.T) {
	databaseSettings, _ := config.GetDatabaseSettings("../../config/ini/example.ini")
	log.Print(utils.GetDatabaseDSN(databaseSettings))
}
