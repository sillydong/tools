package main

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	mrand "math/rand"
	"time"

	jose "github.com/square/go-jose"
	"github.com/urfave/cli"
	"golang.org/x/crypto/pbkdf2"
)

func init() {
	commands = append(commands,
		cli.Command{
			Name:  "jwk",
			Usage: "jwk encode/decode keys",
			Subcommands: []cli.Command{
				cli.Command{
					Name:    "encode",
					Aliases: []string{"e"},
					Subcommands: []cli.Command{
						cli.Command{
							Name:    "public",
							Aliases: []string{"pub"},
							Usage:   "encode public key to jwk",
							Flags: []cli.Flag{
								cli.StringFlag{Name: "key", Usage: "public key file"},
							},
							Action: encodeToJson,
						},
						cli.Command{
							Name:    "private",
							Aliases: []string{"priv"},
							Usage:   "encode private key to jwk",
							Flags: []cli.Flag{
								cli.StringFlag{Name: "key", Usage: "private key file"},
							},
						},
					},
				},
			},
		},
	)
}

func encodeToJson(ctx *cli.Context) {
	file := ctx.String("key")
	if file == "" {
		cli.ShowCommandHelpAndExit(ctx, "public", 1)
	}

	keycontent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	keytype, key, err := LoadPublicKey(keycontent)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("handling " + keytype)
	rsakey := key.(rsa.PublicKey)

	pub := jose.JSONWebKey{Key: &rsakey, Algorithm: "RSA"}
	pubjson, err := pub.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", string(pubjson))
}

const (
	INT_SIZE  = 4
	INT2_SIZE = 2
)

const (
	ENCRYPT_DES_ECB_PKCS5           = 0 // DES with ECB and PKCS5 padding
	ENCRYPT_NONE                    = 1 // no encryption
	ENCRYPT_DES_CBC_PKCS5           = 2 // DES with CBC and PKCS5 padding
	ENCRYPT_PBKDF2_DESEDE_CBC_PKCS5 = 3 // DESede with CBC and PKCS5 padding and PBKDF2 to derive encryption key
	ENCRYPT_PBKDF2_AES_CBC_PKCS5    = 4 // AES with CBC and PKCS5 padding and PBKDF2 to derive encryption key
)

const (
	KEY_ENCODING_DSA_PRIVATE    = "DSA_PRIV_KEY"
	KEY_ENCODING_DSA_PUBLIC     = "DSA_PUB_KEY"
	KEY_ENCODING_DH_PRIVATE     = "DH_PRIV_KEY"
	KEY_ENCODING_DH_PUBLIC      = "DH_PUB_KEY"
	KEY_ENCODING_RSA_PRIVATE    = "RSA_PRIV_KEY"
	KEY_ENCODING_RSACRT_PRIVATE = "RSA_PRIVCRT_KEY"
	KEY_ENCODING_RSA_PUBLIC     = "RSA_PUB_KEY"
)

func LoadPrivateKey(buf, password []byte) (string, crypto.PrivateKey, error) {
	buf, err := decrypt(buf, password)
	if err != nil {
		return "", nil, err
	}
	keytype, key, err := bytesToPrivateKey(buf)
	if err != nil {
		return keytype, nil, err
	}
	return keytype, key, nil
}

func LoadPublicKey(buf []byte) (string, crypto.PublicKey, error) {
	return bytesToPublicKey(buf)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func nonce() []byte {
	mrand.Seed(time.Now().UnixNano())
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[mrand.Intn(len(letterRunes))]
	}
	return []byte(string(b))
}

func readInt(buf []byte, offset int) int {
	return int(buf[offset])<<24 | (0x00ff&int(buf[offset+1]))<<16 | (0x00ff&int(buf[offset+2]))<<8 | (0x00ff & int(buf[offset+3]))
}

func readInt2(buf []byte, offset int) int {
	return int(buf[offset])<<8 | int(buf[offset+1])&0x00ff
}

func readBytes(buf []byte, offset int) []byte {
	length := readInt(buf, offset)
	if len(buf) < length {
		return nil
	}
	return buf[offset+INT_SIZE : length+offset+INT_SIZE]
}

func bytesToBigInt(buf []byte) *big.Int {
	return new(big.Int).SetBytes(buf)
}

func bytesToInt(buf []byte) int {
	if len(buf) > 4 {
		panic(fmt.Errorf("%+v can not be convert to int", buf))
	}
	if len(buf) == 4 {
		return int(binary.BigEndian.Uint32(buf))
	}
	b := make([]byte, 4)
	copy(b[4-len(buf):], buf)
	return int(binary.BigEndian.Uint32(b))
}

