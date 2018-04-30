package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func init() {
	commands = append(commands,
		cli.Command{
			Name:    "current",
			Aliases: []string{"c"},
			Usage:   "return current timestamp",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "nano", Usage: "return nano seconds"},
			},
			Action: current,
		}, cli.Command{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "format time",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "from",
					Usage: "format from timestamp",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "f", Usage: "format"},
					},
					Action: from,
				},
				cli.Command{
					Name:  "to",
					Usage: "format to timestamp",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "f", Usage: "format"},
					},
					Action: to,
				},
				cli.Command{
					Name:   "rule",
					Usage:  "print go time fomat rule",
					Action: rule,
				},
			},
		})
}

func current(ctx *cli.Context) {
	if ctx.Bool("nano") {
		fmt.Println(time.Now().UnixNano())
	} else {
		fmt.Println(time.Now().Unix())
	}
}

func from(ctx *cli.Context) {
	format := ctx.String("f")
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	if f, found := formats[format]; found {
		format = f
	}

	input := ctx.Args()
	if len(input) > 0 {
		for _, arg := range input {
			stamp, err := strconv.ParseInt(arg, 10, 0)
			if err != nil {
				fmt.Printf("%15s : %v\n", arg, err)
			} else {
				fmt.Printf("%15s : %v\n", arg, time.Unix(stamp, 0).Format(format))
			}
		}
	}
}

func to(ctx *cli.Context) {
	format := ctx.String("f")
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	if f, found := formats[format]; found {
		format = f
	}

	input := ctx.Args()
	if len(input) > 0 {
		sinput := strings.Join(input, " ")
		t, err := time.Parse(format, sinput)
		if err != nil {
			fmt.Printf("%15s : %v\n", sinput, err)
		} else {
			fmt.Printf("%s : %d\n", sinput, t.Unix())
		}
	}
}

var formats = map[string]string{
	"ANSIC":       "Mon Jan _2 15:04:05 2006",
	"UnixDate":    "Mon Jan _2 15:04:05 MST 2006",
	"RubyDate":    "Mon Jan 02 15:04:05 -0700 2006",
	"RFC822":      "02 Jan 06 15:04 MST",
	"RFC822Z":     "02 Jan 06 15:04 -0700", // RFC822 with numeric zone
	"RFC850":      "Monday, 02-Jan-06 15:04:05 MST",
	"RFC1123":     "Mon, 02 Jan 2006 15:04:05 MST",
	"RFC1123Z":    "Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123 with numeric zone
	"RFC3339":     "2006-01-02T15:04:05Z07:00",
	"RFC3339Nano": "2006-01-02T15:04:05.999999999Z07:00",
	"Kitchen":     "3:04PM",
	// Handy time stamps.
	"Stamp":      "Jan _2 15:04:05",
	"StampMilli": "Jan _2 15:04:05.000",
	"StampMicro": "Jan _2 15:04:05.000000",
	"StampNano":  "Jan _2 15:04:05.000000000",
}

var chars = map[string]string{
	"Year":   "2006",
	"Month":  "01",
	"Day":    "02",
	"Hour":   "15",
	"Minute": "04",
	"Second": "05",
}

func rule(ctx *cli.Context) {
	fmt.Println("internal formats")
	for name, value := range formats {
		fmt.Printf("%15s : %s\n", name, value)
	}

	fmt.Println("fomat chars")
	for name, value := range chars {
		fmt.Printf("%15s : %s\n", name, value)
	}
}
