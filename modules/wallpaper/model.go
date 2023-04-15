package wallpaper

type Wallpaper struct {
	WallpaperID string   `json:"id"`
	URL         string   `json:"url"`
	ShortURL    string   `json:"short_url"`
	Views       int      `json:"views"`
	Purity      string   `json:"purity"`
	Category    string   `json:"category"`
	DimensionX  int      `json:"dimension_x"`
	DimensionY  int      `json:"dimension_y"`
	Resolution  string   `json:"resolution"`
	File_size   int      `json:"file_size"`
	File_type   string   `json:"file_type"`
	Created_at  string   `json:"created_at"`
	Colors      []string `json:"colors"`
	Path        string   `json:"path"`
	Thumbs      Thumbs   `json:"thumbs"`
	Tags        []Tag    `json:"tags"`
}

type Tag struct {
	TagID       int    `json:"id"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Category_id int    `json:"category_id"`
	Category    string `json:"category"`
	Purity      string `json:"purity"`
	Created_at  string `json:"created_at"`
}

type Thumbs struct {
	Large    string `json:"large"`
	Original string `json:"original"`
	Small    string `json:"small"`
}
