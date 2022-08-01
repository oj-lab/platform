package service

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserInfo(c *gin.Context) {
	account := c.Param("account")
	userInfo, err := model.GetUserInfo(&account, nil, nil)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"userInfo": userInfo,
	})
}

func FindUserInfos(c *gin.Context) {
	account := c.Query("account")
	account = "%" + account + "%"
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	userInfos, err := model.FindUserInfos(account, offset, limit)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"userInfos": userInfos,
	})
}
