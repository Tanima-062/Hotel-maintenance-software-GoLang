package utils

import (
	"fmt"
	"strings"
	"unicode"
)

func UpperAndLowerStrList(str string) map[string]string {
	list := make(map[string]string)
	// 引数(str)が空ならそのまま返す
	if str == "" {
		return list
	}
	// 小文字で初期化すること
	str = strings.ToLower(str)
	lim := (1 << len(str)) - 1

	for i := 0; i <= lim; i++ {
		tstr := str
		for pos, c := range formatBinary(i, lim) {
			if c == rune('1') {
				tstr = toUpperSpecifiedPos(tstr, pos)
			}
		}
		if _, ok := list[tstr]; !ok {
			list[tstr] = tstr
		}
	}
	return list
}

func formatBinary(i, lim int) string {
	sVal := fmt.Sprintf("%b\n", i)
	sLim := fmt.Sprintf("%b\n", lim)

	if len(sVal) != len(sLim) {
		sVal = fmt.Sprintf("%s%s", strings.Repeat("0", len(sLim)-len(sVal)), sVal)
	}

	return sVal
}

func toUpperSpecifiedPos(str string, pos int) string {
	var rstr []rune
	for p, r := range str {
		if p == pos {
			rstr = append(rstr, unicode.ToUpper(r))
			continue
		}

		rstr = append(rstr, r)
	}

	return string(rstr)
}
