package models

import "errors"

//商品类别
type CateGory struct {
	ID        int64  `gorm:"primary_key" form:"CGID"` //并不是递增
	Name      string `gorm:"not null;unique;size:15" form:"CGName"`
	Introduce string `gorm:"type:blob" form:"CGIntro"` //介绍
}

//查询所有商品类别
func SearchAllCG() []CateGory {
	var cgs []CateGory
	DB.Find(&cgs)
	return cgs
}

//删除商品类别
func DELCG(id string) error {
	var cg CateGory
	DB.First(&cg, "id=?", id)
	var sum int
	DB.Model(&Commodity{}).Where("cate_gory_id=?", id).Count(&sum)
	if sum <= 0 {
		DB.Delete(&cg)
		return nil
	} else {
		err := errors.New("该商品类还有商品无法删除！")
		return err
	}
}

//查询某一商品类别
func SearchCGByID(id string) CateGory {
	var cg CateGory
	DB.First(&cg, "id=?", id)
	return cg
}

//修改商品类别
//func UpCg(cgid string,cg CateGory){
//	var upcg CateGory
//	DB.First(&upcg,"id=?",cgid)
//	DB.Model(&upcg).Updates(&cg)
//	DB.Model(&Commodity{}).Where("cate_gory_id=?",cgid).Update("cate_gory_id",cg.ID)
//}
