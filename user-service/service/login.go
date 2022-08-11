package service

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func Login(c *gin.Context) {
	account := c.Param("account")
	password := c.PostForm("password")
	result, err := model.ComparePassword(account, password)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	if result {
		userInfo, err := model.GetUserInfo(&account, nil, nil)
		if err != nil {
			utils.ApplyError(c, err)
			return
		}
		token, err := utils.GenerateTokenString(userInfo.Account, userInfo.Role)
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
