package routers

import (
	"fmt"
	"github.com/StephenChristianW/go-movies-open/services"
	"github.com/StephenChristianW/go-movies-open/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

var basePath string

func init() {
	var err error
	basePath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	basePath = basePath + "/front/htmls/backend/"
}

func setIcon() {
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File(filepath.Join("front", "static", "images", "favicon.ico"))
	})
}

func htmlController() {
	router.GET("/admin/login", func(c *gin.Context) {
		c.File(basePath + "login.html")
	})
	router.GET("/admin/main", func(c *gin.Context) {
		c.File(basePath + "index.html")
	})
	router.GET("/admin//html/backgroundImages", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=1800")
		images, err := getHtmlBgImage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, services.Response{
				Code: -1,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, services.Response{
			Code: 1,
			Msg:  "ok",
			Data: images,
		})
	})
	// 静态目录
	router.Static("/static", filepath.Join("front", "static"))
}

func getHtmlBgImage() ([]string, error) {
	// 本地路径（兼容 Windows/Linux）
	localImageDir := filepath.Join("front", "static", "images", "admin", "adminBackGround")

	// 检查目录是否存在
	if _, err := os.Stat(localImageDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("目录不存在: %s", localImageDir)
	}

	files, err := os.ReadDir(localImageDir)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, file := range files {
		if !file.IsDir() {
			// 返回 Web 访问路径（非本地路径！）
			webPath := "/static/images/admin/adminBackGround/" + file.Name()
			images = append(images, webPath)
		}
	}
	if len(images) <= 5 {
		return images, nil
	}
	return utils.SliceRandN(images, 5), nil

}
