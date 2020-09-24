package main

import (
	"CommoditySalesSystem/controller"
	"CommoditySalesSystem/models"
	_ "CommoditySalesSystem/models"
	"CommoditySalesSystem/tools"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
	//"path/filepath"
	//"github.com/gin-contrib/multitemplate"
)

//中间件
func m1(c *gin.Context) {
	session, err := c.Cookie("user")
	if err != nil {
		c.Redirect(302, "/main?id=1")
		return
	}
	sess := models.SearchSession(session)
	if sess.UserId <= 0 {
		c.Redirect(302, "/main?id=1")
		return
	}
	user := models.SearchUserByID(strconv.Itoa(int(sess.UserId)))
	if user.IsAdministrators == 1 {
		c.Redirect(302, "/main?id=1")
		return
	} else {
		c.Next()
	}
}

func m2(c *gin.Context) {
	session, err := c.Cookie("user")
	if err != nil {
		c.Redirect(302, "/main?id=1")
		return
	}
	sess := models.SearchSession(session)
	if sess.UserId <= 0 {
		c.Redirect(302, "/main?id=1")
		return
	} else {
		c.Next()
	}
}

func m3(c *gin.Context) {
	session, err := c.Cookie("user")
	if err != nil {
		c.Redirect(302, "/main?id=1")
		return
	}
	sess := models.SearchSession(session)
	if sess.UserId <= 0 {
		c.Redirect(302, "/main?id=1")
		return
	}
	user := models.SearchUserByID(strconv.Itoa(int(sess.UserId)))
	if user.IsAdministrators == 2 {
		c.Redirect(302, "/main?id=1")
		return
	} else {
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.SetFuncMap(template.FuncMap{
		"showCGName":     tools.ShowCGName,
		"before":         tools.Before,
		"next":           tools.Next,
		"userName":       tools.ShowUserName,
		"commName":       tools.ShowCommName,
		"userName2":      tools.ShowUserName2,
		"commName2":      tools.ShowCommName2,
		"coomName3":      tools.ShowCommName3,
		"addrName":       tools.ShowAddrName,
		"ShowIsShortage": tools.ShowIsShortage,
	})
	r.LoadHTMLGlob("./templates/*.tmpl")
	//r.HTMLRender = loadTemplates("./templates")
	r.GET("/main", controller.ShowMain)                          //显示主界面
	r.GET("/login", controller.ShowLogin)                        //显示登录界面
	r.POST("/login", controller.HandleLogin)                     //处理登录
	r.GET("/registe", controller.ShowRegist)                     //显示注册页
	r.POST("/registe", controller.HandleRegiste)                 //处理注册页
	r.GET("/exit", controller.Exit)                              //处理退出
	r.POST("/searchComm", controller.SearchCommByName)           //按商品名称搜索
	r.POST("/SearchCommByCGName", controller.SearchCommByCGName) //按商品类别名搜索
	AdmainGroup := r.Group("/Admain")
	AdmainGroup.Use(m1)
	{
		AdmainGroup.GET("/showComInfo", controller.ShowComInfo) //显示所有商品信息
		//AdmainGroup.GET("/addComm",controller.ShowAddComm)//显示添加商品界面
		AdmainGroup.POST("/addComm", controller.HandleComm)            //处理添加商品
		AdmainGroup.GET("/DeleComm", controller.DeleComm)              //删除商品
		AdmainGroup.POST("/upDateComm", controller.UpdateComm)         //修改商品
		AdmainGroup.GET("/showCG", controller.ShowCateGory)            //显示商品类型界面
		AdmainGroup.GET("/showAddCG", controller.ShowAddCG)            //显示添加或修改商品类别界面
		AdmainGroup.POST("/handleAddCG", controller.HandleAddCG)       //处理添加商品类别
		AdmainGroup.GET("/delCG", controller.DeleteCG)                 //删除商品类别
		AdmainGroup.POST("/handleUpdateCG", controller.UpdateCG)       //修改商品类型
		AdmainGroup.GET("/searchOrder", controller.ShowOrder)          //显示待处理订单
		AdmainGroup.GET("/HandleOrder", controller.HandleOrder)        //管理员处理订单
		AdmainGroup.GET("/showSStus", controller.ShowSalesStatus)      //显示查询销售情况界面
		AdmainGroup.POST("/handleSStus", controller.HandleSalesStatus) //处理查询
		AdmainGroup.GET("/showSStuCG")
	}
	Public := r.Group("/public")
	Public.Use(m2)
	{
		Public.GET("/addComm", controller.ShowAddComm)       //显示添加商品界面
		Public.POST("/upDateRegiste", controller.UpdateUser) //修改个人信息
	}
	Guest := r.Group("guest")
	Guest.Use(m3)
	{
		Guest.GET("/showAddr", controller.ShowAddr)              //显示总地址
		Guest.GET("/addAddr", controller.ShowAddAddr)            //显示添加地址页面
		Guest.POST("/addAddr", controller.HandleAddAddr)         //处理添加地址
		Guest.GET("/delAddr", controller.DelteAddr)              //删除地址
		Guest.POST("/upDateAddr", controller.UpdateAddr)         //修改地址
		Guest.GET("/showAddOrder", controller.ShowAddOrder)      //显示购物车添加界面
		Guest.POST("/handleAddOrder", controller.HandleAddOrder) //处理添加购物车
		Guest.GET("/handlechose", controller.HandelChoseAddr)
		Guest.GET("/searchOrder", controller.ShowOrder)                //显示待处理的购物车
		Guest.GET("/settle", controller.SettlementOrder)               //处理结算
		Guest.GET("/delOrder", controller.DelOrder)                    //删除订单
		Guest.GET("/showUpdateComSum", controller.ShowUpdateOrSum)     //显示修改订单商品数量界面
		Guest.POST("/handleUpdateOrSum", controller.HandleUpdateOrSum) //处理修改订单数量
		Guest.GET("/isArrive", controller.MakeSureIsArrive)            //确定收货
	}

	r.Run(":9000")

}
