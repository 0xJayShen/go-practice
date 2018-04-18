//package models
//
//import (
//	"github.com/jinzhu/gorm"
//	_ "github.com/jinzhu/gorm/dialects/mysql"
//	"Seckill/pkg/setting"
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
	"Seckill/pkg/setting"

	"time"
	"Seckill/pkg/logging"

)

var db *gorm.DB

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedOn int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}
func init() {

	var (
		err error
		dbType, dbName, user, password, host string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	fmt.Println(user,password,host)
	//tablePrefix = sec.Key("TABLE_PREFIX").String()

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))
	fmt.Println(11111)
	if err != nil{
		logging.Error(err)
	}
	db.SingularTable(true)

	db.AutoMigrate(&Tag{},&Auth{},&Article{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tag{},&Auth{},&Article{})

	if err != nil {
		logging.Info(err)
	}

	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//	return tablePrefix + defaultTableName;
	//}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}
