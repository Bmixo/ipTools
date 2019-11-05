package ipDatabase

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var db *ipDataBase
func NewipDataBase() *ipDataBase {
	return db
}
func init() {

	db = &ipDataBase{}

	workSpace, err := os.Getwd()
	f, err := os.OpenFile(filepath.Join(workSpace, "ipDatabase/qqwry.dat"), os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err.Error())
		return
	}
	db.fp = f
	p := make([]byte, 4)
	f.Read(p)
	db.firstip = int(binary.LittleEndian.Uint32(p))
	f.Read(p)
	db.lastip = int(binary.LittleEndian.Uint32(p))
	db.totalip = (db.lastip - db.firstip) / 7




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
func (db *ipDataBase) Getlocationfromip(ip uint32) {
	//var msg msgSuss
	if db.fp == nil {
		return
	}

}
func long2ip() {

}

/**
 * 返回读取的3个字节的整型数
 */
func (db *ipDataBase) getint3() int {
	//将读取的little-endian编码的3个字节转化为长整型数
	p := make([]byte, 3)
	db.fp.Read(p)
	p = append(p, 0)
	return int(binary.LittleEndian.Uint32(p))
}
func (db *ipDataBase) getlong() int {
	//将读取的little-endian编码的3个字节转化为长整型数
	p := make([]byte, 4)
	db.fp.Read(p)
	return int(binary.LittleEndian.Uint32(p))
}

//二分查找
func (db *ipDataBase) findIndex(ip, left, right int) int {
	if right-left <= 1 {
		return left
	}
	mid := (left + right) / 2
	midOffset := db.firstip + mid*7
	db.fp.Seek(int64(midOffset), 0)
	targetIP := db.getlong()

	if ip < targetIP {
		return db.findIndex(ip, left, mid)
	}
	return db.findIndex(ip, mid, right)

}


func (db *ipDataBase) readString(offset int)  (result []byte){
	flag:=make([]byte,1)
	db.fp.Seek(int64(offset),0)
	db.fp.Read(flag)
	if flag[0]==0{// 没有区域信息
		return []byte{}
	}else if flag[0]==2 { //重定向
		db.fp.Seek(int64(offset+1),0)
		offset =db.getint3()
		return db.readString(offset)
	}
	db.fp.Seek(int64(offset),0)
	s:= make([]byte,1)

	for true{
		db.fp.Read(s)
		if s[0]==0{
			break
		}
		result=append(result,s[0])
	}
	return


}
func (db *ipDataBase) getRecord(index int)  (country ,area[]byte){
	db.fp.Seek( int64(db.firstip+ index*7+4),0) //+4是ip偏移
	offset :=db.getint3()+4
	db.fp.Seek(int64(offset),0)
	// 标志字节为1，表示国家和区域信息都被同时重定向
	// 标志字节为2，表示国家信息被重定向
	// 否则，表示国家信息没有被重定向

	flag:=make([]byte,1)

	db.fp.Read(flag)
	switch(flag[0]){
	case 1:
		countryOffset := db.getint3()
		country=db.readString(countryOffset)
		db.fp.Read(flag)
		if flag[0] ==2{
			area=db.readString(countryOffset+4)
		}else{
			area=db.readString(countryOffset+len(country))
		}
	case 2:
		countryOffset := db.getint3()
		country=db.readString(countryOffset)
		area= db.readString(offset+4)

	default:
		country=db.readString(offset)
		area= db.readString(offset+len(country))

	}
	return country,area


}
func (db *ipDataBase) getAddr(ip string) (string,string ,string,error) {
	ipAddr, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		return "","","",errors.New("ip invalid")
	}
	index := db.findIndex(ip2int(ipAddr.String()), 0, db.totalip)

	country,area := db.getRecord(index)

	c := simplifiedchinese.GBK.NewDecoder()
	x,err := c.Bytes(country)
	if err != nil {
		return "","","",errors.New("ip invalid")
	}
	xx,err := c.Bytes(area)
	if err != nil {
		return "","","",errors.New("ip invalid")
	}
	return ipAddr.String(),string(x),string(xx),nil
}




func (db *ipDataBase) SearchIP(ip string) (msg []byte,err error){
	ipaddr,country,area,err:= db.getAddr(ip)
	if err!=nil{
		msg,err = json.Marshal(msgError{
			Error:err.Error(),
		})
		return msg,err

	}
	msgsuss := msgSuss{
		Ip:ipaddr,
	}
	if _,ok := provinceList[country]  ;ok { //map查询高效
		msgsuss.Country="中国"
		msgsuss.Area = "中国"+country+area

	}else{
		msgsuss.Country=country
		msgsuss.Area = country+area
	}
	for i:=0;i<len(IspList);i++{
		if strings.Contains(area,IspList[i]){
			msgsuss.Isp=IspList[i]
			break
		}
	}
	msg,err= json.Marshal(msgsuss)
	return msg,err

}
