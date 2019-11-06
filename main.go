package main


import(
	"flag"
	"fmt"
	"github.com/Bmixo/ipTools/ipDatabase"
	"os"
)

var ip string


func init() {

	flag.StringVar(&ip, "i", "", "input you ip ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of params:\n")
		flag.PrintDefaults()
	}
}


func main(){
	flag.Parse()
	db := ipDatabase.NewipDataBase()
	fmt.Println(ip)

	msg ,_:= db.SearchIP(ip)
	fmt.Println(string(msg))

}

