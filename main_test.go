package main

import (
	"github.com/Bmixo/ipTest/ipDatabase"
	"testing"
)




func BenchmarkAll(b *testing.B) {
	db := ipDatabase.NewipDataBase()
	for n := 0; n < b.N; n++ {
		db.SearchIP("163.177.65.160")
	}

}

