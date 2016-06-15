package main

import (
	"fmt"
	"github.com/sillydong/goczd/gotime"
	"strconv"
	"time"
)

func tounixtime(timestr string) {
	unixtime, err := gotime.StrToTime(timestr, gotime.FORMAT_YYYY_MM_DD_HH_II_SS)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(unixtime)
	}
}

func parseunixtime(timestamp string) {
	unixtimestamp, err := strconv.Atoi(timestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(gotime.TimeToStr(int64(unixtimestamp), gotime.FORMAT_YYYY_MM_DD_HH_II_SS))
	}
}

func currentunixtime() {
	fmt.Println(time.Now().Unix())
}

func currentunixnanotime() {
	fmt.Println(time.Now().UnixNano())
}
