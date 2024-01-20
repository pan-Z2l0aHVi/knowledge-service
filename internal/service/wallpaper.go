package service

import (
	"encoding/json"
	"io"
	"knowledge-service/internal/entity"
	"knowledge-service/pkg/consts"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WallpaperService struct{}

func (e *WallpaperService) ChSearchWallpaper(ch chan<- entity.SearchWallpaperAPIRes, query entity.SearchWallpaperQuery, page int) {
	query.Page = strconv.Itoa(page)
	searchRes, err := e.SearchWallpaper(query)
	if err != nil {
		ch <- entity.SearchWallpaperAPIRes{
			Error: err,
		}
	} else {
		ch <- entity.SearchWallpaperAPIRes{
			Result: searchRes,
		}
	}
}

func (e *WallpaperService) SearchWallpaper(query entity.SearchWallpaperQuery) (entity.SearchWallpaperResp, error) {
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
		return entity.SearchWallpaperResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.SearchWallpaperResp{}, err
	}
	var result entity.SearchWallpaperResp
	if err = json.Unmarshal(body, &result); err != nil {
		return entity.SearchWallpaperResp{}, err
	}
	return result, nil
}
