package service

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	account := c.Param("account")

	if serviceSettings.AuthOn {
		tokenString := c.GetHeader("Authorization")
		selfAccount, role, err := utils.ParseTokenString(tokenString)
		if err != nil {
			utils.ApplyError(c, err)
			return
		}
		if selfAccount != account && model.String2Role(role) != model.RoleAdmin {
			c.JSON(401, gin.H{
				"status": "method not allowed",
			})
			return
		}
	}

	err := model.DeleteUser(account)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
	})
}
