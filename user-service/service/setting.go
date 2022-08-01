package service

import "github.com/OJ-lab/oj-lab-services/utils"

var serviceSettings utils.ServiceSettings

func SetupServiceSetting(settings utils.ServiceSettings) {
	serviceSettings = settings
}
