package controller

import (
	"CommoditySalesSystem/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"path"
	"strconv"
	"time"
)

var Sum int

//var Exits bool

func ShowMain(c *gin.Context) {
	var islogin bool
	var isAdmain bool
	var messsgeOrder int
	var makeSureArrive int
	var sessions models.Session
	u := new(models.GuestOrAdministrators)
	session, err := c.Cookie("user")
	fmt.Println(err)
	fmt.Println(session)
	if err != nil {
		islogin = false
	} else {
		islogin = true
		sessions = models.SearchSession(session)
		if sessions.UserId <= 0 {
			islogin = false
		} else {
			var user models.GuestOrAdministrators
			user.ID = sessions.UserId
			u = models.SeachOneGA(&user)
			if u.IsAdministrators == 1 {
				isAdmain = false
				models.DB.Table("orders").Where("is_handel=1 and guest_id= ?", u.ID).Count(&messsgeOrder)
				models.DB.Table("orders").Unscoped().Where("is_handel=3 and guest_id= ?", u.ID).Count(&makeSureArrive)
			} else {
				isAdmain = true
				models.DB.Table("orders").Where("is_handel=2").Count(&messsgeOrder)
			}
		}
	}
	id, _ := strconv.Atoi(c.Query("id"))
	if id == 0 {
		id = 1
	}
	order := c.Query("order")
	orderID, _ := strconv.Atoi(order)
	coms := []models.Commodity{}
	var Byname bool
	var BycgID bool
	commname, ok := c.GetQuery("commName")
	cg, ok2 := c.GetQuery("cgid")
	fmt.Println("-----------")
	fmt.Println(cg)
	if ok2 {
		BycgID = true
		coms = models.SearchCommByCGID(id, 9, orderID, cg)
		Sum = models.SumCommodity3(9, cg)
	}
	if !ok && !ok2 {
		Byname = false
		coms = models.SearchCommodity(id, 9, orderID)
		Sum = models.SumCommodity(9)
	} else if !ok2 {
		Byname = true
		coms = models.SearchCommodity(id, 9, orderID, commname)
		Sum = models.SumCommodity2(9, commname)
	}
	for i, _ := range coms {
		a := []rune(coms[i].Details)
		if len(a) < 15 {
			coms[i].BriefIntroduction = string(a)
		} else {
			coms[i].BriefIntroduction = string(a[0:15]) + "..."
		}
	}

	c.HTML(200, "main.tmpl", gin.H{
		"login":       islogin,
		"coms":        coms,
		"user":        u,
		"isAdmain":    isAdmain,
		"sum":         Sum,
		"id":          id,
		"orderID":     orderID,
		"byName":      Byname,
		"commname":    commname,
		"BycgID":      BycgID,
		"cgid":        cg,
		"messagesum":  messsgeOrder,
		"makesureArr": makeSureArrive,
	})
}

//显示登录处理器
func ShowLogin(c *gin.Context) {
	rand.Seed(time.Now().Unix())
	var c1 rune
	str := ""
	for i := 0; i < 4; i++ {
		a := rand.Intn(26) + 97
		c1 = rune(a)
		str += string(c1)
	}
	if value, ok := c.GetQuery("SCode"); ok {
		c.HTML(200, "login.tmpl", gin.H{
			"str":   str,
			"value": value,
		})
	} else {
		c.HTML(200, "login.tmpl", gin.H{
			"str":   str,
			"value": "200",
		})
	}

}

//处理登录
func HandleLogin(c *gin.Context) {
	var user models.GuestOrAdministrators
	c.ShouldBind(&user)
	u := models.SeachOneGA(&user)
	if u.ID <= 0 {
		c.Redirect(302, "/login?SCode=500")
		return
	}
	value := time.Now().Format("2006/01/02 15:04:05.000000000")
	c.SetCookie("user", value, 60*60*24, "/", "localhost", false, true)
	err := models.InsertSession(value, u.Name, u.ID)
	if err != nil {
		c.JSON(200, "session出错")
		return
	}
	c.Redirect(302, "/main?id=1")
}

