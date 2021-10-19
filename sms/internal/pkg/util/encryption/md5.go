package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func GetMD516Encode(data string) string {
	return EncodeMD5(data)[8:24]
}
