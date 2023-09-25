package main

import (
	"flag"
	"fmt"
	"sync"
)

func main() {
	// 标志
	Tag()
	//定义命令行参数方式
	var key string
	var IP string
	var check bool
	var IPlist []string
	// file := flag.String("f", "", "IP文件")
	// out := flag.String("o", "", "输出的文件名")
	flag.StringVar(&key, "k", "", "微步的key，必须有该参数")
	flag.BoolVar(&check, "c", true, "是否自检，对本机的所有IP进行检测，格式建议用-c=xxxxxx，默认值：true，进行自检")
	//定义命令行参数方式
	flag.StringVar(&IP, "u", "", "自定义IP，该值只在对自定义IP进行检测")
	// 解析命令行参数
	flag.Parse()
	// 判断是否有key
	if key == "" {
		fmt.Println("[\033[1;31m✘\033[0m] 没有设置key")
		flag.Usage()
		return
	}
	// 判断是否有IP,并格式化
	if IP != "" {
		IPlist = ProcessIPs(IP)
	}
	// 如果自检为真，进行自检
	if check {
		// 使用netstat -ano获取IP
		ips, err := GetIP()
		if err != nil {
			fmt.Println("[\033[1;31m✘\033[0m] 获取不到本地IP信息")
		} else {
			IPlist = append(IPlist, ips...)
		}
	}
	// 判断列表有ip
	if len(IPlist) <= 0 {
		fmt.Println("[\033[1;31m✘\033[0m] 没有可查询的IP")
		return
	}
	// 去重
	IPlist, err := DeleteLocalIP(RemoveDuplicate(IPlist))
	if err != nil {
		fmt.Println("[\033[1;31m✘\033[0m] 去重失败")
		return
	}
	fmt.Println("[\033[0;38;5;214m!\033[0m] IP列表为：", IPlist)
	// 设置全局组
	var wg sync.WaitGroup
	for _, ipdata := range IPlist {
		// 增加一个goroutine标志
		wg.Add(1)
		// 判断用户是否使用自己的key
		go func(ipdata string) {
			// 关闭goroutine
			defer wg.Done()
			err = Check_IP(ipdata, key)
			if err != nil {
				fmt.Println("[\033[1;31m✘\033[0m]", err)
			}
		}(ipdata)
	}
	// 等待所有goroutine完成
	wg.Wait()
}
