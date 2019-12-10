package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	commands = append(commands,
		cli.Command{
			Name:  "md5",
			Usage: "md5 encode",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "file", Usage: "generate file md5"},
			},
			Action: md5encode,
		},
		cli.Command{
			Name:  "crc32",
			Usage: "crc32 generate",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "file", Usage: "generate file crc32"},
			},
			Action: crc32encode,
		},
		cli.Command{
			Name:  "sha1",
			Usage: "sha1 generate",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "file", Usage: "generate file sha1"},
			},
			Action: sha1encode,
		},
		cli.Command{
			Name:  "base64",
			Usage: "base64 encode/decode",
			Subcommands: []cli.Command{
				cli.Command{
					Name:    "encode",
					Aliases: []string{"e"},
					Usage:   "base64 encode",
					Action:  base64encode,
				},
				cli.Command{
					Name:    "decode",
					Aliases: []string{"d"},
					Usage:   "base64 decode",
					Action:  base64decode,
				},
			},
		},
		cli.Command{
			Name:  "url",
			Usage: "url encode/decode",
			Subcommands: []cli.Command{
				cli.Command{
					Name:    "encode",
					Aliases: []string{"e"},
					Usage:   "url encode",
					Action:  urlencode,
				},
				cli.Command{
					Name:    "decode",
					Aliases: []string{"d"},
					Usage:   "url decode",
					Action:  urldecode,
				},
			},
		},
		cli.Command{
			Name:   "htpasswd",
			Usage:  "generate htpasswd",
			Action: htpasswd,
		},
	)
}

func md5encode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "md5")
	} else {
		if ctx.Bool("f") {
			// filename
			f, err := os.Open(ctx.Args()[0])
			defer f.Close()
			if err != nil {
				fmt.Println(err)
			} else {
				h := md5.New()
				if _, err := io.Copy(h, f); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%x\n", h.Sum(nil))
				}
			}
		} else {
			h := md5.New()
			io.WriteString(h, ctx.Args()[0])
			fmt.Printf("%x\n", h.Sum(nil))
		}
	}
}

func crc32encode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "crc32")
	} else {
		crc32q := crc32.MakeTable(0xD5828281)
		if ctx.Bool("f") {
			// filename
			body, err := ioutil.ReadFile(ctx.Args()[0])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%08x\n", crc32.Checksum(body, crc32q))
			}
		} else {

			fmt.Printf("%08x\n", crc32.Checksum([]byte(ctx.Args()[0]), crc32q))
		}
	}
}

func sha1encode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "sha1")
	} else {
		if ctx.Bool("f") {
			// filename
			f, err := os.Open(ctx.Args()[0])
			defer f.Close()
			if err != nil {
				fmt.Println(err)
			} else {
				h := sha1.New()
				if _, err := io.Copy(h, f); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%x\n", h.Sum(nil))
				}
			}
		} else {
			h := sha1.New()
			io.WriteString(h, ctx.Args()[0])
			fmt.Printf("%x\n", h.Sum(nil))
		}
	}
}

func base64encode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "encode")
	} else {
		for _, str := range ctx.Args() {
			fmt.Printf("%s : %s\n", str, base64.StdEncoding.EncodeToString([]byte(str)))

		}
	}
}

func base64decode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "decode")
	} else {
		for _, str := range ctx.Args() {
			dstr, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				fmt.Printf("%s : %v\n", str, err)
			} else {
				fmt.Printf("%s : %s\n", str, dstr)
			}
		}
	}
}

func urlencode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "encode")
	} else {
		for _, u := range ctx.Args() {
			fmt.Printf("%s :\n\t%s\n", u, url.QueryEscape(u))
		}
	}
}

func urldecode(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "decode")
	} else {
		for _, u := range ctx.Args() {
			ud, err := url.QueryUnescape(u)
			if err != nil {
				fmt.Printf("%s :\n\t%v\n", u, err)
			} else {
				fmt.Printf("%s :\n\t%s\n", u, ud)
			}
		}
	}
}

func htpasswd(ctx *cli.Context) {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "htpasswd")
	} else {
		for _, pass := range ctx.Args() {
			sha := hashSha(pass)
			fmt.Printf("sha : {SHA}%s\n", sha)
			bcrypt := hashBcrypt(pass)
			fmt.Printf("bcrypt : %s\n", bcrypt)

		}
	}
}

func hashSha(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	passwordSum := []byte(s.Sum(nil))
	return base64.StdEncoding.EncodeToString(passwordSum)
}

func hashBcrypt(password string) (hash string) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(passwordBytes)
}
