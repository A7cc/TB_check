package main

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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

// 获取IP
func GetIP() ([]string, error) {
	cmd := exec.Command("netstat", "-ano")
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	// IP的正则表达式
	rege, _ := regexp.Compile(`[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}`)
	// 设置IP列表
	IPlist := rege.FindAllString(string(buf), -1)
	// 去重和排除内网
	return IPlist, nil
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
func DeleteLocalIP(iplist []string) ([]string, error) {
	// 用于输出排除内网的字符串切片
	var out []string
	// 循环遍历iplist
	for _, v := range iplist {
		// 判断元素是否是内网地址，如果不是增加到out切片
		rege, err := regexp.Compile(`^(127\.0\.0\.1)|(localhost)|(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(0\.0\.0\.0)|(172\.((1[6-9])|(2\d)|(3[01]))\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})$`)
		if err != nil {
			return nil, err
		}
		matched := rege.MatchString(v)
		if !matched {
			out = append(out, v)
		}
	}
	return out, nil
}

// 处理ip
func ProcessIPs(ips string) (hostlist []string) {
	// 判断是否有逗号
	if strings.Contains(ips, ",") {
		// 如果有逗号将其划分多个IP表
		IPList := strings.Split(ips, ",")
		// 存放IP用与hosts变量一致
		var ips []string
		// 循环处理IP表
		for _, ip := range IPList {
			ips = parseIP(ip)
			hostlist = append(hostlist, ips...)
		}
	} else {
		hostlist = parseIP(ips)
	}
	return RemoveDuplicate(hostlist)
}

// 根据用户给出的ip形式进行分类
func parseIP(ip string) []string {
	reg := regexp.MustCompile(`[a-zA-Z]+`)
	switch {
	case strings.HasSuffix(ip, "/8"):
		// 扫描/8时，由于A段太多了，只扫网关和随机IP，避免扫描过多IP
		return parseIP8(ip)
	case strings.Contains(ip, "/"):
		// 解析 /24 /16等
		return parseIP2(ip)
	case reg.MatchString(ip):
		// 域名用lookup获取ip
		host, err := net.LookupHost(ip)
		if err != nil {
			return nil
		}
		return host
	case strings.Contains(ip, "-"):
		// 处理192.168.1.1-192.168.1.100或者192.168.1.1-24
		return parseIP1(ip)
	default:
		// 处理单个ip
		testIP := net.ParseIP(ip)
		if testIP == nil {
			return nil
		}
		return []string{ip}
	}
}

// 把 192.168.x.x/xx 转换成IP列表
func parseIP2(host string) (hosts []string) {
	// 使用 net.ParseCIDR() 方法解析给定的网段，返回网段的 IP 地址和子网掩码
	// 检查给定的网段是否正确
	ipone, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		return
	}
	// 把 192.168.x.x/xx 转换成 192.168.x.x-192.168.x.x 并转成IP列表
	hosts = parseIP1(IPRange(ipone.String(), ipNet))
	return
}

// 解析ip段: 192.168.111.1-255，192.168.111.1-192.168.112.255
func parseIP1(ip string) []string {
	// 如果有逗号将其划分多个
	IPRangelist := strings.Split(ip, "-")
	// 确认该IP格式是否为正确IP
	testIP := net.ParseIP(IPRangelist[0])
	// 创建一个存储所有IP列表
	var allIP []string
	// 通过len函数来确认IPRangelist[1]是192.168.1.255形式还是数字形式
	if len(IPRangelist[1]) < 4 {
		// 处理数字形式
		// 将字符串转成数字
		Range, err := strconv.Atoi(IPRangelist[1])
		// 判断合理性
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		// 分离IP
		SplitIP := strings.Split(IPRangelist[0], ".")
		// 转换为数字
		ip1, err1 := strconv.Atoi(SplitIP[3])
		// 拼接
		PrefixIP := SplitIP[0] + "." + SplitIP[1] + "." + SplitIP[2]
		// 判断合理性
		if ip1 > Range || err1 != nil {
			return nil
		}
		// 循环拼接IP
		for i := ip1; i <= Range; i++ {
			allIP = append(allIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		// 处理192.168.1.255形式
		// 分离IP
		SplitIP1 := strings.Split(IPRangelist[0], ".")
		SplitIP2 := strings.Split(IPRangelist[1], ".")
		// 判断合理性
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		// 用于存放起始IP和结束IP列表
		start, end := [4]int{}, [4]int{}
		// 循环读取4段IP
		for i := 0; i < 4; i++ {
			// 转换为数字
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			// 判断合理性
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			// 添加到起始IP和结束IP列表
			start[i], end[i] = ip1, ip2
		}
		// 通过移位运算，将地址改为数字的形式
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		// 通过循环将数字转成地址
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			allIP = append(allIP, ip)
		}
	}
	return allIP
}

// 获取把 192.168.x.x/xx 转换成 192.168.x.x-192.168.x.x的起始IP、结束IP
func IPRange(start string, c *net.IPNet) string {
	// 16进制子网掩码
	mask := c.Mask
	// 创建一个net.ip类型
	bcst := make(net.IP, len(c.IP))
	// 将dreams值给bcst
	copy(bcst, c.IP)
	// 获取结束IP
	for i := len(mask) - 1; i >= 0; i-- {
		bcst[i] = c.IP[i] | ^mask[i]
	}
	end := bcst.String()
	// 返回用-表示的ip段,192.168.1.1-192.168.255.255
	return fmt.Sprintf("%s-%s", start, end)
}

// 处理B段IP
func parseIP8(ip string) []string {
	// 去掉最后的/8
	realIP := ip[:len(ip)-2]
	// net.ParseIP 这个方法用来检查 ip 地址是否正确，如果不正确，该方法返回 nil
	testIP := net.ParseIP(realIP)
	// 判断该IP是否为正常IP
	if testIP == nil {
		return nil
	}
	// 获取IP的头部
	IP8head := strings.Split(ip, ".")[0]
	// 存放B段IP
	var allIP []string
	// 构造B段的随机IP表
	for a := 0; a <= 255; a++ {
		for b := 0; b <= 255; b++ {
			// 一般情况下网关为1或者254
			allIP = append(allIP, fmt.Sprintf("%s.%d.%d.%d", IP8head, a, b, 1))
			allIP = append(allIP, fmt.Sprintf("%s.%d.%d.%d", IP8head, a, b, RandInt(2, 80)))
			allIP = append(allIP, fmt.Sprintf("%s.%d.%d.%d", IP8head, a, b, RandInt(81, 170)))
			allIP = append(allIP, fmt.Sprintf("%s.%d.%d.%d", IP8head, a, b, RandInt(171, 253)))
			allIP = append(allIP, fmt.Sprintf("%s.%d.%d.%d", IP8head, a, b, 254))
		}
	}
	return allIP
}

// 随机数
func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
