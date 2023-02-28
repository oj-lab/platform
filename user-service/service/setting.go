package service

import "github.com/OJ-lab/oj-lab-services/config"

var serviceSettings config.ServiceSettings

func SetupServiceSetting(settings config.ServiceSettings) {
	serviceSettings = settings
}
