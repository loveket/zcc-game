package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	getpeople "xiuianserver/config"
)

const kapoolpath = "./config/kaPool.json"

func init() {
	rand.Seed(time.Now().Unix())
	NewKaPool()
}

type KaPool struct {
}

var ConfigMap map[int]*KaGroup
var ConfigSlice *KaPoolSlice

type KaGroup struct {
	kId         int
	WeightCount int
	KaConfig    []*Config
}
type Config struct {
	Id     int `json:"id"`
	Weight int `json:"weight"`
	Result int `json:"result"`
	End    int `json:"end"`
}

type KaPoolSlice struct {
	KaPool []*Config `json:"ka_pool"`
}

func NewKaPool() {
	file, err := os.Open(kapoolpath)
	if err != nil {
		fmt.Println("file read err", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("io read err", err)
	}
	kp := &KaPoolSlice{}
	err = json.Unmarshal(data, &kp)
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}
	file.Close()
	ConfigSlice = kp
	GetKaGroupMap()
}
func GetKaGroupMap() {
	ConfigMap = make(map[int]*KaGroup)
	for _, v := range ConfigSlice.KaPool {
		dm, ok := ConfigMap[v.Id]
		if !ok {
			dm = &KaGroup{
				kId: v.Id,
			}
			ConfigMap[v.Id] = dm
		}
		dm.WeightCount += v.Weight
		dm.KaConfig = append(dm.KaConfig, v)
	}
}
func GetRandPool(k *KaGroup) *Config {
	rNums := rand.Intn(k.WeightCount)
	rNow := 0
	for _, v := range k.KaConfig {
		rNow += v.Weight
		if rNums < rNow {
			if v.End == 1 {
				return v
			}
			pool := ConfigMap[v.Result]
			if pool == nil {
				return nil
			}
			return GetRandPool(pool)
		}
	}
	return nil
}
func (kp *KaPool) OneLuckyDraw() (kPool []string) {
	pool := ConfigMap[1000]
	for {
		config := GetRandPool(pool)
		fmt.Println(getpeople.GetPeopleName(config.Result))
		kPool = append(kPool, getpeople.GetPeopleName(config.Result))
		return kPool
	}
}
func (kp *KaPool) TenLuckyDraw() (kPool []string) {
	pool := ConfigMap[1000]
	num := 0
	for {
		config := GetRandPool(pool)
		kPool = append(kPool, getpeople.GetPeopleName(config.Result))
		num++
		if num >= 10 {
			return kPool
		}
	}
}
