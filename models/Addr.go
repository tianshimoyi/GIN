package models

import "github.com/jinzhu/gorm"

type Addr struct {
	gorm.Model
	Country         string `gorm:"size:30" form:"AddrCoun"`   //国家
	Province        string `gorm:"size:30" form:"AddrPro"`    //省份
	City            string `gorm:"size:30" form:"AddrCity"`   //市
	County          string `gorm:"size:30" form:"AddrCy"`     //县
	DetailedAddress string `gorm:"size:150" form:"AddrDeAdd"` //详情
	GuestID         uint   `form:"AddrGuID"`                  //用户ID
}

//按用户ID查询地址
func SearchAddrByGID(id string) []Addr {
	var addrs []Addr
	DB.Find(&addrs, "guest_id=?", id)
	return addrs
}

//删除地址
func DelAddr(id string) {
	var sum int
	DB.Model(Order{}).Where("addr_id=?", id).Count(&sum)
	if sum <= 0 {
		DB.Unscoped().Delete(Addr{}, "id=?", id)
	}
}

//按ID查询地址
func SearchAddr(addrID uint) Addr {
	var addr Addr
	DB.First(&addr, "id=?", addrID)
	return addr
}
