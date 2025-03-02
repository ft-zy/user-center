package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"my-go-user-center/src/common"
	"my-go-user-center/src/model"
	"my-go-user-center/src/service"
)

// 添加标签
func CreateTag(c *gin.Context) {
	req := &model.YpTags{}
	// 从前端接收请求
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("CreateTag request json err %v", err)
		common.RespFail(c.Writer, err.Error())
		return
	}
	common.RespOK(c.Writer, service.CreateTag(req), "成功")
}

// 修改标签
func UpdateTag(c *gin.Context) {
	req := &model.YpTags{}
	// 从前端接收请求
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("CreateTag request json err %v", err)
		common.RespFail(c.Writer, err.Error())
		return
	}
	common.RespOK(c.Writer, service.UpdateTag(req), "成功")
}

// 删除标签
func DeleteTag(c *gin.Context) {
	req := &model.YpTags{}
	// 从前端接收请求
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("CreateTag request json err %v", err)
		common.RespFail(c.Writer, err.Error())
		return
	}
	common.RespOK(c.Writer, service.DeleteTag(req), "成功")
}

// 查标签
func ListTag(c *gin.Context) {

	common.RespOK(c.Writer, service.ListTag(), "成功")
}
