package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"sync"
	"ticker-server/dto"
	"time"
)

type CoinHandler struct {
	lock     sync.RWMutex
	priceMap map[string]*CoinInfo
}

type CoinInfo struct {
	QueryTime    time.Time `json:-`
	Time         time.Time `json:"time"`
	AssetIdBase  string    `json:"asset_id_base"`
	AssetIdQuote string    `json:"asset_id_quote"`
	Rate         float64   `json:"rate"`
}

func NewCoinHandler() *CoinHandler {

	return &CoinHandler{
		lock:     sync.RWMutex{},
		priceMap: make(map[string]*CoinInfo, 0),
	}
}

// curl https://rest.coinapi.io/v1/exchangerate/BTC/USD --header "X-CoinAPI-Key: B89898B1-1DFC-4D44-AB49-4D56856A3627"
func HttpGet(urlStr string) (string, error) {
	var result string

	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("X-CoinAPI-Key", "B89898B1-1DFC-4D44-AB49-4D56856A3627")
	resp, err := (&http.Client{}).Do(req)

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
	oh.lock.RLock()
	if _, ok := oh.priceMap[currency]; ok {
		coin := oh.priceMap[currency]
		now := time.Now()
		if now.Sub(coin.QueryTime).Seconds() > 10 {
			isQuery = true
		}
	} else {
		isQuery = true
	}
	oh.lock.RUnlock()

	if isQuery {
		resp, err := HttpGet("https://rest.coinapi.io/v1/exchangerate/" + currency + "/USD")
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
		oh.lock.Lock()
		oh.priceMap[currency] = &data
		oh.priceMap[currency].QueryTime = time.Now()
		oh.lock.Unlock()
	}
	glog.Info("isQuery api:", isQuery, ",resp:", oh.priceMap[currency].Rate)
	res.Rate = oh.priceMap[currency].Rate
	res.AssetIdBase = oh.priceMap[currency].AssetIdBase
	res.AssetIdQuote = oh.priceMap[currency].AssetIdQuote
	res.Time = oh.priceMap[currency].Time

	return
}
