package pubsub

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Image struct {
	ID            uint32 `json:"img_id,omitempty" bson:"img_id,omitempty"`
	FakeID        *UID   `json:"id,omitempty" bson:"-"`
	Url           string `json:"url" bson:"url"`
	FileName      string `json:"file_name,omitempty" bson:"file_name,omitempty"`
	OriginWidth   int    `json:"org_width" bson:"org_width"`
	OriginHeight  int    `json:"org_height" bson:"org_height"`
	OriginUrl     string `json:"org_url" bson:"org_url"`
	CloudName     string `json:"cloud_name,omitempty" bson:"cloud_name"`
	CloudId       string `json:"cloud_id,omitempty" bson:"cloud_id"`
	DominantColor string `json:"dominant_color" bson:"dominant_color"`
	RequestId     string `json:"request_id,omitempty" bson:"-"`
	FileSize      uint32 `json:"file_size,omitempty" bson:"-"`
}

func (i *Image) HideSomeInfo() *Image {
	if i != nil {
		// i.CloudID = ""
		i.CloudId = ""
	}

	return i
}

func (i *Image) Fulfill(domain string) {
	i.Url = fmt.Sprintf("%s%s", domain, i.CloudId)
}

// This method for mapping Image to json data type in sql
func (i *Image) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Image from date data type in sql
func (i *Image) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	if err := json.Unmarshal(v, i); err != nil {
		return err
	}
	return nil
}

type Images []Image

// This method for mapping Images to json array data type in sql
func (is *Images) Value() (driver.Value, error) {
	if is == nil {
		return nil, nil
	}

	b, err := json.Marshal(is)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Images from json array type in sql
func (is *Images) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	var imgs []Image

	if err := json.Unmarshal(v, &imgs); err != nil {
		return err
	}

	*is = Images(imgs)
	return nil
}
