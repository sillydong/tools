package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		help(args)
	} else {
		action := args[1]
		switch action {
		case "fromtime":
			if checkargs(3, args) {
				parseunixtime(args[2])
			}
		case "totime":
			if checkargs(3, args) {
				tounixtime(args[2])
			}
		case "unixtime":
			if checkargs(2, args) {
				currentunixtime()
			}
		case "unixnano":
			if checkargs(2, args) {
				currentunixnanotime()
			}
		case "md5":
			if checkargs(3, args) {
				md5(args[2])
			}
		case "base64encode":
			if checkargs(3, args) {
				base64_encode(args[2])
			}
		case "base64decode":
			if checkargs(3, args) {
				base64_decode(args[2])
			}
		case "urlencode":
			if checkargs(3, args) {
				urlencode(args[2])
			}
		case "urldecode":
			if checkargs(3, args) {
				urldecode(args[2])
			}
		case "ip":
			if checkargs(3, args) {
				ip(args[2])
			}
		default:
			help(args)
		}
	}
}

func help(args []string) {
	fmt.Println("run: ", args[0], " command args")
	fmt.Println(`
commands:
	fromtime timestamp     parse timestamp to string
	totime timestring      parse time string to timestamp
	unixtime               get current unix timestamp
	unixnano               get current unix nano timestamp
	md5 string             get md5
	base64encode string    get base64 ecnode
	base64decode string    get base64 decode
	urlencode string       get urlencode
	urldecode string       get urldecode
	ip 0.0.0.0             get location for ip address
`)
}

func checkargs(lenth int, args []string) bool {
	if len(args) == lenth {
		return true
	}
	help(args)
	return false
}
