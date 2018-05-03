//package models
//
//import (
//	"github.com/jinzhu/gorm"
//	_ "github.com/jinzhu/gorm/dialects/mysql"
//	"gin-docker-mysql/pkg/setting"
//	"fmt"
//	"log"
//	"time"
//)
//
//type BaseModel struct {
//	ID        uint `gorm:"primary_key"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}
//type Post struct {
//	BaseModel
//	Title       string                // title
//	Body        string                // body
//	//View        int                   // view count
//	//IsPublished bool                  // published or not
//	Tags        []*Tag     `gorm:"-"` // tags of post
//	Comments    []*Comment `gorm:"-"` // comments of post
//}
//type Tag struct {
//	BaseModel
//	Name  string         // tag name
//	//Total int `gorm:"-"` // count of post
//}
//
//// table post_tags
//type PostTag struct {
//	BaseModel
//	PostId uint // post id
//	TagId  uint // tag id
//}
//type Comment struct {
//	BaseModel
//	//UserID    uint                      // 用户id
//	Content   string                    // 内容
//	PostID    uint                      // 文章id
//	//ReadState bool `gorm:"default:'0'"` // 阅读状态
//	//Replies []*Comment // 评论
//	//NickName  string `gorm:"-"`
//	//AvatarUrl string `gorm:"-"`
//	//GithubUrl string `gorm:"-"`
//}
//
//var DB *gorm.DB
//
//func InitDB() (*gorm.DB, error) {
//
//	//db.CreateTable(&Product{})
//	sec, err := setting.Cfg.GetSection("database")
//	dbType := sec.Key("TYPE").String()
//	dbName := sec.Key("NAME").String()
//	user := sec.Key("USER").String()
//	password := sec.Key("PASSWORD").String()
//	host := sec.Key("HOST").String()
//	DB, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
//		user,
//		password,
//		host,
//		dbName))
//	if err != nil {
//		log.Println(err)
//	}
//	DB.DB().SetMaxIdleConns(10)
//	DB.DB().SetMaxOpenConns(100)
//
//	DB.AutoMigrate(&Post{},&Tag{},&PostTag{},&Comment{})
//
//	return nil, err
//}
//
//func CloseDB() {
//	defer DB.Close()
//}

package models

import (
	"log"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gin-docker-mysql/pkg/setting"

	"time"
	//"gin-docker-mysql/pkg/logging"
)

var DB *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedAt *time.Time
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host  string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		fmt.Println("error")
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	//tablePrefix = sec.Key("TABLE_PREFIX").String()

	DB, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}
	fmt.Println("数据库连接")
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return tablePrefix + defaultTableName
	//}
	DB.SingularTable(true)

	DB.AutoMigrate(&Auth{}, &Article{}, &Tag{},&Comment{})


	DB.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DB.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer DB.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
