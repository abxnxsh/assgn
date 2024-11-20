
package models

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"user_id"`
	ProductName           string    `json:"product_name"`
	ProductDescription    string    `json:"product_description"`
	ProductImages         []string  `json:"product_images"`
	CompressedProductImages []string `json:"compressed_product_images"`
	ProductPrice          float64   `json:"product_price"`
	CreatedAt             time.Time `json:"created_at"`
}

func (p *Product) UnmarshalJSON(data []byte) error {

	type Alias Product
	aux := &struct {
		CompressedProductImages json.RawMessage `json:"compressed_product_images"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if string(aux.CompressedProductImages) == "null" || len(aux.CompressedProductImages) == 0 {
		p.CompressedProductImages = nil 
	} else {
		if err := json.Unmarshal(aux.CompressedProductImages, &p.CompressedProductImages); err != nil {
			return err
		}
	}
	return nil
}
