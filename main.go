package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/melbahja/goph"
)

func main() {
	//获取命令行参数
	var host, user, password string
	var httpPort int

	flag.StringVar(&host, "h", "", "host 默认为空")
	flag.StringVar(&user, "u", "root", "user 默认为root")
	flag.StringVar(&password, "p", "", "password 默认为空")
	flag.IntVar(&httpPort, "P", 9100, "port 默认为9100")
	flag.Parse()

	//登陆aruba交换机
	client, err := goph.New(user, host, goph.Password(password))
	if err != nil {
		log.Panic(err)
	}

	defer client.Close()

	out, err := client.Run("show cpuinfo")
	if err != nil {
		log.Panic(err)
	}
	res, err := ioutil.ReadAll(out)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(res)

}
