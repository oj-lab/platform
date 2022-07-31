package tests

import (
	sharedTools "github.com/OJ-lab/oj-lab-services/shared-tools"
	"log"
	"testing"
)

func TestIniBasicUsage(t *testing.T) {
	databaseSettings := sharedTools.GetDatabaseSettings("../../config/example.ini")
	log.Print(sharedTools.GetDatabaseDSN(databaseSettings))
}
