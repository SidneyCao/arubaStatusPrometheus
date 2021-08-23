package main

import (
	"flag"
	"fmt"
)

func main() {
	var host, user, password string
	var httpPort int

	flag.StringVar(&host, "h", "", "host 默认为空")
	flag.StringVar(&user, "u", "root", "user 默认为root")
	flag.StringVar(&password, "p", "", "password 默认为空")
	flag.IntVar(&httpPort, "P", 9100, "port 默认为9100")

	flag.Parse()
	fmt.Print(host, ',', user, ',', password, ',', httpPort)
}
