package models

type Comment struct {
	Model

	ArticleID int `json:"Article_id" gorm:"index"`
	Article   Tag `json:"article"`

	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}
func GetComments(maps interface{}) (comments []Comment) {
	DB.Preload("Article").Where(maps).Find(&comments)

	return
}