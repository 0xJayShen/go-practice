package models

import "fmt"

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func ExistTagByName(name string) bool {
	var tag Tag
	DB.Select("id").Where("name = ?", name, 0).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func AddTag(name string, state int, createdBy string) bool {
	DB.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})

	return true
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
	DB.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistTagByID(id int) bool {
	var tag Tag
	DB.Select("id").Where("id = ? ", id,).First(&tag)
	fmt.Println(tag.ID)
	if tag.ID > 0 {
		return true
	}

	return false
}

func DeleteTag(id int) bool {
	DB.Where("id = ?", id).Delete(&Tag{})

	return true
}

func EditTag(id int, data interface{}) bool {
	DB.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
}
