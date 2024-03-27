package api

import (
	"encoding/base64"
	"net/http"
	"os"
	"speech/service"

	"github.com/gin-gonic/gin"
)

var voiceMap = map[string]byte{
	"wuu-CN-XiaotongNeural":         1,
	"wuu-CN-YunzheNeural":           1,
	"zh-CN-XiaoxiaoNeural":          1,
	"zh-CN-YunxiNeural":             1,
	"zh-CN-YunjianNeural":           1,
	"zh-CN-XiaoyiNeural":            1,
	"zh-CN-YunyangNeural":           1,
	"zh-CN-XiaochenNeural":          1,
	"zh-CN-XiaohanNeural":           1,
	"zh-CN-XiaomengNeural":          1,
	"zh-CN-XiaomoNeural":            1,
	"zh-CN-XiaoqiuNeural":           1,
	"zh-CN-XiaoruiNeural":           1,
	"zh-CN-XiaoshuangNeural":        1,
	"zh-CN-XiaoxuanNeural":          1,
	"zh-CN-XiaoyanNeural":           1,
	"zh-CN-XiaozhenNeural":          1,
	"zh-CN-YunfengNeural":           1,
	"zh-CN-YunhaoNeural":            1,
	"zh-CN-YunxiaNeural":            1,
	"zh-CN-YunzeNeural":             1,
	"zh-CN-henan-YundengNeural":     1,
	"zh-CN-liaoning-XiaobeiNeural":  1,
	"zh-CN-shaanxi-XiaoniNeural":    1,
	"zh-CN-shandong-YunxiangNeural": 1,
	"zh-CN-sichuan-YunxiNeural":     1,
	"zh-HK-HiuMaanNeural":           1,
	"zh-HK-WanLungNeural":           1,
	"zh-HK-HiuGaaiNeural":           1,
	"zh-TW-HsiaoChenNeural":         1,
	"zh-TW-YunJheNeural":            1,
	"zh-TW-HsiaoYuNeural":           1,
}

func Speech(c *gin.Context) {
	var body struct {
		Content string `json:"content"`
		Voice   string `json:"voice"`
	}

	if c.ShouldBindJSON(&body) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if _, ok := voiceMap[body.Voice]; body.Voice != "" && !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if body.Content == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// voices := []string{
	// 	"wuu-CN-XiaotongNeural",
	// 	"wuu-CN-YunzheNeural",
	// 	"zh-CN-XiaoxiaoNeural",
	// 	"zh-CN-YunxiNeural",
	// 	"zh-CN-YunjianNeural",
	// 	"zh-CN-XiaoyiNeural",
	// 	"zh-CN-YunyangNeural",
	// 	"zh-CN-XiaochenNeural",
	// 	"zh-CN-XiaohanNeural",
	// 	"zh-CN-XiaomengNeural",
	// 	"zh-CN-XiaomoNeural",
	// 	"zh-CN-XiaoqiuNeural",
	// 	"zh-CN-XiaoruiNeural",
	// 	"zh-CN-XiaoshuangNeural",
	// 	"zh-CN-XiaoxuanNeural",
	// 	"zh-CN-XiaoyanNeural",
	// 	"zh-CN-XiaozhenNeural",
	// 	"zh-CN-YunfengNeural",
	// 	"zh-CN-YunhaoNeural",
	// 	"zh-CN-YunxiaNeural",
	// 	"zh-CN-YunzeNeural",
	// 	"zh-CN-henan-YundengNeural",
	// 	"zh-CN-liaoning-XiaobeiNeural",
	// 	"zh-CN-shaanxi-XiaoniNeural",
	// 	"zh-CN-shandong-YunxiangNeural",
	// 	"zh-CN-sichuan-YunxiNeural",
	// 	"zh-HK-HiuMaanNeural",
	// 	"zh-HK-WanLungNeural",
	// 	"zh-HK-HiuGaaiNeural",
	// 	"zh-TW-HsiaoChenNeural",
	// 	"zh-TW-YunJheNeural",
	// 	"zh-TW-HsiaoYuNeural",
	// }

	rspChan := make(chan service.SpeechResponse, 1)
	service.WaitChan <- service.SpeechRequest{
		Content:      body.Content,
		VoiceName:    body.Voice,
		ResponseChan: rspChan,
	}

	rsp, ok := <-rspChan
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	allbytes, err := os.ReadFile(rsp.FileName)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	voiceBase64 := base64.StdEncoding.EncodeToString(allbytes)

	os.Remove(rsp.FileName)

	c.JSON(http.StatusOK, gin.H{
		"base64":   voiceBase64,
		"duration": rsp.Duration,
	})
}
