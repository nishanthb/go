package main

import (
	"fmt"
	"time"
)

func main() {
	//format := "Mon Jan 2 15:04:05.000 -0700 MST 2006"
	format := "Mon Jan 2 15:04:05.000 MST 2006"
	t := time.Now()
	fmt.Println(t.Format(format))
}
