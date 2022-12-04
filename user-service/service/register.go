package service

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	roleStringArray := c.PostFormArray("roles")
	roles := []model.Role{}
	for _, roleString := range roleStringArray {
		roles = append(roles, model.String2Role(roleString))
	}
	err := business.CreateUser(account, password, roles)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
	})
}
