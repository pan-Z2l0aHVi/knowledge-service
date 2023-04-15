package wallpaper

type SearchParams struct {
	Categories  string `form:"categories" binding:"required"`
	Purity      string `form:"purity" binding:"required"`
	Sorting     string `form:"sorting" binding:"required"`
	Order       string `form:"order" binding:"required"`
	TopRange    string `form:"topRange" binding:"required"`
	AtLeast     string `form:"atleast" binding:"required"`
	Resolutions string `form:"resolutions" binding:"required"`
	Ratios      string `form:"ratios" binding:"required"`
	Colors      string `form:"colors" binding:"required"`
	Page        string `form:"page" binding:"required"`
}

type GetInfoParams struct {
	URL string `form:"url" binding:"required"`
}

type WallhavenSearchResp struct {
	Data []Wallpaper `json:"data"`
}

type WallhavenInfoResp struct {
	Data Wallpaper `json:"data"`
}
