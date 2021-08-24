package main

import (
	"fmt"
	"github.com/Bmixo/ipTools"
)

func main() {
	db := ipTools.Init()
	result, err := db.SearchIP("60.195.153.98")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))
}
