package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Image struct {
	Id        int64  `json:"id" gorm:"column:id;" db:"id"`
	Url       string `json:"url" gorm:"column:url;" db:"url"`
	Width     int64  `json:"width" gorm:"column:width;" db:"width"`
	Height    int64  `json:"height" gorm:"column:height;" db:"height"`
	Extension string `json:"extension" gorm:"column:extension;" db:"extension"`
}

func (*Image) TableName() string { return "images" }

func (img *Image) Fulfill(domain string) {
	img.Url = fmt.Sprintf("%s/%s", domain, strings.Split(img.Url, "/")[1])
}

func (img *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.WithStack(errors.New(fmt.Sprintf("Failed to unmarshal data from DB: %s", value)))
	}

	var i Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return errors.WithStack(err)
	}

	*img = i
	return nil
}

func (img *Image) Value() (driver.Value, error) {
	if img == nil {
		return nil, nil
	}
	return json.Marshal(img)
}

type Images []Image

func (i *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.WithStack(errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %s", value)))
	}

	var data []Image
	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.WithStack(err)
	}

	*i = data
	return nil
}

func (i *Images) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}
	return json.Marshal(i)
}
