package tools

import (
	"CommoditySalesSystem/controller"
	"CommoditySalesSystem/models"
	"strconv"
)

func ShowCGName(id int64) string {
	var cg models.CateGory
	models.DB.First(&cg, "id=?", id)
	return cg.Name
}

//上一页
func Before(id int) int {
	if id == 1 {
		return 1
	} else {
		id--
		return id
	}
}

//下一页
func Next(id int) int {
	if id == controller.Sum {
		return id
	} else {
		id++
		return id
	}
}

//显示用户名
func ShowUserName(uid string) string {
	user := models.SearchUserByID(uid)
	return user.Name
}

func ShowUserName2(uid uint) string {
	user := models.SearchUserByID(strconv.Itoa(int(uid)))
	return user.Name
}

//显示商品名
func ShowCommName(commid string) string {
	comm := models.SearchCoomByID(commid)
	return comm.Name
}

func ShowCommName2(commid int64) string {
	comm := models.SearchCoomByID(strconv.Itoa(int(commid)))
	return comm.Img
}

func ShowCommName3(commid int64) string {
	comm := models.SearchCoomByID(strconv.Itoa(int(commid)))
	return comm.Name
}

//显示地址
func ShowAddrName(addrid uint) string {
	addr := models.SearchAddr(addrid)
	return addr.Country + addr.Province + addr.City + addr.County + addr.DetailedAddress
}

//判断是否缺货
func ShowIsShortage(i int64) string {
	switch {
	case i == 0:
		return "暂时无货"
	case i < 1000:
		return "缺货"
	default:
		return "库存充足"
	}
}
