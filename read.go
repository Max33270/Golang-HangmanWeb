package piscine

import (
	"fmt"
	"os"
)

func Read(c string) []byte {
	file, err := os.Open(c)
	if err != nil {
		fmt.Println(err)
	}
	arr := make([]byte, 1200)
	file.Read(arr)
	return arr
}
