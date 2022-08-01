package user_tests

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
)

func SetupTestDatabase() {
	dataBaseSettings, err := utils.GetDatabaseSettings("../../../config/test.ini")
	if err != nil {
		panic("failed to get database settings")
	}
	model.OpenConnection(dataBaseSettings)
}
