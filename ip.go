package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands,
		cli.Command{
			Name:   "myip",
			Usage:  "show my ip address",
			Action: myip,
		},
		// cli.Command{
		// 	Name:   "ipaddr",
		// 	Usage:  "show address of ip",
		// 	Action: ipaddr,
		// },
	)
}

func myip(ctx *cli.Context) {
	resp, body, err := gorequest.New().Get("http://myip.ipip.net").End()
	if err != nil {
		fmt.Println(err)
	} else {
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%d - %s\n", resp.StatusCode, resp.Status)
		} else {
			fmt.Println(strings.TrimSpace(body))
		}
	}
}

// func ipaddr(ctx *cli.Context) {

// }
