package backend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/evercyan/letitgo/crypto"
	"github.com/evercyan/letitgo/request"
	"github.com/evercyan/letitgo/util"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails"
)

var (
	emojiUrl   = "https://raw.githubusercontent.com/evercyan/cantor/master/assets/emoji"
	appLogFile = "/tmp/emoji.log"
	ffmpegPath = "/usr/local/bin/ffmpeg"
)

/**
 * ********************************
 */

type App struct {
	RT  *wails.Runtime
	Log *logrus.Logger
	Url map[string]string
}

func (a *App) WailsInit(runtime *wails.Runtime) error {
	a.RT = runtime
	a.initLogger()
	a.Url = map[string]string{
		"config": emojiUrl + "/template.json",
		"ass":    emojiUrl + "/%s/template.ass",
		"mp4":    emojiUrl + "/%s/template.mp4",
	}
	a.Log.Info("WailsInit url ", crypto.JsonEncode(a.Url))
	return nil
}

func (a *App) WailsShutdown() {
	a.Log.Info("WailsShutdown")
}

func (a *App) initLogger() {
	a.Log = logrus.New()
	logFile, err := os.OpenFile(appLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic("无法创建日志文件: " + appLogFile)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	a.Log.SetOutput(mw)
	a.Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func (a *App) resp(code int, data interface{}) string {
	return crypto.JsonEncode(map[string]interface{}{
		"code": code,
		"data": data,
	})
}

/**
 * ********************************
 */

func (a *App) GetTplList(param string) string {
	a.Log.Info("GetTplList param ", param)
	content, err := request.Get(a.Url["config"])
	if err != nil {
		return a.resp(-1, "读取配置文件失败")
	}
	replaceMap := map[string]string{
		"\n": "",
		" ":  "",
	}
	for k, v := range replaceMap {
		content = strings.Replace(content, k, v, -1)
	}
	return a.resp(0, content)
}

func (a *App) BuildGif(param string) string {
	a.Log.Info("BuildGif param ", param)
	req := struct {
		Code     string   `json:"code"`
		TextList []string `json:"text_list"`
	}{}
	if err := json.Unmarshal([]byte(param), &req); err != nil {
		return a.resp(-1, err.Error())
	}

	// 模板标识 + 台词生成 md5
	md5 := util.Md5(req.Code + crypto.JsonEncode(req.TextList))
	outputGif := fmt.Sprintf("/tmp/%s_%s.gif", req.Code, md5)
	a.Log.Info("BuildGif outputGif ", outputGif)
	if !util.IsExist(outputGif) {
		assContent, err := request.Get(fmt.Sprintf(a.Url["ass"], req.Code))
		if err != nil {
			return a.resp(-1, "读取模板台词文件失败")
		}

		// 生成台词文件
		outputAss := fmt.Sprintf("/tmp/%s_%s.ass", req.Code, md5)
		file, err := os.Create(outputAss)
		if err != nil {
			return a.resp(-1, err.Error())
		}
		tpl := template.Must(template.New("assfile").Parse(assContent))
		err = tpl.Execute(file, map[string][]string{
			"sentences": req.TextList,
		})
		if err != nil {
			return a.resp(-1, err.Error())
		}

		// 调用 ffmpeg 生成 gif
		mp4FileUrl := fmt.Sprintf(a.Url["mp4"], req.Code)
		cmdStr := fmt.Sprintf("%s -i %s -r 8 -vf ass=%s,scale=300:-1 -y %s", ffmpegPath, mp4FileUrl, outputAss, outputGif)
		a.Log.Info("BuildGif cmdStr ", cmdStr)
		args := []string{
			"-i", mp4FileUrl, "-r", "8", "-vf", fmt.Sprintf("ass=%s,scale=300:-1", outputAss), "-y", outputGif,
		}
		_, err = exec.Command(ffmpegPath, args...).Output()
		os.Remove(outputAss)
		if err != nil {
			return a.resp(-1, err.Error())
		}
	}

	// 不允许读本地图片, base64 再出去
	return a.resp(0, map[string]string{
		"path": outputGif,
		"data": "data:image/gif;base64," + crypto.Base64Encode(util.ReadFile(outputGif)),
	})
}

func (a *App) DownloadGif(gifPath string) string {
	selectPath := a.RT.Dialog.SelectDirectory()
	a.Log.Info("DownloadGif selectPath ", selectPath)
	if selectPath == "" {
		return a.resp(0, "")
	}
	if !util.IsExist(gifPath) {
		return a.resp(-1, "文件不存在")
	}
	content := util.ReadFile(gifPath)
	if content == "" {
		return a.resp(-1, "读取文件失败")
	}

	filePath := selectPath + "/" + path.Base(gifPath)
	a.Log.Info("DownloadGif filePath ", filePath)
	if err := util.WriteFile(filePath, content); err != nil {
		return a.resp(-1, "下载文件失败: "+err.Error())
	}
	return a.resp(0, "下载成功")
}
