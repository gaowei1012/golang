package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

var db * gorm.DB

// 初始化数据库
func init()  {
	var err error
	db, err = gorm.Open("mysql", "root:gao129231wei@/blogs?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect databases")
	}

	db.AutoMigrate(&shopModel{})
}

type (
	shopModel struct {
		gorm.Model
		Title string `json:"title"`
		Completed int `json:"completed"`
	}

	transformedShop struct {
		ID 	uint `json:"id"`
		Title	string	`json:"title"`
		Completed	bool	`json:"completed"`
	}
)

// create
func create(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	shop := shopModel{Title: c.PostForm("title"), Completed: completed}

	db.Save(&shop)

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "shop apis created successfully",  "resourceId": shop.ID})
}

// home
func fetchHome(c *gin.Context)  {
	var shops []shopModel
	var _shops []transformedShop

	db.Find(&shops)
		if len(shops) <= 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not shops Found"})
			return
		}

	for _, item := range shops {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}

		_shops = append(_shops, transformedShop{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": _shops,
	})
}

// 获取单个 shop
func fetchSingle(c *gin.Context)  {
	var shop shopModel

	shipID := c.Param("id")
	db.First(&shop, shipID)

	if shop.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No shop found!",
		})
		return
	}

	completed := false
	if shop.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_shop := transformedShop{ID: shop.ID, Title: shop.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": _shop,
	})
}

// 更新
func shopUpdate(c *gin.Context)  {
	var shop shopModel
	shopID := c.Param("id")
	db.First(&shop, shopID)

	if shop.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No shop found!",
		})
		return
	}

	db.Model(&shop).Update("title", c.PostForm("title"))
		completed, _ := strconv.Atoi(c.PostForm("completed"))
		db.Model(&shop).Update("completed", completed)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "更新成功!",
		})
}

// 删除
func shopDelete(c *gin.Context)  {
	var shop shopModel
	shopID := c.Param("id")

	db.First(&shop, shopID)
	if shop.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "删除失败! 没有匹配数据!",
		})
		return
	}
	db.Delete(&shop)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "删除成功!",
		})
}

func main() {

	router := gin.Default()
	v1 := router.Group("api/v1")

	{
		v1.POST("/", create)
		v1.GET("/", fetchHome)
		v1.GET("/:id", fetchSingle)
		v1.PUT("/:id", shopUpdate)
		v1.DELETE("/:id", shopDelete)
	}

	router.Run(":16680")

}
