package main

import (
	"github.com/Bmixo/ipTest/ipDatabase"
	"testing"
)

//func TestSearchIP(t *testing.T){
//	db := ipDatabase.NewipDataBase()
//	for i:=0;i<255;i++{ //panic测试 绝对不panic
//		for j:=0;j<255;j++{
//			for k:=0;k<255;k++{
//					ii,jj,kk:=strconv.Itoa(i),strconv.Itoa(j),strconv.Itoa(k)
//					db.SearchIP(ii+"."+jj+"."+kk+"."+"0")
//			}
//		}
//	}
//
//}

func TestSearchIPInput(t *testing.T){
	db := ipDatabase.NewipDataBase()
	_,err:= db.SearchIP("^%^*(*()()(0")
	if err==nil{
		t.Errorf("ip test fail")
	}

	_,err= db.SearchIP("1.1.1.1")
	if err!=nil{
		t.Errorf("ip test fail")
	}
}

func BenchmarkAll(b *testing.B) {
	db := ipDatabase.NewipDataBase()
	for n := 0; n < b.N; n++ {
		db.SearchIP("60.195.153.98")
	}

}

