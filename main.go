package main


import(
	"flag"
	"fmt"
	"github.com/Bmixo/ipTools/ipDatabase"
	"os"
	"github.com/gin-gonic/gin"
)

var ip string

var address string
func init() {

	flag.StringVar(&ip, "i", "", "input you ip ")
	flag.StringVar(&address, "w", "", "web port")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "输入例示:\n")
		flag.PrintDefaults()
	}
}


func main(){
	flag.Parse()
	db := ipDatabase.NewipDataBase()

	if address==""{
		msg ,_:= db.SearchIP(ip)
		fmt.Println(string(msg))
		return
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ip/:addr", func(c *gin.Context) {
		msgs ,_:= db.SearchIP(c.Param("addr"))
		fmt.Println(c.Param("addr"))
		c.String(200, string(msgs))
	})
	r.Run(address)

}

