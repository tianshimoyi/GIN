package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

//订单
type Order struct {
	gorm.Model
	GuestID       uint  `gorm:"not null" form:"OrGuestID"` //客户ID
	ComID         int64 `gorm:"not null" form:"OrComID"`   //商品ID
	CGID          int64 //商品类别ID
	Tel           string
	ComSum        int64   `gorm:"not null" form:"OrderSum"` //商品数量
	Price         float64 `gorm:"not null" form:"OrPrice"`  //价钱
	AddrID        uint    `form:"OrAddrID"`                 //地址ID
	ArriveTime    string  `form:"OrArrTime"`                //到达时间
	IsHandel      int64   //是否处理 1 未处理 2待处理 3处理 4处理完成
	DeliveryOrNot bool    //是否发货
	Persent       float64 //销量百分比
}

//查询订单
func SearchOrder(i string, uid uint) ([]Order, error) {
	var orders []Order
	err := DB.Unscoped().Where("guest_id=? and is_handel=?", uid, i).Find(&orders).Error
	return orders, err
}

//查询所有待处理订单
func SearchALLOrder(i string) []Order {
	var orders []Order
	DB.Where("is_handel=?", i).Find(&orders)
	return orders
}

func SearchOneOrder(id string) Order {
	var order Order
	DB.First(&order, "id=?", id)
	return order
}

//删除订单
func DelteOrder(orderID string) {
	DB.Unscoped().Delete(Order{}, "id=?", orderID)
}

//将字符串转换为时间
func StrToTime(str string) (time.Time, error) {
	local, _ := time.LoadLocation("Asia/Shanghai")
	timeObj, err := time.ParseInLocation("2006-01-02", str, local)
	return timeObj, err
}

func NoRepead(orders []Order) []Order {
	if len(orders) == 0 {
		return nil
	}
	var orders2 = []Order{}
	orders2 = append(orders2, orders[0])
	for i := 1; i < len(orders); i++ {
		var flag = true
		for j := 0; j < len(orders2); j++ {
			if orders[i].ComID == orders2[j].ComID {
				orders2[j].ComSum += orders[i].ComSum
				orders2[j].Price += orders[i].Price
				flag = false
				continue
			}
		}
		if flag == true {
			orders2 = append(orders2, orders[i])
		}
	}
	return orders2
}
