package ipTools

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var IpDataBase *MIpDataBase

func Init() *MIpDataBase {
	return IpDataBase
}
func init() {

	IpDataBase = &MIpDataBase{}

	f, err := os.OpenFile("./data/qqwry.dat", os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err.Error())
		return
	}
	IpDataBase.fp = f
	p := make([]byte, 4)
	f.Read(p)
	IpDataBase.firstip = int(binary.LittleEndian.Uint32(p))
	f.Read(p)
	IpDataBase.lastip = int(binary.LittleEndian.Uint32(p))
	IpDataBase.totalip = (IpDataBase.lastip - IpDataBase.firstip) / 7

}
func ip2int(ip string) int {
	var ipArry [4]int
	ipArryString := strings.Split(ip, ".")
	for i := 0; i < 4; i++ {
		result, err := strconv.Atoi(ipArryString[i])
		if err != nil {
			result = 1 //这里看php的函数说错误返回1
		}
		ipArry[i] = result
	}
	return (ipArry[0] << 24) | (ipArry[1] << 16) | (ipArry[2] << 8) | ipArry[3]
}

/**
 * 根据所给 IP 地址或域名返回所在地区信息
 */
func (m *MIpDataBase) Getlocationfromip(ip uint32) {
	//var msg msgSuss
	if m.fp == nil {
		return
	}

}
func long2ip() {

}

/**
 * 返回读取的3个字节的整型数
 */
func (m *MIpDataBase) getint3() int {
	//将读取的little-endian编码的3个字节转化为长整型数
	p := make([]byte, 3)
	m.fp.Read(p)
	p = append(p, 0)
	return int(binary.LittleEndian.Uint32(p))
}
func (m *MIpDataBase) getlong() int {
	//将读取的little-endian编码的3个字节转化为长整型数
	p := make([]byte, 4)
	m.fp.Read(p)
	return int(binary.LittleEndian.Uint32(p))
}

//二分查找
func (m *MIpDataBase) findIndex(ip, left, right int) int {
	if right-left <= 1 {
		return left
	}
	mid := (left + right) / 2
	midOffset := m.firstip + mid*7
	m.fp.Seek(int64(midOffset), 0)
	targetIP := m.getlong()

	if ip < targetIP {
		return m.findIndex(ip, left, mid)
	}
	return m.findIndex(ip, mid, right)

}

func (m *MIpDataBase) readString(offset int) (result []byte) {
	flag := make([]byte, 1)
	m.fp.Seek(int64(offset), 0)
	m.fp.Read(flag)
	if flag[0] == 0 { // 没有区域信息
		return []byte{}
	} else if flag[0] == 2 { //重定向
		m.fp.Seek(int64(offset+1), 0)
		offset = m.getint3()
		return m.readString(offset)
	}
	m.fp.Seek(int64(offset), 0)
	s := make([]byte, 1)

	for true {
		m.fp.Read(s)
		if s[0] == 0 {
			break
		}
		result = append(result, s[0])
	}
	return

}
func (m *MIpDataBase) getRecord(index int) (country, area []byte) {
	m.fp.Seek(int64(m.firstip+index*7+4), 0) //+4是ip偏移
	offset := m.getint3() + 4
	m.fp.Seek(int64(offset), 0)
	// 标志字节为1，表示国家和区域信息都被同时重定向
	// 标志字节为2，表示国家信息被重定向
	// 否则，表示国家信息没有被重定向

	flag := make([]byte, 1)

	m.fp.Read(flag)
	switch flag[0] {
	case 1:
		countryOffset := m.getint3()
		country = m.readString(countryOffset)
		m.fp.Read(flag)
		if flag[0] == 2 {
			area = m.readString(countryOffset + 4)
		} else {
			area = m.readString(countryOffset + len(country))
		}
	case 2:
		countryOffset := m.getint3()
		country = m.readString(countryOffset)
		area = m.readString(offset + 4)

	default:
		country = m.readString(offset)
		area = m.readString(offset + len(country))

	}
	return country, area

}
func (m *MIpDataBase) getAddr(ip string) (string, string, string, error) {
	if len(ip) == 0 {
		return "", "", "", errors.New("ip invalid")
	}
	ipAddr, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		return "", "", "", errors.New("ip invalid")
	}
	index := m.findIndex(ip2int(ipAddr.String()), 0, m.totalip)

	country, area := m.getRecord(index)

	c := simplifiedchinese.GBK.NewDecoder()
	x, err := c.Bytes(country)
	if err != nil {
		return "", "", "", errors.New("ip invalid")
	}
	xx, err := c.Bytes(area)
	if err != nil {
		return "", "", "", errors.New("ip invalid")
	}
	return ipAddr.String(), string(x), string(xx), nil
}

func (m *MIpDataBase) SearchIP(ip string) (msg []byte, err error) {

	isChina := false
	seperatorSheng := "省"
	seperatorShi := "市"
	seperatorXian := "县"
	seperatorQu := "区"

	ipaddr, country, area, errOld := m.getAddr(ip)
	if errOld != nil {
		msg, err = json.Marshal(msgError{
			Error: errOld.Error(),
		})
		return msg, errOld

	}
	msgsuss := msgSuss{
		Ip: ipaddr,
	}
	//存在 省 标志 xxx省yyyy 中的yyyy
	if strings.Contains(country, seperatorSheng) {
		isChina = true
		msgsuss.Country = "中国"
		x := strings.Split(country, seperatorSheng)
		if len(x) >= 1 { //省
			msgsuss.Province = x[0] + seperatorSheng
			x = strings.Split(x[0], seperatorShi)
			if len(x) >= 1 { //市
				msgsuss.City = x[0] + seperatorShi
				x = strings.Split(x[0], seperatorXian)
				if len(x) >= 1 { //县
					msgsuss.County = x[0] + seperatorXian
				}
			}

		}

	} else { //处理内蒙古不带省份类型的和直辖市
		if strings.Contains(country, "内蒙古") {
			isChina = true
			msgsuss.Country = "中国"
			msgsuss.Province = "内蒙古"
			x := strings.Split(country, "内蒙古")
			if len(x) >= 2 {
				msgsuss.City = x[1]
			}

		}

		for i := 0; i < len(cityDirectly); i++ {
			if strings.Contains(country, cityDirectly[i]) {
				isChina = true
				msgsuss.Country = "中国"
				x := strings.Split(country, seperatorShi)
				if len(x) >= 1 {
					msgsuss.Province = x[0] + seperatorShi
					if len(x) >= 2 {
						x := strings.Split(x[1], seperatorQu) //处理区
						if len(x) > 1 {
							msgsuss.City = x[0] + seperatorQu
						}
					}
				}
				break
			}
		}

	}
	if isChina {
		msgsuss.Area = msgsuss.Country + country + area
	} else { //国外格式
		msgsuss.Country = country
		msgsuss.Area = country + area
	}
	for i := 0; i < len(IspList); i++ {
		if strings.Contains(area, IspList[i]) {
			msgsuss.Isp = IspList[i]
			break
		}
	}
	msg, err = json.Marshal(msgsuss)
	return msg, err

}
