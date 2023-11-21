package service

import (
	"encoding/json"
	"io"
	"knowledge-service/internal/api"
	"knowledge-service/pkg/consts"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WallpaperService struct{}

func (e *WallpaperService) ChSearchWallpaper(ch chan<- api.SearchWallpaperAPIRes, query api.SearchWallpaperQuery, page int) {
	query.Page = strconv.Itoa(page)
	searchRes, err := e.SearchWallpaper(query)
	if err != nil {
		ch <- api.SearchWallpaperAPIRes{
			Error: err,
		}
	} else {
		ch <- api.SearchWallpaperAPIRes{
			Result: searchRes,
		}
	}
}

func (e *WallpaperService) SearchWallpaper(query api.SearchWallpaperQuery) (api.SearchWallpaperResp, error) {
	v := url.Values{}
	v.Set("apikey", consts.APIKey)
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
	resp, err := client.Get(consts.WallhavenAPI + "/search?" + v.Encode())
	if err != nil {
		return api.SearchWallpaperResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.SearchWallpaperResp{}, err
	}
	var result api.SearchWallpaperResp
	if err = json.Unmarshal(body, &result); err != nil {
		return api.SearchWallpaperResp{}, err
	}
	return result, nil
}
