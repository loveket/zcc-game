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
	RandTest()
}

func GetRandPool(k *KaGroup) *Config {
	rand.Seed(time.Now().Unix())
	rNums := rand.Intn(k.WeightCount)
	rNow := 0
	for _, v := range k.KaConfig {
		rNow += v.Weight
		if rNums < rNow {
			return v
		}
	}
	return nil
}

func RandTest() {
	pool := ConfigMap[1000]
	for {
		config := GetRandPool(pool)
		//fmt.Println(config.End, "**", config.Id, "***", config.Result)
		if config.End == 0 {
			fmt.Println("**********", config.Result)
			for _, v := range ConfigSlice.KaPool {
				if v.Id == config.Result {
					fmt.Println(getpeople.GetPeopleName(v.Result))

				}
			}
		}
		// cm := ConfigMap[config.Result]
		// if cm == nil {
		// 	break
		// }
		time.Sleep(500 * time.Millisecond)
	}
}
