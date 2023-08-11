package main

import (
	"flag"
	"fmt"
	"sync"
)

// 设置全局组
var wg sync.WaitGroup

func main() {
	// 标志
	Tag()
	//定义命令行参数方式
	var key string
	var IP string
	var check bool
	// file := flag.String("f", "", "IP文件")
	// out := flag.String("o", "", "输出的文件名")
	flag.StringVar(&key, "k", "", "微步的key，必须有该参数")
	flag.BoolVar(&check, "c", true, "是否自检，对本机的所有IP进行检测，格式建议用-c=xxxxxx，默认值：true，进行自检")
	//定义命令行参数方式
	flag.StringVar(&IP, "u", "", "自定义IP，该值只在对自定义IP进行检测")
	// 解析命令行参数
	flag.Parse()
	// 如果自检为真，进行自检
	if check && (key != "") {
		// 使用netstat -ano获取IP
		IPlist := GetIP()
		if IP != "" {
			IPlist = append(IPlist, IP)
		}
		fmt.Println("[\033[33;1m*\033[0m] IP列表为：", IPlist)
		for _, v := range IPlist {
			// 增加一个goroutine标志
			wg.Add(1)
			// 判断用户是否使用自己的key
			go func() {
				Check_IP(v, key)
			}()
		}
	} else if (IP != "") && (key != "") {
		// 增加一个goroutine标志
		wg.Add(1)
		// 用户是否检查IP和使用key
		Check_IP(IP, key)
	} else {
		fmt.Println("[\033[31;1m-\033[0m] 没有输入参数！")
	}
	// 等待所有goroutine完成
	wg.Wait()
}
