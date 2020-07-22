package hash

import (
	"crypto"
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func EncodeSHA256(value string) string {
	m := crypto.SHA256.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
