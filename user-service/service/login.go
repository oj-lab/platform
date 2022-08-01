package service

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
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
		c.JSON(200, gin.H{
			"token": token,
		})
	} else {
		c.JSON(403, gin.H{
			"status": "failed",
		})
	}
}
