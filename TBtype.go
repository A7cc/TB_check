// TBtype.go
package main

// 关于asn信息结构体
type Asn struct {
	Rank   int    `json:"rank"`
	Info   string `json:"info"`
	Number int    `json:"number"`
}

// ip对应的位置信息
type Location struct {
	// 国家
	Country string
	// 省
	Province string
	// 城市
	City string
	// 经度
	Lng string
	// 纬度
	Lat string
	// 国家代码
	Country_code string
}

// IP附属信息
type Basic struct {
	// 运营商
	Carrier string
	// ip对应的位置信息
	Location Location
}

// 相关攻击团伙或安全事件信息
type Tags_classes struct {
	// 标签类别
	Tags []string
	// 具体的攻击团伙或安全事件标签
	Tags_type string
}

// 关于IP的信息结构体
type IP struct {
	// 严重级别
	Severity string `json:"severity"`
	// 恶意的类型
	Judgments []string `json:"judgments"`
	// 相关攻击团伙或安全事件信息
	Tags_classes []Tags_classes `json:"tags_classes"`
	// IP附属信息
	Basic Basic
	// asn信息
	Asn Asn `json:"asn"`
	// 应用场景
	Scene string `json:"scene"`
	// 可信度
	Confidence_level string `json:"confidence_level"`
	// 是否为恶意IP
	Is_malicious bool `json:"is_malicious"`
	// 情报的最近更新时间
	Update_time string `json:"update_time"`
}

// 关于Threatbook的信息结构体
type Threatbook struct {
	// Threatbook返回的信息
	Data map[string]IP
	// 状态码
	Response_code int `json:"response_code"`
	// 状态码信息
	Verbose_msg string `json:"verbose_msg"`
}
