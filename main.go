package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/melbahja/goph"
)

//定义命令行参数
var (
	host     = flag.String("h", "", "host 默认为空")
	user     = flag.String("u", "root", "user 默认为root")
	password = flag.String("p", "", "password 默认为空")
	httpPort = flag.Int("P", 9100, "port 默认为9100")
)

func main() {
	//获取命令行参数
	flag.Parse()

	//登陆aruba交换机
	client, err := goph.New(*user, *host, goph.Password(*password))
	if err != nil {
		log.Panic(err)
	}

	defer client.Close()

	cpuinfo, err := client.Run("show cpuload")
	if err != nil {
		log.Panic(err)
	}
	meminfo, err := client.Run("show memory")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(cpuinfo), string(meminfo))

}
