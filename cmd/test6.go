package main

import (
	"fmt"
	"strings"
)

func main() {
	r := strings.NewReplacer(" ", "")
	fmt.Println(r.Replace("10 10 10"))
}
