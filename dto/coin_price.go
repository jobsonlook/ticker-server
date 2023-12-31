package dto

import "time"

type CoinPriceReq struct {
}

type CoinPriceResp struct {
	Time         time.Time `json:"time"`
	AssetIdBase  string    `json:"asset_id_base"`
	AssetIdQuote string    `json:"asset_id_quote"`
	Rate         float64   `json:"rate"`
}
