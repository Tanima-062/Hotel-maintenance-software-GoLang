package utils

import (
	"testing"
)

func Test_Crypto(t *testing.T) {
	t.Run("暗号化できることのテスト", func(t *testing.T) {
		test := map[string]string{
			"3Y3~VD4NnhN": "26Ae+mMwivgThASRFEpzcQ==",
			"akano":       "G9dRkNo5Lm8=",
			"AKANO":       "EBHfMF3EXS4=",
			"Akano":       "A2ojr5a54Bw=",
			"akanO":       "EdDV2UcSqY0=",
			"aKano":       "nj7IUapVvoI=",
			"AkanO":       "dDrhsdPmLYs=",
			"aKaNo":       "0VScTNT6/IA=",
			"qz06zisl":    "iO2yvu/VkcM=",
			"tomohiro":    "ngYPjCFtXj0=",
			"TOMOHIRO":    "c76+s91jKhA=",
			"Tomohiro":    "nJnxMsmPL/c=",
			"sato":        "qyD5MDvoohM=",
		}

		for plain, enc := range test {
			value, err := Encrypt(plain)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if value != enc {
				t.Fatalf("暗号化できませんでした。想定: %s, 実際: %s", enc, value)
			}
		}
	})

	t.Run("復号できることのテスト", func(t *testing.T) {
		test := map[string]string{
			"G9dRkNo5Lm8=":             "akano",
			"EBHfMF3EXS4=":             "AKANO",
			"A2ojr5a54Bw=":             "Akano",
			"EdDV2UcSqY0=":             "akanO",
			"nj7IUapVvoI=":             "aKano",
			"dDrhsdPmLYs=":             "AkanO",
			"0VScTNT6/IA=":             "aKaNo",
			"ngYPjCFtXj0=":             "tomohiro",
			"c76+s91jKhA=":             "TOMOHIRO",
			"nJnxMsmPL/c=":             "Tomohiro",
			"qyD5MDvoohM=":             "sato",
			"26Ae+mMwivgThASRFEpzcQ==": "3Y3~VD4NnhN",
			"iO2yvu/VkcM=":             "qz06zisl",
		}

		for enc, plain := range test {
			value, err := Decrypt(enc)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if value != plain {
				t.Fatalf("復号できませんでした。想定: %s, 実際: %s", plain, value)
			}
		}
	})
}
