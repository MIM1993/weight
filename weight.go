/*
@Time : 2021/6/4 下午10:53
@Author : MuYiMing
@File : weight
@Software: GoLand
*/
package weight

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"math/rand"
	"sort"
)

const (
	//计算出当前平台的int最大值
	intSize = 32 << (^uint(0) >> 63) // cf. strconv.IntSize
	maxInt  = 1<<(intSize-1) - 1
)

//控制器
type WeightManager struct {
	//权重总和，代表生成的随机数的最大值
	Total int
	//权重节点集合
	WeightNodes []*WeightNode
	//权重值集合
	weightVals []int
}

//权重节点
type WeightNode struct {
	WeightVal int
	Item      string
}

//创建node
func NewNode(wv int, val string) *WeightNode {
	return &WeightNode{
		WeightVal: wv,
		Item:      val,
	}
}

//创建控制器
func NewManager(nodes ...*WeightNode) (*WeightManager, error) {
	return newManager(nodes...)
}

func newManager(nodes ...*WeightNode) (*WeightManager, error) {
	if len(nodes) > 0 {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].WeightVal < nodes[j].WeightVal
		})
	}

	tmpArr := make([]int, len(nodes))
	total := 0

	//生成索引
	for i, v := range nodes {
		if (maxInt - total) <= v.WeightVal {
			return nil, errors.New("sum of Choice Weights exceeds max int")
		}
		total += v.WeightVal
		tmpArr[i] = total
	}

	if total < 1 {
		//todo：权重总和不能小于1，可以在初始化配置时进行校验
		return nil, errors.New("zero Choices with Weight >= 1")
	}

	res := new(WeightManager)
	res.Total = total
	res.WeightNodes = nodes
	res.weightVals = tmpArr
	return res, nil
}


type ConfigStruct struct {
	Weights []*WeightNode
}


//根据配置文件生成控制器
func NewManagerWithCfgFile(filePath string) (*WeightManager, error) {
	var conf *ConfigStruct
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		return nil, fmt.Errorf("read config file err:%v", err)
	}
	return newManager(conf.Weights...)
}

//并发不安全
func (wm *WeightManager) Pink() (string, error) {
	n := rand.Intn(wm.Total) + 1
	idx := searchInts(wm.weightVals, n)
	if idx > len(wm.WeightNodes) {
		return "", errors.New("index over")
	}
	return wm.WeightNodes[idx].Item, nil
}

//并发安全,从函数外传入随机数种子
func (wm *WeightManager) PinkSource(rs *rand.Rand) (string, error) {
	n := rs.Intn(wm.Total) + 1
	idx := searchInts(wm.weightVals, n)
	if idx > len(wm.WeightNodes) {
		return "", errors.New("index over")
	}
	return wm.WeightNodes[idx].Item, nil
}

//从golang sort中直接抄来的
// The standard library sort.SearchInts() just wraps the generic sort.Search()
// function, which takes a function closure to determine truthfulness. However,
// since this function is utilized within a for loop, it cannot currently be
// properly inlined by the compiler, resulting in non-trivial performance
// overhead.
//
// Thus, this is essentially manually inlined version.  In our use case here, it
// results in a up to ~33% overall throughput increase for Pick().
func searchInts(a []int, x int) int {
	// Possible further future optimization for searchInts via SIMD if we want
	// to write some Go assembly code: http://0x80.pl/articles/simd-search.html
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
