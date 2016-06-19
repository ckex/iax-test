package main

import (
	logger "github.com/alecthomas/log4go"
	"time"
	"flag"
	"math/rand"
	"runtime"
	"math"
	"iax-test/service"
	"encoding/json"
	"iax-test/models"
	"strings"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var coroutines = flag.Int("coroutines", 1, "并发数")
var reqnum = flag.Int("reqnum", 1, "请求数")
var interval = flag.Int("interval", 500, "每次请求间隔")
var url = flag.String("url", "http://127.0.0.1:9090/v1/bid", "请求目录地址")

func init() {
	logger.LoadConfiguration("./conf/log4go.xml")
}

func main() {
	defer logger.Close()
	flag.Parse()
	cpus := runtime.NumCPU()

	logger.Info("CPU核心数=%d,\t并发数=%d,\t单线程请求数=%d,\t每次请求间隔时间=%d,\t请求地址=%s", cpus, *coroutines, *reqnum, *interval, *url)

	timeStart := time.Now().UnixNano()
	runtime.GOMAXPROCS(cpus)

	coroutinesNum := *coroutines
	subCoroutinesNum := *reqnum
	count := coroutinesNum * subCoroutinesNum

	ch := make(chan service.IaxResponse, count)
	defer close(ch)

	for i := 0; i < coroutinesNum; i++ {
		go func(index int) {
			// Start Coroutines
			for j := 0; j < subCoroutinesNum; j++ {
				st := time.Now().UnixNano()
				req := service.IaxRequest{Id:index, SubId:j, Descr:"haha", Url:*url}
				res := Request(req)
				end := time.Now().UnixNano()
				useTimeNano := end - st
				if *interval > 0 && useTimeNano < int64(*interval * 1000 * 1000) {
					sleepTime := int64(*interval * 1000 * 1000) - useTimeNano
					time.Sleep(time.Nanosecond * time.Duration(sleepTime))
				}
				end = time.Now().UnixNano()
				useTime := float64(end - st) / 1000.00 / 1000.00
				res.UseTime = float32(round(useTime, 2))
				ch <- res
			}
		}(i)
	}

	resultList := make([]service.IaxResponse, count)
	next := true
	var index int = 0
	for next {
		select {
		case result := <-ch:
			resultList[index] = result
			index += 1
			if index == count {
				next = false
			}
		}

	}

	timeEnd := time.Now().UnixNano()
	useTime := timeEnd - timeStart

	minUseTime := resultList[0] // 最小用时
	maxUseTime := minUseTime // 最大用时
	success := 0 //成功总数
	successMinUseTime := minUseTime // 状态成功最小用时
	successMaxUseTime := successMinUseTime // 状态成功最大用时

	successResult := 0
	noAdsCount := 0
	for index, value := range resultList {
		if value.StatusCode == 200 {
			success += 1
			if value.UseTime < successMinUseTime.UseTime {
				successMinUseTime = value
			}
			if value.UseTime > successMaxUseTime.UseTime {
				successMaxUseTime = value
			}
			var result models.Result
			bodyText := value.Body
			bodyText = strings.Replace(bodyText, "window.admaxADMAX_1.serve.callback(", "", -1)
			bodyText = strings.Replace(bodyText, ")", "", -1)
			bodyText = strings.Replace(bodyText, "\n", "", -1)
			bodyText = strings.Replace(bodyText, "\t", "", -1)
			err := json.Unmarshal([]byte(bodyText), &result)
			if err != nil {
				logger.Error("%+v", err)
			}
			value.Result = result
			if result.Success {
				successResult += 1
			}
			if len(result.Ads[0].Link) < 100 || strings.HasPrefix(result.Ads[0].Link, "http") {
				noAdsCount += 1
				logger.Debug(" index=%d use=%.2f-ms\t ----> %s, %s\n", index, value.UseTime,result.Ads[0].BidId,result.Ads[0].Link)
			}

		}
		if value.UseTime < minUseTime.UseTime {
			minUseTime = value
		}
		if value.UseTime > maxUseTime.UseTime {
			maxUseTime = value
		}
		//logger.Info("%d,%+v", index, value)
	}

	successRate := float64(success) / float64(count) * 100.0 // 成功比例
	successAdRate := float64(count - noAdsCount) / float64(count) * 100.0 //DSP参与广告竟价的比例
	//successRate = round(successRate, 2)

	allUseTime := float64(useTime) / 1000.00 / 1000.00 // 总用时

	avgUseTime := allUseTime / float64(count)// 平均用时

	allTimeSecond := float64(allUseTime) / 1000.00
	if allTimeSecond <= 0 {
		allTimeSecond = 1
	}
	qps := float64(count) / allTimeSecond
	if qps < 1 && count > 0 {
		qps = 1
	}

	logger.Info("\n\tGame over, All use Time=%.2f-ms, avg use time=%.2f-ms, min time=%.2f, max time=%.2f, success min time=%.2f, success max time=%.2f, request count=%d, success count=%d, result success count=%d, no ads count=%d, failure=%d, success rate=%.2f%s, success ad rate=%.2f%s, QPS=%d \n", allUseTime, avgUseTime, minUseTime.UseTime, maxUseTime.UseTime, successMinUseTime.UseTime, successMaxUseTime.UseTime, count, success, successResult, noAdsCount, (count - success), successRate, "%", successAdRate, "%", int(qps))
}

func round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

func Request(req service.IaxRequest) service.IaxResponse {
	return service.Get(req)
}
