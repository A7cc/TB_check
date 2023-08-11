package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func Tag() {
	fmt.Println("\033[31;1m====================================================\033[0m\033[34;1m")
	fmt.Println("        ______  ____\033[31;1m_\033[0m\033[34;1m__\033[31;1m___\033[0m\033[34;1m_____\033[31;1m___")
	fmt.Println("       /\033[0m\033[34;1m      ||   。       \033[33;1m。\033[0m\033[34;1m   /")
	fmt.Println("      /  _    \033[31;1m|\033[0m\033[34;1m|____\033[31;1m__\033[0m\033[34;1m_     _  _/")
	fmt.Println("     \033[31;1m/\033[0m\033[34;1m  \033[31;1m/ |\033[0m\033[34;1m   |       /    /\033[31;1m( )\033[0m\033[34;1m  T00ls:微步自检工具")
	fmt.Println("    /  \033[31;1m/_\033[0m\033[34;1m_|   |      \033[31;1m/\033[0m\033[34;1m    /   ( )       ")
	fmt.Println("   /  ____    \033[31;1m|\033[0m\033[34;1m  \033[35;1m<-—+—++—+--}\033[0m\033[34;1m\033[31;1m( )\033[0m\033[34;1m____\033[31;1m/|\033[0m\033[34;1m")
	fmt.Println("  \033[31;1m/\033[0m\033[34;1m  /    |   |    /    /    ( \033[33;1m.\033[0m\033[34;1m   . \033[31;1m)\033[0m\033[34;1m")
	fmt.Println(" /\033[31;1m_\033[0m\033[34;1m_/     |_\033[31;1m__|\033[0m\033[34;1m   \033[31;1m/_\033[0m\033[34;1m___/     (\033[31;1m__\033[0m\033[34;1m__=___)  \033[31;1m❤\033[0m")
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
	return DeleteLocalIP(Unique(IPlist))
}

// 去重
func Unique(str1 []string) []string {
	// 将第一个切片元素赋值，构建一个新的字符串切片
	out := str1[:1]
	// 遍历str1的元素
	for _, word := range str1 {
		i := 0
		for ; i < len(out); i++ {
			if word == out[i] {
				break
			}
		}
		if i == len(out) {
			out = append(out, word)
		}
	}
	return out
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
