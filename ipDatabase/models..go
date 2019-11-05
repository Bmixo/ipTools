package ipDatabase

import "os"

type msgSuss struct {
	Ip       string `json:ip`
	Country  string `json:country`
	Province string `json:province`
	City     string `json:city`
	County   string `json:county`
	Isp      string `json:county`
	Area     string `json:area`
}
type msgError struct {
	Error string `json:error`
}
type ipDataBase struct {
	fp *os.File

	firstip int
	lastip  int
	totalip int
}

var (
	IspList = []string{
		"联通",
		"移动",
		"铁通",
		"电信",
		"长城",
		"聚友",
	}

	provinceList = map[string]bool{
		"北京":true,
		"天津":true,
		"重庆":true,
		"上海":true,
		"河北":true,
		"山西":true,
		"辽宁":true,
		"吉林":true,
		"黑龙江":true,
		"江苏":true,
		"浙江":true,
		"安徽":true,
		"福建":true,
		"江西":true,
		"山东":true,
		"河南":true,
		"湖北":true,
		"湖南":true,
		"广东":true,
		"海南":true,
		"四川":true,
		"贵州":true,
		"云南":true,
		"陕西":true,
		"甘肃":true,
		"青海":true,
		"台湾":true,
		"内蒙古":true,
		"广西":true,
		"宁夏":true,
		"新疆":true,
		"西藏":true,
		"香港":true,
		"澳门":true,
	}
	cityDirectly = []string{
		"北京",
		"天津",
		"重庆",
		"上海",
	}
)
