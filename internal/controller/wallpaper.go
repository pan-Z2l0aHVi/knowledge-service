package controller

import (
	"encoding/json"
	"io"
	"knowledge-service/internal/api"
	"knowledge-service/internal/service"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type WallpaperController struct{}

func (e *WallpaperController) Search(ctx *gin.Context) {
	var query api.SearchWallpaperQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	page, err := strconv.Atoi(query.Page)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	// wallpaper api 的 page_size 固定为 24，接口转发时改为 48
	ch1 := make(chan api.SearchWallpaperAPIRes)
	ch2 := make(chan api.SearchWallpaperAPIRes)
	wallpaperS := service.WallpaperService{}
	go wallpaperS.ChSearchWallpaper(ch1, query, 2*page-1)
	go wallpaperS.ChSearchWallpaper(ch2, query, 2*page)
	res1 := <-ch1
	res2 := <-ch2
	if res1.Error != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	if res2.Error != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	data := append(res1.Result.Data, res2.Result.Data...)
	tools.RespSuccess(ctx, data)
}

func (e *WallpaperController) GetInfo(ctx *gin.Context) {
	var query api.GetWallpaperInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	decodedURL, err := url.QueryUnescape(query.URL)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	url := strings.Replace(decodedURL, "https://wallhaven.cc", consts.WallhavenAPI, 1)
	resp, err := client.Get(url)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	var info api.GetWallpaperInfoResp
	if err := json.Unmarshal(body, &info); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, info.Data)
}
