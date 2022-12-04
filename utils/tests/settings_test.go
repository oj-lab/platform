package tests

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/utils"
)

func TestIniBasicUsage(t *testing.T) {
	databaseSettings, err := config.GetDatabaseSettings("../../config/ini/test.ini")
	if err != nil {
		t.Error("Fail to load DB settings: ", err)
	}
	log.Print(utils.MustGetDatabaseDSN(*databaseSettings))
}
