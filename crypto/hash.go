package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func MD5(b []byte) string {
	en := md5.New()
	en.Write(b)
	return fmt.Sprintf("%x", en.Sum(nil))
}

func SHA1(b []byte) string {
	en := sha1.New()
	en.Write(b)
	return fmt.Sprintf("%x", en.Sum(nil))
}

func SHA256(b []byte) string {
	en := sha256.New()
	en.Write(b)
	return fmt.Sprintf("%x", en.Sum(nil))
}

func SHA512(b []byte) string {
	en := sha512.New()
	en.Write(b)
	return fmt.Sprintf("%x", en.Sum(nil))
}