//显示注册
func ShowRegist(c *gin.Context) {
	var IsAdmain bool
	var Update bool
	var user models.GuestOrAdministrators
	_, ok1 := c.GetQuery("isAdmain")
	if ok1 {
		IsAdmain = true
	}
	uid, ok := c.GetQuery("uid")
	if ok {
		Update = true
		user = models.SearchUserByID(uid)
	}
	c.HTML(200, "registe.tmpl", gin.H{
		"isUpdate": Update,
		"user":     user,
		"isAdmain": IsAdmain,
	})
}

//处理注册
func HandleRegiste(c *gin.Context) {
	var u models.GuestOrAdministrators
	c.ShouldBind(&u)
	img, err := c.FormFile("img")
	if err != nil {
		c.Redirect(302, "/registe")
		return
	}
	img.Filename = strconv.Itoa(int(time.Now().Unix())) + img.Filename
	addr := path.Join("./static/img/" + img.Filename)
	c.SaveUploadedFile(img, addr)
	u.Img = "/" + addr
	isAdmain := c.PostForm("isAdmain")
	if isAdmain == "true" {
		u.IsAdministrators = 2
	} else {
		u.IsAdministrators = 1
	}
	err = models.InsertGA(&u)
	fmt.Println("------------")
	fmt.Println(err)
	if err != nil {
		c.Redirect(302, "/registe")
		return
	}
	c.Redirect(302, "/login")
}

//修改个人信息
func UpdateUser(c *gin.Context) {
	id := c.PostForm("id")
	upuser := models.SearchUserByID(id)
	var user models.GuestOrAdministrators
	c.ShouldBind(&user)
	file, err := c.FormFile("img")
	if err == nil {
		file.Filename = strconv.Itoa(int(time.Now().Unix())) + file.Filename
		addr := path.Join("./static/img/", file.Filename)
		c.SaveUploadedFile(file, addr)
		user.Img = "/" + addr
	}
	models.DB.Model(&upuser).Updates(&user)
	c.Redirect(302, "/main?id=1")
}

//退出
func Exit(c *gin.Context) {
	session, err := c.Cookie("user")
	if err != nil {
		c.Redirect(302, "/main?id=1")
		return
	}
	sessions := models.SearchSession(session)
	if sessions.UserId <= 0 {
		c.Redirect(302, "/main?id=1")
		return
	}
	models.DeleteSession(&sessions)
	c.SetCookie("user", "", -1, "/", "localhost", false, true)
	c.Redirect(302, "/main?id=1")
	//Exits=false
}

//显示商品管理界面
func ShowComInfo(c *gin.Context) {
	a := c.Query("id")
	id, _ := strconv.Atoi(a)
	coms := models.SearchCommodity(id, 12, 0)
	Sum = models.SumCommodity(12)
	c.HTML(200, "CommodityInformation.tmpl", gin.H{
		"coms": coms,
		"sum":  Sum,
		"id":   id,
	})
}

//添加商品
func HandleComm(c *gin.Context) {
	//c.SetCookie("cg","",-1,"/","localhost",false,true)
	var comm models.Commodity
	c.ShouldBind(&comm)
	//fmt.Println(comm.Details)
	stock := c.PostForm("Stock")
	comStock, _ := strconv.Atoi(stock)
	comm.Stock = &sql.NullInt64{Int64: int64(comStock), Valid: true}
	comm.SalesVolume = &sql.NullInt64{Int64: 0, Valid: true}
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(200, "错误")
		return
	}
	file.Filename = strconv.Itoa(int(time.Now().Unix())) + file.Filename
	addr := path.Join("./static/img/", file.Filename)
	c.SaveUploadedFile(file, addr)
	comm.Img = "/" + addr
	//fmt.Printf("%+v",comm)
	err = models.InsertCoom(&comm)
	if err != nil {
		c.JSON(200, "错误")
		return
	}
	c.Redirect(302, "/Admain/showComInfo?id=1")
}

