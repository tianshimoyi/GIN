package models

import (
	"database/sql"
	"fmt"
	"math"
)

//商品信息
type Commodity struct {
	ID                int64          `gorm:"primary_key" form:"comID"`            //ID并不是自增
	Name              string         `grom:"not null;size:30" form:"comName"`     //名字
	Img               string         `form:"not null;unique" form:"comImg"`       //照片
	Stock             *sql.NullInt64 `gorm:"deafult:0" form:"comStock"`           //库存
	CateGoryID        int64          `gorm:"size:10;not null" form:"comCateGory"` //类别
	Price             float64        `gorm:"not null" form:"comPrice"`            //价格
	Details           string         `gorm:"type:blob" form:"comDetails"`         //详情介绍
	BriefIntroduction string         //简介
	SalesVolume       *sql.NullInt64 `gorm:"default:0" form:"comSaVm"` //销量
	//UserID uint
}

//分页查询商品
func SearchCommodity(id, num, orderID int, name ...string) []Commodity {
	var coms []Commodity
	if len(name) == 0 {
		switch orderID {
		case 1:
			DB.Order("price asc").Offset((id - 1) * num).Limit(num).Find(&coms)
		case 2:
			DB.Order("price desc").Offset((id - 1) * num).Limit(num).Find(&coms)
		case 3:
			DB.Order("sales_volume desc").Offset((id - 1) * num).Limit(num).Find(&coms)
		default:
			DB.Offset((id - 1) * num).Limit(num).Find(&coms)
		}
	} else {
		switch orderID {
		case 1:
			DB.Where("name regexp ?", "["+name[0]+"]").Order("price asc").Offset((id - 1) * num).Limit(num).Find(&coms)
		case 2:
			DB.Where("name regexp ?", "["+name[0]+"]").Order("price desc").Offset((id - 1) * num).Limit(num).Find(&coms)
		case 3:
			DB.Where("name regexp ?", "["+name[0]+"]").Order("sales_volume desc").Offset((id - 1) * num).Limit(num).Find(&coms)
		default:
			DB.Where("name regexp ?", "["+name[0]+"]").Offset((id - 1) * num).Limit(num).Find(&coms)
		}
	}
	return coms
}

func SearchCommByCGID(id, num, orderID int, cgid string) []Commodity {
	var cgs []CateGory
	DB.Where("name like ?", "%"+cgid+"%").Find(&cgs)
	ids := []int64{}
	for _, v := range cgs {
		ids = append(ids, v.ID)
	}
	fmt.Printf("%#v", ids)
	var coms []Commodity
	switch orderID {
	case 1:
		DB.Where("cate_gory_id in (?)", ids).Order("price asc").Offset((id - 1) * num).Limit(num).Find(&coms)
	case 2:
		DB.Where("cate_gory_id in (?)", ids).Order("price desc").Offset((id - 1) * num).Limit(num).Find(&coms)
	case 3:
		DB.Where("cate_gory_id in (?)", ids).Order("sales_volume desc").Offset((id - 1) * num).Limit(num).Find(&coms)
	default:
		DB.Where("cate_gory_id in (?)", ids).Offset((id - 1) * num).Limit(num).Find(&coms)
	}
	return coms
}

//查询商品总数
func SumCommodity(a int) (sum int) {
	DB.Table("commodities").Count(&sum)
	b := math.Ceil(float64(sum) / float64(a))
	return int(b)
}

func SumCommodity2(a int, name string) (sum int) {
	DB.Table("commodities").Where("name regexp ?", "["+name+"]").Count(&sum)
	b := math.Ceil(float64(sum) / float64(a))
	return int(b)
}

func SumCommodity3(a int, id string) (sum int) {
	var cgs []CateGory
	DB.Where("name like ?", "%"+id+"%").Find(&cgs)
	ids := []int64{}
	for _, v := range cgs {
		ids = append(ids, v.ID)
	}
	DB.Table("commodities").Where("cate_gory_id in (?)", ids).Count(&sum)
	b := math.Ceil(float64(sum) / float64(a))
	return int(b)
}

//插入shangpin
func InsertCoom(comm *Commodity) error {
	return DB.Create(comm).Error
}

//按ID查询某一商品
func SearchCoomByID(id string) Commodity {
	var comm Commodity
	DB.First(&comm, "id=?", id)
	return comm
}

//修改商品类型
func UpdateComCGS(cgid, commid string) {
	var comm Commodity
	DB.First("&comm", "id=?", commid)
	DB.Model(&comm).Update("cate_gory_id=?", cgid)
}
