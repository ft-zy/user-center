package service

import (
	log "github.com/sirupsen/logrus"
	"my-go-user-center/src/config"
	"my-go-user-center/src/model"
)

func CreateTag(req *model.YpTags) bool {
	err := config.Db.Model(&model.YpTags{}).Create(&req).Error
	return err == nil
}

func UpdateTag(req *model.YpTags) bool {
	err := config.Db.Model(&model.YpTags{}).Updates(&req).Error
	return err == nil
}

func DeleteTag(req *model.YpTags) bool {
	var tag model.YpTags
	result := config.Db.Model(&model.YpTags{}).Find(&tag).Where("id = ? and isDelete = 0", req.ID)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Warning("节点已不存在")
		return false
	}
	result = config.Db.Model(&model.YpTags{}).Find(&tag).Where("id = ? and isDelete = 0", req.ParentID)
	if result.Error == nil && result.RowsAffected > 0 {
		log.Warning("父节点存在，无法删除")
		return false
	}
	config.Db.Model(&model.User{}).Where("id=?", req.ID).Update("isDelete", 1)
	return true
}

func ListTag() []model.YpTags {
	var list []model.YpTags
	config.Db.Model(&model.User{}).Find(&list)
	return list
}