//显示商品添加或者修改界面
func ShowAddComm(c *gin.Context) {
	var isUpDate bool
	var readonly bool
	var comm models.Commodity
	var cgid int
	commid, ok := c.GetQuery("commid")
	_, ok2 := c.GetQuery("readonly")
	if ok2 {
		readonly = true
	}
	if ok {
		isUpDate = true
		comm = models.SearchCoomByID(commid)
	} else {
		isUpDate = false
	}
	id, ok := c.GetQuery("cgid")
	//c.SetCookie("cg",id,60*60*24,"/","localhost",false,true)
	if ok {
		cgid, _ = strconv.Atoi(id)
	}
	c.HTML(200, "addComm.tmpl", gin.H{
		"isUpDate": isUpDate,
		"comm":     comm,
		"cgid":     int64(cgid),
		"readonly": readonly,
	})
}

//修改商品
func UpdateComm(c *gin.Context) {
	var com models.Commodity
	id := c.PostForm("comID")
	models.DB.First(&com, "id=?", id)
	var comm models.Commodity
	c.ShouldBind(&comm)
	cgName := c.PostForm("leibie")
	var cg models.CateGory
	models.DB.First(&cg, "name=?", cgName)
	if cg.ID <= 0 {
		c.Redirect(302, "/Admain/addComm?commid="+id)
		return
	}
	comm.CateGoryID = cg.ID
	stock := c.PostForm("Stock")
	comStock, _ := strconv.Atoi(stock)
	comm.Stock = &sql.NullInt64{Int64: int64(comStock), Valid: true}
	savm := c.PostForm("SaVm")
	comSaVm, _ := strconv.Atoi(savm)
	comm.SalesVolume = &sql.NullInt64{Int64: int64(comSaVm), Valid: true}
	file, err := c.FormFile("img")
	if err == nil {
		file.Filename = strconv.Itoa(int(time.Now().Unix())) + file.Filename
		addr := path.Join("./static/img/", file.Filename)
		c.SaveUploadedFile(file, addr)
		comm.Img = "/" + addr
	}
	models.DB.Model(&com).Updates(&comm)
	c.Redirect(302, "/Admain/showComInfo?id=1")
}

//删除商品
func DeleComm(c *gin.Context) {
	id := c.Query("id")
	models.DB.Where("id=?", id).Delete(models.Commodity{})
	c.Redirect(302, "/Admain/showComInfo?id=1")
}

//显示商品类别和添加商品选中商品类别
func ShowCateGory(c *gin.Context) {
	cgs := models.SearchAllCG()
	var isChoose bool
	var err bool
	_, ok1 := c.GetQuery("ischoose")
	if !ok1 {
		isChoose = false
	} else {
		isChoose = true
	}
	if _, ok2 := c.GetQuery("err"); ok2 {
		err = true
		fmt.Println("----------")
	}

	c.HTML(200, "CateGory.tmpl", gin.H{
		"cgs":      cgs,
		"ischoose": isChoose,
		"err":      err,
	})
}

//显示添加商品类别界面
func ShowAddCG(c *gin.Context) {
	var isUpdate bool
	cgid, ok := c.GetQuery("cgid")
	if ok {
		isUpdate = true
	}
	cg := models.SearchCGByID(cgid)
	c.HTML(200, "ShowAddCG.tmpl", gin.H{
		"isUpdate": isUpdate,
		"cg":       cg,
	})
}

//处理商品类别添加
func HandleAddCG(c *gin.Context) {
	var cg models.CateGory
	c.ShouldBind(&cg)
	err := models.DB.Create(&cg).Error
	if err != nil {
		c.Redirect(302, "/Admain/showAddCG")
		return
	}
	c.Redirect(302, "/Admain/showCG")
}

