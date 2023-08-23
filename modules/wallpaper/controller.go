package wallpaper

import (
	"encoding/json"
	"io"
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	WallhavenAPI = "https://wallhaven.cc/api/v1"
)

func (e *Wallpaper) Search(ctx *gin.Context) {
	var query SearchQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	page, err := strconv.Atoi(query.Page)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	// wallpaper api 的 page_size 固定为 24，接口转发时改为 48
	ch1 := make(chan SearchAPIRes)
	ch2 := make(chan SearchAPIRes)
	go chSearchWallpaper(ch1, query, 2*page-1)
	go chSearchWallpaper(ch2, query, 2*page)
	res1 := <-ch1
	res2 := <-ch2
	if res1.Error != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	if res2.Error != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	data := append(res1.Result.Data, res2.Result.Data...)
	tools.RespSuccess(ctx, data)
}

func (e *Wallpaper) GetInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	decodedURL, err := url.QueryUnescape(query.URL)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	url := strings.Replace(decodedURL, "https://wallhaven.cc", WallhavenAPI, 1)
	resp, err := http.Get(url)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	var info WallpaperInfoResp
	if err := json.Unmarshal(body, &info); err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, info.Data)
}

func chSearchWallpaper(ch chan<- SearchAPIRes, query SearchQuery, page int) {
	query.Page = strconv.Itoa(page)
	searchRes, err := searchWallpaper(query)
	if err != nil {
		ch <- SearchAPIRes{
			Error: err,
		}
	} else {
		ch <- SearchAPIRes{
			Result: searchRes,
		}
	}
}

func searchWallpaper(query SearchQuery) (SearchResp, error) {
	v := url.Values{}
	v.Set("q", query.Keywords)
	v.Set("ai_art_filter", query.AIArtFilter)
	v.Set("categories", query.Categories)
	v.Set("purity", query.Purity)
	v.Set("sorting", query.Sorting)
	v.Set("order", query.Order)
	v.Set("topRange", query.TopRange)
	v.Set("atleast", query.AtLeast)
	v.Set("resolutions", query.Resolutions)
	v.Set("ratios", query.Ratios)
	v.Set("colors", query.Colors)
	v.Set("page", query.Page)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(WallhavenAPI + "/search?" + v.Encode())
	if err != nil {
		return SearchResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SearchResp{}, err
	}
	var result SearchResp
	if err = json.Unmarshal(body, &result); err != nil {
		return SearchResp{}, err
	}
	return result, nil
}
