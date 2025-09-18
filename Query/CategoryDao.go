package Query

import (
	"baitapweek2/Db"
	"fmt"
)

func GetCategory() ([]Db.Category, error) {
	var category []Db.Category

	result := Db.DB.Find(&category)
	if result.Error != nil {
		fmt.Println("Lỗi không thể lấy category")
		return nil, result.Error
	}
	return category, nil
}