//删除商品类别
func DeleteCG(c *gin.Context) {
	id := c.Query("cgid")
	err := models.DELCG(id)
	if err == nil {
		c.Redirect(302, "/Admain/showCG")
	} else {
		c.Redirect(302, "/Admain/showCG?err=true")
	}
}

//修改商品类别
func UpdateCG(c *gin.Context) {
	cgid := c.PostForm("CGID")
	var upcg models.CateGory
	models.DB.First(&upcg, "id=?", cgid)
	var cg models.CateGory
	c.ShouldBind(&cg)
	models.DB.Model(&upcg).Updates(&cg)
	c.Redirect(302, "/Admain/showCG")
}

//主界面按商品名称搜索
func SearchCommByName(c *gin.Context) {
	commName := c.PostForm("commName")
	c.Redirect(302, "/main?id=1&order=0&commName="+commName)
}

//主界面按商品类型名称搜索
func SearchCommByCGName(c *gin.Context) {
	cgName := c.PostForm("cgName")
	/*var cg models.CateGory
	models.DB.First(&cg, "name like ?", cgName)*/
	c.Redirect(302, "/main?id=1&order=0&cgid="+cgName)
}

//显示地址界面
func ShowAddr(c *gin.Context) {
	orderID, ok2 := c.GetQuery("orderID")
	if ok2 {
		c.SetCookie("order", orderID, 60*60*24, "/", "localhost", false, true)
	}
	uid := c.Query("uid")
	var ischoose bool
	_, ok := c.GetQuery("ischoose")
	if ok {
		ischoose = true
	}
	addrs := models.SearchAddrByGID(uid)
	c.HTML(200, "Addr.tmpl", gin.H{
		"addres":   addrs,
		"uid":      uid,
		"ischoose": ischoose,
	})
}

//显示添加或者修改地址界面
func ShowAddAddr(c *gin.Context) {
	var isUpdate bool
	var addr models.Addr
	_, ok := c.GetQuery("isUpdate")
	if ok {
		id := c.Query("addrid")
		models.DB.First(&addr, "id=?", id)
		isUpdate = true
	} else {
		isUpdate = false
	}
	uid := c.Query("uid")
	fmt.Println("update: ", isUpdate)
	c.HTML(200, "AddAddr.tmpl", gin.H{
		"uid":      uid,
		"isUpdate": isUpdate,
		"addr":     addr,
	})
}

//处理添加地址
func HandleAddAddr(c *gin.Context) {
	var addr models.Addr
	c.ShouldBind(&addr)
	models.DB.Create(&addr)
	c.Redirect(302, "showAddr?uid="+strconv.Itoa(int(addr.GuestID)))
}

//删除地址
func DelteAddr(c *gin.Context) {
	AddrID := c.Query("Addrid")
	uid := c.Query("uid")
	models.DelAddr(AddrID)
	c.Redirect(302, "showAddr?uid="+uid)
}

//修改地址
func UpdateAddr(c *gin.Context) {
	addrID := c.PostForm("AddrID")
	var updateAddr models.Addr
	models.DB.First(&updateAddr, "id=?", addrID)
	var addr models.Addr
	c.ShouldBind(&addr)
	models.DB.Model(&updateAddr).Updates(addr)
	c.Redirect(302, "showAddr?uid="+strconv.Itoa(int(updateAddr.GuestID)))
}

//显示添加购物车界面
func ShowAddOrder(c *gin.Context) {
	var max bool
	sessID, _ := c.Cookie("user")
	session := models.SearchSession(sessID)
	uid := session.UserId
	commid := c.Query("commid")
	comm := models.SearchCoomByID(commid)
	sum := comm.Stock.Int64
	if sum == 0 {
		max = true
	}
	c.HTML(200, "Addorder.tmpl", gin.H{
		"uid":    uid,
		"commid": commid,
		"sum":    sum,
		"max":    max,
	})
}

