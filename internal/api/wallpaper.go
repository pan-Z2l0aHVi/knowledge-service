package api

import "knowledge-service/internal/model"

type SearchWallpaperQuery struct {
	Keywords    string `form:"q"`
	Categories  string `form:"categories" binding:"required"`
	Purity      string `form:"purity" binding:"required"`
	Sorting     string `form:"sorting" binding:"required"`
	Order       string `form:"order" binding:"required"`
	TopRange    string `form:"topRange"`
	AtLeast     string `form:"atleast"`
	Resolutions string `form:"resolutions"`
	Ratios      string `form:"ratios"`
	Colors      string `form:"colors"`
	AIArtFilter string `form:"ai_art_filter" binding:"required"`
	Page        string `form:"page" binding:"required"`
}

type GetWallpaperInfoQuery struct {
	URL string `form:"url" binding:"required"`
}

type SearchWallpaperResp struct {
	Data []model.Wallpaper `json:"data"`
}

type GetWallpaperInfoResp struct {
	Data model.Wallpaper `json:"data"`
}

type SearchWallpaperAPIRes struct {
	Error  error
	Result SearchWallpaperResp
}
