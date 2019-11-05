package main


import(
	"fmt"
	"github.com/Bmixo/ipTest/ipDatabase"
)





func main(){
	db := ipDatabase.NewipDataBase()
	msg ,_:= db.SearchIP("www.baidu.com")
	fmt.Println(string(msg))

}