//处理添加购物车
func HandleAddOrder(c *gin.Context) {
	var order models.Order
	c.ShouldBind(&order)
	order.IsHandel = 1
	user := models.SearchUserByID(strconv.Itoa(int(order.GuestID)))
	comm := models.SearchCoomByID(strconv.Itoa(int(order.ComID)))
	sum := comm.Stock.Int64
	//fmt.Println("sum: ",sum)
	//fmt.Println("orderSum",order.ComSum)
	if sum < order.ComSum {
		c.Redirect(302, "/guest/showAddOrder?uid="+strconv.Itoa(int(user.ID))+"&commid="+strconv.Itoa(int(comm.ID)))
		return
	}
	price := float64(order.ComSum) * (1 - models.Discount(&user)) * comm.Price
	order.Price = price
	order.Tel = user.Tel
	order.CGID = comm.CateGoryID
	models.DB.Create(&order)
	c.SetCookie("order", strconv.Itoa(int(order.ID)), 60*60*24, "/", "localhost", false, true)
	c.Redirect(302, "/guest/showAddr?uid="+strconv.Itoa(int(order.GuestID))+"&ischoose=true")
}

func HandelChoseAddr(c *gin.Context) {
	addrid := c.Query("addrid")
	addrID, _ := strconv.Atoi(addrid)
	orderid, _ := c.Cookie("order")
	order := models.SearchOneOrder(orderid)
	models.DB.Model(&order).Update("addr_id", addrID)
	c.SetCookie("addrid", "", -1, "/", "localhost", false, true)
	c.Redirect(302, "/guest/searchOrder?isHandle=1")
}

//显示未结算的订单信息
func ShowOrder(c *gin.Context) {
	//addrid,ok:=c.GetQuery("addrID")
	var MakeSure bool
	var WaitHandle bool
	var orders []models.Order
	sessionid, _ := c.Cookie("user")
	session := models.SearchSession(sessionid)
	isHandle := c.Query("isHandle")
	if isHandle == "3" {
		orders, _ = models.SearchOrder(isHandle, session.UserId)
		MakeSure = true
	} else if isHandle == "2" {
		WaitHandle = true
		orders = models.SearchALLOrder(isHandle)
	} else {
		orders, _ = models.SearchOrder(isHandle, session.UserId)
	}
	//fmt.Println("makesure: ",MakeSure)
	c.HTML(200, "showOrders.tmpl", gin.H{
		"orders":     orders,
		"makeSure":   MakeSure,
		"waitHandle": WaitHandle,
	})
}

//删除订单
func DelOrder(c *gin.Context) {
	orderID := c.Query("orderID")
	models.DelteOrder(orderID)
	c.Redirect(302, "/guest/searchOrder?isHandle=1")
}

//结算订单
func SettlementOrder(c *gin.Context) {
	orderid := c.Query("orderID")
	order := models.SearchOneOrder(orderid)
	user := models.SearchUserByID(strconv.Itoa(int(order.GuestID)))
	comm := models.SearchCoomByID(strconv.Itoa(int(order.ComID)))
	user.ConsumptionAmount += order.Price
	models.DB.Save(&user)
	comm.SalesVolume.Int64 += order.ComSum
	comm.Stock.Int64 -= order.ComSum
	models.DB.Save(&comm)
	order.IsHandel = 2
	models.DB.Save(&order)
	c.Redirect(302, "/guest/searchOrder?isHandle=1")
}

//显示修改商品数量界面
func ShowUpdateOrSum(c *gin.Context) {
	orderID := c.Query("orderID")
	order := models.SearchOneOrder(orderID)
	coom := models.SearchCoomByID(strconv.Itoa(int(order.ComID)))
	c.HTML(200, "UpdateOderSum.tmpl", gin.H{
		"order": order,
		"comm":  coom,
	})
}

