package model

import (
	"math/rand"
	"strconv"
	"time"
)

var WeaponData *WeaponList

type WeaponList struct {
	Weapons []*Weapon `json:"weapons"`
}

type Weapon struct {
	WeaponUser      int        `json:"weaponuser"`      // 武器适配分类 eg:1长剑，2短剑，3巨剑
	WeaponName      string     `json:"weaponname"`      // 武器名称
	WeaponAttribute *Attribute `json:"weaponattribute"` // 武器属性
	Image           string     `json:"image,omitempty"` // 图片

}
type Attribute struct {
	AttackVal  int `json:"attackval"`
	DefenseVal int `json:"defenseval"`
	HpVal      int `json:"hpval"`
	MpVal      int `json:"mpval"`
}

func init() {
	rand.Seed(time.Now().Unix())
	ProductWeapon()
}

// 测试生成数据
func ProductWeapon() {
	w := &WeaponList{}
	for i := 1; i <= 5000; i++ {
		w1 := &Weapon{
			WeaponUser:      1,
			WeaponName:      "长剑" + strconv.Itoa(i),
			WeaponAttribute: &Attribute{AttackVal: rand.Intn(500) + 500, DefenseVal: rand.Intn(50) + 50, HpVal: rand.Intn(500) + 500, MpVal: rand.Intn(500) + 500},
			Image:           "C:\\Photos\\a.jpg",
		}
		w.Weapons = append(w.Weapons, w1)
	}
	for i := 1; i <= 5000; i++ {
		w1 := &Weapon{
			WeaponUser:      2,
			WeaponName:      "短剑" + strconv.Itoa(i),
			WeaponAttribute: &Attribute{AttackVal: rand.Intn(300) + 300, DefenseVal: rand.Intn(50) + 50, HpVal: rand.Intn(500) + 500, MpVal: rand.Intn(500) + 500},
			Image:           "",
		}
		w.Weapons = append(w.Weapons, w1)
	}
	for i := 1; i <= 5000; i++ {
		w1 := &Weapon{
			WeaponUser:      3,
			WeaponName:      "巨剑" + strconv.Itoa(i),
			WeaponAttribute: &Attribute{AttackVal: rand.Intn(800) + 800, DefenseVal: rand.Intn(100) + 100, HpVal: rand.Intn(500) + 500, MpVal: rand.Intn(500) + 500},
			Image:           "",
		}
		w.Weapons = append(w.Weapons, w1)
	}
	WeaponData = w
}
