package core

type Image struct {
	SQLModel
	Url       string `json:"url" gorm:"column:url;" db:"url"`
	Width     int64  `json:"width" gorm:"column:width;" db:"width"`
	Height    int64  `json:"height" gorm:"column:height;" db:"height"`
	Extension string `json:"extension" gorm:"column:extension;" db:"extension"`
}

func (Image) TableName() string { return "images" }

func NewImage(id int, url string, width, height int64, extension string) *Image {
	return &Image{
		SQLModel:  SQLModel{Id: id},
		Url:       url,
		Width:     width,
		Height:    height,
		Extension: extension,
	}
}

type Audio struct {
	SQLModel
	Url      string `json:"url" gorm:"column:url;" db:"url"`
	Format   string `json:"format" gorm:"column:format;" db:"format"`
	Duration int64  `json:"duration" gorm:"column:duration;" db:"duration"`
}

func (Audio) TableName() string { return "audio_files" }

func NewAudio(id int, url string, format string, duration int64) *Audio {
	return &Audio{
		SQLModel: SQLModel{Id: id},
		Url:      url,
		Format:   format,
		Duration: duration,
	}
}
