package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/melbahja/goph"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//定义命令行参数
var (
	host     = flag.String("h", "", "host 默认为空")
	user     = flag.String("u", "root", "user 默认为root")
	password = flag.String("p", "", "password 默认为空")
	httpPort = flag.String("P", ":9100", "port 默认为9100")
)

//定义Prometheus Metric
var (
	cpuLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpuloadCollectedByGO",
		Help: "Current CPU Load collected by prometheus go client",
	})
)

func init() {
	prometheus.MustRegister(cpuLoad)
}

func main() {
	//获取命令行参数
	flag.Parse()

	//登陆aruba交换机
	cpuinfo, meminfo, err := sshTo(*user, *password, *host)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(cpuinfo), string(meminfo))

	cpuLoad.Set(99)
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
