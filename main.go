package main

import (
	"flag"
	"github.com/golang/glog"
	"ticker-server/config"
	"ticker-server/server"
	"time"
)

/*
完成如下小任务，请当成工作场景，你可以查阅任何资料，使用任何工具，没有任何限制：


curl https://rest.coinapi.io/v1/exchangerate/BTC/USD --header "X-CoinAPI-Key: B89898B1-1DFC-4D44-AB49-4D56856A3627"

这是一个数字货币汇率的第三方API，path中有一对数字货币的符号，比如现在是 BTC 的美元价格。

1. 在golang中请求这个 API，设计一个 struct，把结果存入其中，并打印出来。
2. 实现一个 http 服务器，框架不限，简化这个API为 /price/:currency ，返回任何一个数字货币的美元价格。
3. 为了更快的响应，和节约费用，为这个 API 增加一个缓存，可以让已经请求过的结果在 10秒内不用再调用上游 API 。
*/

var (
	ConfigURL string
)

func init() {
	flag.StringVar(&ConfigURL, "config", "./config.yaml", "config url")
	flag.Parse()
	config.InitConfig(ConfigURL)
}

// go run main.go -alsologtostderr
func main() {
	tickerGlog := time.NewTicker(time.Second * 5)
	go func() {
		for range tickerGlog.C {
			glog.Flush()
		}
	}()
	glog.Info("run")

	server.StartServer()
}
