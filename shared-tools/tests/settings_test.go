package tests

import (
	sharedtools "github.com/OJ-lab/oj-lab-services/shared-tools"
	"testing"
)

func TestIniBasicUsage(t *testing.T) {
	sharedtools.GetDatabaseSettings("../config/default.ini")
}
