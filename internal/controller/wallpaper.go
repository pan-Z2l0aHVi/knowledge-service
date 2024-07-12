package controller

import (
	"encoding/json"
	"io"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
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

// @Summary 搜索壁纸
// @Description 使用但不校验登录态
// @Produce json
// @Param query query entity.SearchWallpaperQuery true "query参数"
// @Success 200 {array} entity.WallpaperItem "ok" "壁纸列表"
// @Router /wallpaper/search [get]
func (e *WallpaperController) Search(ctx *gin.Context) {
	var query entity.SearchWallpaperQuery
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
	ch1 := make(chan entity.SearchWallpaperAPIRes)
	ch2 := make(chan entity.SearchWallpaperAPIRes)
	wallpaperS := service.WallpaperService{}
	go wallpaperS.ChSearchWallpaper(ch1, query, 2*page-1)
	go wallpaperS.ChSearchWallpaper(ch2, query, 2*page)
	res1 := <-ch1
	res2 := <-ch2
	if res1.Error != nil {
		tools.RespFail(ctx, consts.Fail, res1.Error.Error(), nil)
		return
	}
	if res2.Error != nil {
		tools.RespFail(ctx, consts.Fail, res2.Error.Error(), nil)
		return
	}
	data := append(res1.Result.Data, res2.Result.Data...)

	if uid, exist := ctx.Get("uid"); exist {
		userD := dao.UserDAO{}
		user, err := userD.FindByUserID(ctx, uid.(string))
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		collectedWallpapers := user.CollectedWallpapers
		res := []entity.WallpaperItem{}
		for _, item := range data {
			collected := false
			for _, wallpaper := range collectedWallpapers {
				if item.WallpaperID == wallpaper.WallpaperID {
					collected = true
					break
				}
			}
			res = append(res, entity.WallpaperItem{
				Wallpaper: item.Wallpaper,
				Collected: collected,
			})
		}
		tools.RespSuccess(ctx, res)
	} else {
		tools.RespSuccess(ctx, data)
	}
}

// @Summary 获取壁纸详情
// @Description 使用但不校验登录态
// @Produce json
// @Param url query string true "壁纸链接"
// @Success 200 {object} entity.GetWallpaperInfoResp "ok" "壁纸详情"
// @Router /wallpaper/info [get]
func (e *WallpaperController) GetInfo(ctx *gin.Context) {
	var query entity.GetWallpaperInfoQuery
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
	url := strings.Replace(decodedURL, "https://wallhaven.cc", consts.WallhavenAPI, 1) + "?apikey=" + consts.APIKey
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
	var info entity.GetWallpaperInfoResp
	if err := json.Unmarshal(body, &info); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, info.Data)
}
