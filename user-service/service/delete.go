package service

import (
	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
	"github.com/OJ-lab/oj-lab-services/packages/utils"
	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	account := c.Param("account")

	if serviceSettings.AuthOn {
		tokenString := c.GetHeader("Authorization")
		selfAccount, roles, err := business.ParseTokenString(tokenString)
		if err != nil {
			utils.ApplyError(c, err)
			return
		}
		if selfAccount != account && model.RoleInRoles(model.RoleAdmin, model.Array2Roles(roles)) {
			c.JSON(401, gin.H{
				"status": "method not allowed",
			})
			return
		}
	}

	err := business.DeleteUser(account)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
	})
}
