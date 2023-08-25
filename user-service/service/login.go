package service

import (
	"time"

	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/OJ-lab/oj-lab-services/packages/utils"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
	"github.com/gin-gonic/gin"
)

const cookieAgeProp = "service.cookie_age"

var cookieAge string

func init() {
	cookieAge = application.AppConfig.GetString(cookieAgeProp)
}

func Login(c *gin.Context) {
	account := c.Param("account")
	password := c.PostForm("password")
	result, err := business.ComparePassword(account, password)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	if result {
		userInfo, err := business.GetUserInfo(&account, nil, nil)
		if err != nil {
			utils.ApplyError(c, err)
			return
		}
		token, err := business.GenerateTokenString(userInfo.Account, []model.Role{model.RoleUser})
		if err != nil {
			utils.ApplyError(c, err)
			return
		}
		duration, err := time.ParseDuration(cookieAge)
		if err != nil {
			utils.ApplyError(c, err)
			return
		}
		c.SetCookie("jwt", token, int(duration.Seconds()), "", "", false, true)
		c.JSON(200, gin.H{
			"status": "success",
		})
	} else {
		c.JSON(403, gin.H{
			"status": "failed",
		})
	}
}
