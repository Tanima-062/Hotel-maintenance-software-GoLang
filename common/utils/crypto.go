package utils

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"errors"
	"regexp"
)

// ゼロパディングの際に末尾にくっつく固定値（削除する）
var re = regexp.MustCompile("0QD0jfRLfLw==$")

func Encrypt(src string) (string, error) {
	// 文字列が空ならそのまま返す
	if src == "" {
		return "", nil
	}
	key := []byte{0x68, 0x33, 0x4B, 0x46, 0x69, 0x41, 0x4A, 0x46}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	bsrc := ZeroPadding([]byte(src), bs)

	if len(bsrc)%bs != 0 {
		return "", errors.New("Blocksize Error")
	}

	out := make([]byte, len(bsrc))
	dst := out
	for len(bsrc) > 0 {
		block.Encrypt(dst, bsrc[:bs])
		bsrc = bsrc[bs:]
		dst = dst[bs:]
	}

	b64 := base64.StdEncoding.EncodeToString(out)

	// ゼロパティングで追加された文字列を削って返却する
	return re.ReplaceAllString(b64, "="), nil
}

func Decrypt(encrypted string) (string, error) {
	key := []byte{0x68, 0x33, 0x4B, 0x46, 0x69, 0x41, 0x4A, 0x46}
	//iv := []byte{0x83, 0xCA, 0x46, 0x49, 0x89, 0x7F, 0x76, 0x2A}

	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	bsrc, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	if len(bsrc)%bs != 0 {
		return "", errors.New("block size error")
	}

	out := make([]byte, len(bsrc))
	dst := out
	for len(bsrc) > 0 {
		block.Decrypt(dst, bsrc[:bs])
		bsrc = bsrc[bs:]
		dst = dst[bs:]
	}

	out = ZeroUnPadding(out)

	return string(out), nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