//处理商品数量修改
func HandleUpdateOrSum(c *gin.Context) {
	orderID := c.PostForm("orderID")
	order := models.SearchOneOrder(orderID)
	sum := c.PostForm("sum")
	orderSum, _ := strconv.Atoi(sum)
	comm := models.SearchCoomByID(strconv.Itoa(int(order.ComID)))
	if int64(orderSum) > comm.Stock.Int64 {
		c.Redirect(302, "/guest/showUpdateComSum?orderID="+strconv.Itoa(int(order.ID)))
		return
	}
	order.ComSum = int64(orderSum)
	user := models.SearchUserByID(strconv.Itoa(int(order.GuestID)))
	price := float64(order.ComSum) * (1 - models.Discount(&user)) * comm.Price
	order.Price, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", price), 64)
	models.DB.Save(order)
	c.Redirect(302, "/guest/searchOrder?isHandle=1")
}

//管理员处理订单
func HandleOrder(c *gin.Context) {
	orderID := c.Query("orderID")
	order := models.SearchOneOrder(orderID)
	order.IsHandel = 3
	fmt.Println("------------------")
	now := time.Now().Add(time.Hour * 24 * 7).Format("2006年01月02日")
	order.ArriveTime = now
	order.DeliveryOrNot = true
	models.DB.Debug().Save(&order)
	models.DB.Delete(&order)
	c.Redirect(302, "/Admain/searchOrder?isHandle=2")
}

//用户确认收货
func MakeSureIsArrive(c *gin.Context) {
	orderID := c.Query("orderID")
	var order models.Order
	models.DB.Unscoped().Where("id=?", orderID).First(&order)
	models.DB.Debug().Unscoped().Model(&order).Update("is_handel", 4)
	c.Redirect(302, "/guest/searchOrder?isHandle=3")
}

//管理员查询销售状况
func ShowSalesStatus(c *gin.Context) {
	c.HTML(200, "Salesstatus.tmpl", nil)
}

//处理查询
func HandleSalesStatus(c *gin.Context) {
	var orders []models.Order
	t := c.PostForm("time")
	if t == "" {
		t = "2006-01-01"
	}
	timeObj, _ := models.StrToTime(t)
	cgName := c.PostForm("cg")
	var cg models.CateGory
	if cgName == "" {
		cgName = "全部"
		models.DB.Debug().Unscoped().Find(&orders, "deleted_at>=? and delivery_or_not=?", timeObj, true)
		//fmt.Printf("%#v\n",orders)
	} else {
		models.DB.Debug().First(&cg, "name=?", cgName)
		if cg.ID <= 0 {
			c.Redirect(302, "/Admain/showSStus")
			return
		}
		models.DB.Debug().Unscoped().Find(&orders, "deleted_at>=? and delivery_or_not=true and cg_id=?", timeObj, cg.ID)
	}
	var price float64
	var comSum int64 = 0
	for _, value := range orders {
		price += value.Price
		comSum += value.ComSum
		fmt.Println(value.ComID)
	}
	fmt.Println(len(orders))
	fmt.Println("--------------------------------")
	orders = models.NoRepead(orders)
	fmt.Println(len(orders))
	for i, _ := range orders {
		orders[i].Persent = (float64(orders[i].ComSum) / float64(comSum)) * 100
		s := fmt.Sprintf("%.2f", orders[i].Persent)
		p, _ := strconv.ParseFloat(s, 64)
		orders[i].Persent = p
		s = fmt.Sprintf("%.2f", orders[i].Price)
		p, _ = strconv.ParseFloat(s, 64)
		orders[i].Price = p
		fmt.Println(orders[i].ComID)
	}
	//fmt.Printf("%#v\n",persent)
	c.HTML(200, "showSalesStatus.tmpl", gin.H{
		"time":   timeObj.Format("2006/01/02") + " 至 " + time.Now().Format("2006/01/02"),
		"cg":     cgName,
		"price":  price,
		"comSum": comSum,
		"orders": orders,
	})
}
