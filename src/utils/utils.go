package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"my-go-user-center/src/constant"
	"my-go-user-center/src/model"
	"reflect"
)

// testfdsfdsfsdf
// md5 的加密算法
func EncryptMd5(userPassword string) string {
	h := md5.New()
	h.Write([]byte(constant.SALT + userPassword))
	return hex.EncodeToString(h.Sum(nil))
}

// 用户信息脱敏
func GetSafetyUser(user *model.User) model.User {
	safetyUser := model.User{}
	safetyUser.Id = user.Id
	safetyUser.Username = user.Username
	safetyUser.UserAccount = user.UserAccount
	safetyUser.AvatarUrl = user.AvatarUrl
	safetyUser.Gender = user.Gender
	safetyUser.Phone = user.Phone
	safetyUser.Email = user.Email
	safetyUser.UserStatus = user.UserStatus
	safetyUser.UserRole = user.UserRole
	safetyUser.CreateTime = user.CreateTime
	return safetyUser
}

// CopyStructFields copies fields from src to dst if they have the same name and type.
// It only copies exported fields (those that start with an uppercase letter).
func CopyStructFields(src, dst interface{}) error {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	// Ensure both are structs
	if dstVal.Kind() != reflect.Struct || srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("CopyStructFields expects both arguments to be structs")
	}

	// Iterate over fields in the source struct
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldName := srcVal.Type().Field(i).Name

		// Get the corresponding field in the destination struct
		dstField := dstVal.FieldByName(srcFieldName)
		if !dstField.IsValid() || !dstField.CanSet() {
			// Field doesn't exist or is not exportable/settable
			continue
		}

		// Check if types are assignable
		if !dstField.Type().AssignableTo(srcField.Type()) {
			// Types are not compatible
			continue
		}

		// Perform the copy
		dstField.Set(srcField)
	}

	return nil
}
