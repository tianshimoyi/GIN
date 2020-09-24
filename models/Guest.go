package models

import "github.com/jinzhu/gorm"

//客人或者管理员基本信息
type GuestOrAdministrators struct {
	gorm.Model
	Name              string  `gorm:"unique;size:20;not null" form:"GAName"`  //姓名
	Password          string  `gorm:"not null" form:"GAPassword"`             //密码
	Img               string  `gorm:"not null"`                               //照片
	Gender            string  `gorm:"size:3;not null" form:"GAGender"`        //性别
	Tel               string  `gorm:"not null;unique;size:11" form:"GATel"`   //手机号
	IDNum             string  `gorm:"not null;unique;size:20" form:"GAIDNum"` //身份证
	ConsumptionAmount float64 `form:"GAConAmount"`                            //消费总金额
	AddrID            uint    `form:"GAAddr"`                                 //地址
	IsAdministrators  int64   `gorm:"default:1;size:1" form:"isAd"`           //是否是管理员1代表不是2代表是
}

//查找一个用户或管理员
func SeachOneGA(user *GuestOrAdministrators) *GuestOrAdministrators {
	var u GuestOrAdministrators
	DB.Where(user).First(&u)
	return &u
}

//插入
func InsertGA(user *GuestOrAdministrators) error {
	return DB.Create(user).Error
}

//根据id查询
func SearchUserByID(id string) GuestOrAdministrators {
	var u GuestOrAdministrators
	DB.First(&u, "id=?", id)
	return u
}

//折扣
func Discount(user *GuestOrAdministrators) float64 {
	switch {
	case user.ConsumptionAmount >= 10000:
		return 0.2
	case user.ConsumptionAmount >= 8000:
		return 0.15
	case user.ConsumptionAmount >= 6000:
		return 0.1
	case user.ConsumptionAmount >= 4000:
		return 0.05
	default:
		return 0
	}
}
