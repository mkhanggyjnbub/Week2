package Query

import (
	"baitapweek2/Db"
	"errors"
)

func CheckLogin(email string, password string, role string) (*Db.Users, error) {
	var userQ Db.Users
	Result := Db.DB.Where("email = ?  and  password_hash  = ? and role = ?", email, password, role).First(&userQ)
	if Result.Error != nil {
		return nil, errors.New("lỗi tài khoản và mật khẩu")
	}
	return &userQ, nil
}
func InsertUser(userQ Db.Users) (int64, error) {
	Result := Db.DB.Create(&userQ)
	if Result.Error != nil {
		return 0, errors.New("lỗi tài khoản và mật khẩu")
	}
	return Result.RowsAffected, nil
}
