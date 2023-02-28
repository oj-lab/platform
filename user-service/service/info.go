package service

import (
	"strconv"

	"github.com/OJ-lab/oj-lab-services/user-service/business"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	account, _, err := business.ParseTokenString(token)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	userInfo, err := business.GetUserInfo(&account, nil, nil)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"userInfo": userInfo,
	})
}

func GetUserInfo(c *gin.Context) {
	account := c.Param("account")
	userInfo, err := business.GetUserInfo(&account, nil, nil)
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
	userInfos, err := business.FindUserInfos(account, offset, limit)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	total, err := business.CountUser(account)
	if err != nil {
		utils.ApplyError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"userInfos": userInfos,
		"total":     total,
	})
}
