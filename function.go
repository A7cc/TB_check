package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func Tag() {
	fmt.Println("\033[31;1m====================================================\033[0m\033[34;1m")
	fmt.Println("  __________             __              __  ")
	fmt.Println(" /_  __/ __ )      _____/ /_  ___  _____/ /__")
	fmt.Println("  / / / __  |_____/ ___/ __ \\/ _ \\/ ___/ //_/")
	fmt.Println(" / / / /_/ /_____/ /__/ / / /  __/ /__/ ,<")
	fmt.Println("/_/ /_____/      \\___/_/ /_/\\___/\\___/_/|_|")
	fmt.Println("\033[0m\033[31;1m====================================================\033[0m")
}

// 处理错误方法
func HandlingErrors(err error, text string) {
	if err != nil {
		fmt.Println("[\033[31;1m-\033[0m]", text, "Error：", err)
		os.Exit(1)
	}
}

// 获取IP
func GetIP() []string {
	cmd := exec.Command("netstat", "-ano")
	buf, err := cmd.CombinedOutput()
	HandlingErrors(err, "cmd.CombinedOutput")
	// IP的正则表达式
	rege, _ := regexp.Compile(`[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}`)
	// 设置IP列表
	IPlist := rege.FindAllString(string(buf), -1)
	// 去重和排除内网
	return DeleteLocalIP(RemoveDuplicate(IPlist))
}

// 去重
func RemoveDuplicate[T any](old []T) (result []T) {
	temp := map[any]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return
}

// 排除内网IP
func DeleteLocalIP(iplist []string) []string {
	// 用于输出排除内网的字符串切片
	var out []string
	// 循环遍历iplist
	for _, v := range iplist {
		// 判断元素是否是内网地址，如果不是增加到out切片
		rege, err := regexp.Compile(`^(127\.0\.0\.1)|(localhost)|(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(0\.0\.0\.0)|(172\.((1[6-9])|(2\d)|(3[01]))\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})$`)
		HandlingErrors(err, "regexp.Compile")
		matched := rege.MatchString(v)
		if !matched {
			out = append(out, v)
		}
	}
	return out
}
