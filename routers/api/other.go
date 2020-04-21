package api

import (
	"github.com/gin-gonic/gin"
	"learning/conf"
	"learning/utils"
	"log"
	"net/http"
	"os"
	"time"
)

func PlayVideo(c *gin.Context) {
	videoName := c.Param("name")
	localPath := conf.AppConfig.Path.Video + "/" + videoName
	video, err := os.Open(localPath)
	defer func() {
		if err := video.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.Header("Content-Type", "video/mp4")
	http.ServeContent(c.Writer, c.Request, "", time.Now(), video)
}

type ExecuteProgramForm struct {
	Language string `form:"language" binding:"required"`
	Input    string `form:"input" binding:"required"`
}

func ExecuteProgram(c *gin.Context) {
	var form ExecuteProgramForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	if out, err := utils.ExecuteProgramSubject(form.Language, form.Input); err == nil {
		c.JSON(http.StatusOK, out)
	} else {
		c.JSON(http.StatusInternalServerError, "")
	}
}
