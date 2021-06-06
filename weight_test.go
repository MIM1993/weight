/*
@Time : 2021/6/4 下午11:36
@Author : MuYiMing
@File : weight_test
@Software: GoLand
*/
package weight

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestWeight(t *testing.T) {
	//cfg := map[int]string{
	//	1:"1",
	//	2:"2",
	//	3:"3",
	//}
	//cfg := map[string]weightCfgNode{
	//	"hello":{weight:1,info:"hello"},
	//	"world":{weight:4,info:"world"},
	//	"mim":{weight:2,info:"mim"},
	//}
	//setWeight(cfg)
	//for k,v := range wns{
	//	fmt.Printf("weight:%d | info:%#v\n",k,v)
	//}
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!
	m, err := NewManager(
		&WeightNode{weightVal: 1, Item: "🍒"},
		&WeightNode{weightVal: 4, Item: "🍋"},
		&WeightNode{weightVal: 2, Item: "🍊"},
		&WeightNode{weightVal: 5, Item: "🥑"},
	)

	if err != nil {
		panic(fmt.Sprintf("NewManager err:%s", err))
	}

	fmt.Println("===================")
	fmt.Printf("total:%d\n",m.Total)
	for i,v := range m.WeightNodes{
		fmt.Printf("node%d: %#v\n",i,v)
	}
	fmt.Println("===================")


	fruits := make([]string, 40*18)
	for i:=0;i<len(fruits);i++{
		tmp,err := m.Pink()
		if err != nil{
			panic(fmt.Sprintf("Pink err:%s", err))
		}
		fruits[i]=tmp
	}

	fmt.Println(fruits)


	freqs:=make(map[string]int)
	for _,v := range fruits  {
		freqs[v]++
	}
	fmt.Println("================================")
	fmt.Printf("🍒:%d\n🍋:%d\n🍊:%d\n🥑:%d\n",
		freqs["🍒"],freqs["🍋"],freqs["🍊"],freqs["🥑"])

}
