package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"ticker-server/dto"
	"time"
)

type CoinHandler struct {
	priceMap map[string]CoinInfo
}

type CoinInfo struct {
	Time         time.Time `json:"time"`
	AssetIdBase  string    `json:"asset_id_base"`
	AssetIdQuote string    `json:"asset_id_quote"`
	Rate         float64   `json:"rate"`
}

func NewCoinHandler() *CoinHandler {

	return &CoinHandler{
		priceMap: make(map[string]CoinInfo, 0),
	}
}

// curl https://rest.coinapi.io/v1/exchangerate/BTC/USD --header "X-CoinAPI-Key: B89898B1-1DFC-4D44-AB49-4D56856A3627"
func HttpGet(urlStr string, proxy string) (string, error) {
	var result string
	httpTransport := &http.Transport{}
	if proxy != "" {
		uri, _ := url.Parse(proxy)
		httpTransport.Proxy = http.ProxyURL(uri)
	}
	client := &http.Client{
		Transport: httpTransport,
	}
	resp, err := client.Get(urlStr)
	if err != nil {
		return result, err
	}
	resp.Header.Set("X-CoinAPI-Key", "B89898B1-1DFC-4D44-AB49-4D56856A3627")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	return string(body), nil
}

func (oh *CoinHandler) GetPrice(c *gin.Context) {
	res := &dto.CoinPriceResp{}

	currency := c.Param("currency")
	var reErr error
	defer func() {
		if reErr != nil {
			dto.RespFail(c, 100, reErr.Error(), nil)
			return
		} else {
			dto.RespSucess(c, res)
		}
	}()
	if currency == "" {
		reErr = errors.New("currency is empty")
		return
	}
	isQuery := false
	if _, ok := oh.priceMap[currency]; ok {
		coin := oh.priceMap[currency]
		now := time.Now()
		if now.Sub(coin.Time) > 10 {
			isQuery = true
		}
	} else {
		isQuery = true
	}
	if isQuery {
		resp, err := HttpGet("https://rest.coinapi.io/v1/exchangerate/"+currency+"/USD", "")
		if err != nil {
			reErr = errors.New("api is error")
			return
		}
		data := CoinInfo{}
		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			reErr = errors.New("json  is error")
			return
		}
		oh.priceMap[currency] = data
	}
	res.Rate = oh.priceMap[currency].Rate
	res.AssetIdBase = oh.priceMap[currency].AssetIdBase
	res.AssetIdBase = oh.priceMap[currency].AssetIdBase
	res.Time = oh.priceMap[currency].Time

	return
}