func decrypt(buf, password []byte) ([]byte, error) {
	enctype := readInt(buf, 0)
	switch enctype {
	case ENCRYPT_DES_ECB_PKCS5:
		return nil, fmt.Errorf("%d not implemented", enctype)
	case ENCRYPT_DES_CBC_PKCS5:
		return nil, fmt.Errorf("%d not implemented", enctype)
	case ENCRYPT_PBKDF2_DESEDE_CBC_PKCS5:
		return nil, fmt.Errorf("%d not implemented", enctype)
	case ENCRYPT_PBKDF2_AES_CBC_PKCS5:
		offset := 4
		salt := readBytes(buf, offset)
		offset += INT_SIZE + len(salt)
		iteration := readInt(buf, offset)
		offset += INT_SIZE
		keyLength := readInt(buf, offset)
		offset += INT_SIZE
		iv := readBytes(buf, offset)
		offset += INT_SIZE + len(iv)
		ciphertext := readBytes(buf, offset)

		secretKey := pbkdf2.Key(password, salt, iteration, keyLength, sha1.New)
		if len(secretKey) > 16 {
			secretKey = secretKey[0:16]
		}

		block, err := aes.NewCipher(secretKey[0:16])
		if err != nil {
			return nil, err
		}
		if len(ciphertext) < aes.BlockSize {
			return nil, errors.New("Ciphertext block size too short")
		}
		stream := cipher.NewCBCDecrypter(block, iv)
		stream.CryptBlocks(ciphertext, ciphertext)
		return ciphertext, nil
	case ENCRYPT_NONE:
		return buf[4:], nil
	default:
		return nil, fmt.Errorf("%d not supported", enctype)
	}
}

func bytesToPrivateKey(priv []byte) (string, crypto.PrivateKey, error) {
	keytype := string(readBytes(priv, 0))
	offset := INT_SIZE + len(keytype)
	switch keytype {
	case KEY_ENCODING_DSA_PRIVATE:
		x := readBytes(priv, offset)
		offset += INT_SIZE + len(x)
		p := readBytes(priv, offset)
		offset += INT_SIZE + len(p)
		q := readBytes(priv, offset)
		offset += INT_SIZE + len(q)
		g := readBytes(priv, offset)
		offset += INT_SIZE + len(g)

		key := dsa.PrivateKey{
			PublicKey: dsa.PublicKey{
				Parameters: dsa.Parameters{
					P: bytesToBigInt(p),
					Q: bytesToBigInt(q),
					G: bytesToBigInt(g),
				},
			},
			X: bytesToBigInt(x),
		}
		return keytype, key, nil
	case KEY_ENCODING_RSA_PRIVATE:
		m := readBytes(priv, offset)
		offset += INT_SIZE + len(m)
		exp := readBytes(priv, offset)
		offset += INT_SIZE + len(exp)

		// TODO need to be tested, but handle-core just generate RSAPrivateCrtKey
		key := rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: bytesToBigInt(m),
				E: 65537, // Default length
			},
			D: bytesToBigInt(exp),
		}
		err := key.Validate()
		if err != nil {
			return keytype, nil, err
		}
		return keytype, &key, nil
	case KEY_ENCODING_RSACRT_PRIVATE:
		n := readBytes(priv, offset)
		offset += INT_SIZE + len(n)
		pubEx := readBytes(priv, offset)
		offset += INT_SIZE + len(pubEx)
		ex := readBytes(priv, offset)
		offset += INT_SIZE + len(ex)
		p := readBytes(priv, offset)
		offset += INT_SIZE + len(p)
		q := readBytes(priv, offset)
		offset += INT_SIZE + len(q)
		exP := readBytes(priv, offset)
		offset += INT_SIZE + len(exP)
		exQ := readBytes(priv, offset)
		offset += INT_SIZE + len(exQ)
		coeff := readBytes(priv, offset)
		offset += INT_SIZE + len(coeff)

		key := rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: bytesToBigInt(n),
				E: bytesToInt(pubEx),
			},
			D: bytesToBigInt(ex),

			Primes: []*big.Int{
				bytesToBigInt(p),
				bytesToBigInt(q),
			},
		}

		err := key.Validate()
		if err != nil {
			return keytype, nil, err
		}
		return keytype, &key, nil
	default:
		return keytype, nil, errors.New(keytype + " not supported")
	}
}

func bytesToPublicKey(pub []byte) (string, crypto.PublicKey, error) {
	keytype := string(readBytes(pub, 0))
	offset := INT_SIZE + len(keytype)
	// flags := readInt2(pub, offset) // currently not used... reserved
	offset += INT2_SIZE
	switch keytype {
	case KEY_ENCODING_DSA_PUBLIC:
		q := readBytes(pub, offset)
		offset += INT_SIZE + len(q)
		p := readBytes(pub, offset)
		offset += INT_SIZE + len(p)
		g := readBytes(pub, offset)
		offset += INT_SIZE + len(g)
		y := readBytes(pub, offset)
		offset += INT_SIZE + len(y)

		key := dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: bytesToBigInt(p),
				Q: bytesToBigInt(q),
				G: bytesToBigInt(g),
			},
			Y: bytesToBigInt(y),
		}

		return keytype, key, nil
	case KEY_ENCODING_RSA_PUBLIC:
		ex := readBytes(pub, offset)
		offset += INT_SIZE + len(ex)
		m := readBytes(pub, offset)
		offset += INT_SIZE + len(m)

		key := rsa.PublicKey{
			N: bytesToBigInt(m),
			E: bytesToInt(ex),
		}

		if key.E < 2 || key.E > 1<<31-1 {
			return keytype, nil, errors.New("public key error")
		}

		return keytype, key, nil
	case KEY_ENCODING_DH_PUBLIC:
		return keytype, nil, errors.New(keytype + " not implemented")
	default:
		return keytype, nil, errors.New(keytype + " not supported")
	}
}
