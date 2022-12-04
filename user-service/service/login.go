package service

import (
	"time"

	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
)

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
		duration, err := time.ParseDuration(serviceSettings.CookieAge)
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
