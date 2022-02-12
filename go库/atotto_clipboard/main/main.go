package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func main() {
	clipboard.WriteAll("日本語")
	text, _ := clipboard.ReadAll()
	fmt.Println(text)
}
