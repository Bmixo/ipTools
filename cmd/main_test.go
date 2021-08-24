package main

import (
	"github.com/Bmixo/ipTools"
	"testing"
)

func TestSearchIPInput(t *testing.T) {
	db := ipTools.Init()
	_, err := db.SearchIP("^%^*(*()()(0")
	if err == nil {
		t.Errorf("ip test fail")
	}
	_, err = db.SearchIP("1.1.1.1")
	if err != nil {
		t.Errorf("ip test fail")
	}
}

func BenchmarkAll(b *testing.B) {
	db := ipTools.Init()
	for n := 0; n < b.N; n++ {
		db.SearchIP("60.195.153.98")
	}

}
