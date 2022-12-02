package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Check_IP(IP, apikey string) {
	// 关闭goroutine
	defer wg.Done()
	// 严重级别
	severity := map[string]string{"critical": "严重", "high": "高", "medium": "中", "low": "低", "info": "无危胁"}
	// 可信度
	confidence_level := map[string]string{"high": "高", "medium": "中", "low": "低"}
	// 状态码
	response_code := map[int]string{0: "成功", 1: "部分成功", 2: "没有数据", 3: "任务进行中", 4: "未发现报告", 5: "没有反病毒扫描引擎检测数据", 6: "URL 下载文件失败", 7: "URL 下载文件中", 8: "URL 下载文件上传沙箱失败", -1: "权限受限或请求出错", -2: "请求无效", -3: "请求参数缺失", -4: "超出请求限制", -5: "系统错误"}
	// scene := map[string]string{"CDN": "CDN", "University": "学校单位", "Mobile Network": "移动网络", "Unused": "已路由-未使用", "Unrouted": "已分配-未路由", "WLAN": "WLAN", "Anycast": "Anycast", "Infrastructure": "基础设施", "Internet Exchange": "交换中心", "Company": "企业专线", "Hosting": "数据中心", "Satellite Communication": "卫星通信", "Residence": "住宅用户", "Special Export": "专用出口", "Institution": "组织机构", "Cloud Provider": "云厂商"}
	// judgments := map[string]string{"Spam": "垃圾邮件", "Zombie": "傀儡机", "Scanner": "扫描", "Exploit": "漏洞利用", "Botnet": "僵尸网络", "Suspicious": "可疑", "Brute Force": "暴力破解", "Proxy": "代理", "Whitelist": "白名单", "Info": "基础信息"}

	// 请求微步的API接口
	url := fmt.Sprintf("https://api.threatbook.cn/v3/scene/ip_reputation?apikey=%v&resource=%v", apikey, IP)
	req, err1 := http.NewRequest("GET", url, nil)
	HandlingErrors(err1, "http.NewRequest")
	res, err2 := http.DefaultClient.Do(req)
	HandlingErrors(err2, "http.DefaultClient.Do")
	// 关闭请求
	defer res.Body.Close()
	body, err3 := ioutil.ReadAll(res.Body)
	HandlingErrors(err3, "http.DefaultClient.Do")
	// 定义一个Threatbook结构体存储返回的数据
	val_threat := Threatbook{}
	if err := json.Unmarshal(body, &val_threat); err != nil {
		fmt.Println("\033[31;1m[-]\033[0m 反序列化存在错误", err)
		os.Exit(1)
	}
	if (val_threat.Data[IP].Is_malicious) && (val_threat.Response_code >= 0) {
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n\033[31;1m[-]\033[0m IP：", IP, "\n-----------------------------------------------------------------\n标签类别：", val_threat.Data[IP].Tags_classes, "\n请求状态：", response_code[val_threat.Response_code], "\nIP危害级别：", severity[val_threat.Data[IP].Severity], "\n恶意的类型：", val_threat.Data[IP].Judgments, "\n运营商：", val_threat.Data[IP].Basic.Carrier, "\n国家城市：", val_threat.Data[IP].Basic.Location, "\n可信度：", confidence_level[val_threat.Data[IP].Confidence_level], "\n情报的最近更新时间：", val_threat.Data[IP].Update_time, "\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	} else if val_threat.Response_code < 0 {
		fmt.Println("\033[31;1m[-]\033[0m 请求状态：", response_code[val_threat.Response_code])
	}
}
