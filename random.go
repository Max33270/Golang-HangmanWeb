package piscine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func Random(u string) string {
	rand.Seed(time.Now().UnixNano())
	t := ""
	s := ""
	var a int
	if u == "EASY" {
		a = rand.Intn(37) + 1
		t = "words.txt"
	} else if u == "NORMAL" {
		a = rand.Intn(23) + 1
		t = "words2.txt"
	} else if u == "HARD" {
		a = rand.Intn(24) + 1
		t = "words3.txt"
	} else {
		fmt.Println("Erreur")
		os.Exit(1)
	}
	file, err := os.Open(t)
	if err != nil {
		fmt.Println(err)
	}
	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		if i == a {
			s = scanner.Text()
		}
	}
	return s
}
