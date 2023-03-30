package utils

import (
	"fmt"
	"os"
	"strings"
)

// 返回项目根目录
func GetOsPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	s := strings.Split(dir, "zcc-game")
	result := s[0] + "\\zcc-game"
	return result
}
