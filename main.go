package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/melbahja/goph"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cpuValid = regexp.MustCompile(`[a-z]+\s*(.*?)%`)
var memValid = regexp.MustCompile(`[0-9]+`)

//定义命令行参
var (
	host     = flag.String("h", "", "host 默认为空")
	user     = flag.String("u", "root", "user 默认为root")
	password = flag.String("p", "", "password 默认为空")
	httpPort = flag.String("P", ":9100", "port 默认为9100")
)

//定义Prometheus Metric
var (
	cpuLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpuloadCollectedByGO",
			Help: "Current CPU Load collected by prometheus go client",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(cpuLoad)
}

func main() {
	//获取命令行参数
	flag.Parse()

	go func() {
		for {
			//登陆aruba交换机
			cpuinfo, meminfo, err := sshTo(*user, *password, *host)
			if err != nil {
				log.Panic(err)
			}
			fmt.Println(string(meminfo))
			cpuSlice := cpuValid.FindAllStringSubmatch(cpuinfo, -1)
			memSlice := memValid.FindAllString(meminfo, -1)

			fmt.Println(memSlice)

			cpuUser, err := strconv.ParseFloat(cpuSlice[0][1], 64)
			if err != nil {
				log.Panic(err)
			}

			cpuSystem, err := strconv.ParseFloat(cpuSlice[1][1], 64)
			if err != nil {
				log.Panic(err)
			}

			cpuIdle, err := strconv.ParseFloat(cpuSlice[2][1], 64)
			if err != nil {
				log.Panic(err)
			}

			cpuLoad.With(prometheus.Labels{"type": "user"}).Set(cpuUser)
			cpuLoad.With(prometheus.Labels{"type": "system"}).Set(cpuSystem)
			cpuLoad.With(prometheus.Labels{"type": "idle"}).Set(cpuIdle)
			time.Sleep(10 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Panic(http.ListenAndServe(*httpPort, nil))

}

func sshTo(user string, password string, host string) (string, string, error) {
	client, err := goph.New(user, host, goph.Password(password))
	if err != nil {
		return "", "", err
	}

	defer client.Close()

	cpuinfo, err := client.Run("show cpuload")
	if err != nil {
		return "", "", err
	}
	meminfo, err := client.Run("show memory")
	if err != nil {
		return string(cpuinfo), "", err
	}
	return string(cpuinfo), string(meminfo), nil
}
