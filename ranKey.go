package main

import (
	"fmt"
	"math/rand"
	// "strconv"
	"time"
)

func main() {

	var listData [32]int

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 32; i++ {
		x := rand.Intn(255)

		listData[i] = x
	}

	var mData string

	mData = "{"

	for i := 0; i < len(listData); i++ {
		// fmt.Println(listData[i] + "\n")

		mData += fmt.Sprintf("0x%02x", listData[i])

		if i < len(listData)-1 {
			mData += ", "
		} else {
			mData += "}"
		}

	}

	fmt.Println(mData)
}
