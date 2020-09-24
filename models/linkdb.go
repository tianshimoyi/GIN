package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	DB, err = gorm.Open("mysql", "root:lclyzw17@(127.0.0.1:3306)/CommoditySalesSystem?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("连接失败")
	} else {
		fmt.Println("连接成功")
	}
	DB.AutoMigrate(Commodity{})
	DB.AutoMigrate(GuestOrAdministrators{})
	DB.AutoMigrate(Order{})
	DB.AutoMigrate(Session{})
	DB.AutoMigrate(CateGory{})
	DB.AutoMigrate(Addr{})
	//DB.Unscoped().Delete(Order{}, "id between ? and ?", 21, 22)
	//com:= Commodity{
	//	ID:          1,
	//	Name:        "手电",
	//	Img:         "/static/img/手电.jpg",
	//	Stock:       &sql.NullInt64{1005,true},
	//	CateGoryID:    1,
	//	Price:       124,
	//	SalesVolume: &sql.NullInt64{50000,true},
	//	Details:"这是一款超级好用的电灯。",
	//}
	//DB.Create(&com)
	//cg:=CateGory{
	//	ID:        2,
	//	Name:      "小家具",
	//	Introduce: "各种方便使用的小家具",
	//}
	//DB.Create(&cg)
	//user:=GuestOrAdministrators{
	//	Name:              "张雨",
	//	Password:          "123456",
	//	Img:               "/static/img/demacia-silvermere.jpg",
	//	Gender:            "女",
	//	Tel:               "18841299421",
	//	IDNum:             "411322199809093121",
	//	ConsumptionAmount: 0,
	//	IsAdministrators:  2,
	//}
	//DB.Create(&user)
}
