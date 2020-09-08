package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//var AesIV = []byte{0x19, 0x20, 0x84, 0x11, 0x02, 0x01, 0x11, 0x29, 0x78, 0x12, 0x07, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

// AesCBCEncrypt .
func AesCBCEncrypt(EncSrc, key, iv []byte) ([]byte, error) {
	if len(key) == 0 || len(iv) == 0 {
		err := fmt.Errorf("AesCBCEncrypt, key or iv null")
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) < blockSize {
		err = fmt.Errorf("AesCBCEncrypt, iv length error")
		return nil, err
	}
	EncSrc = PKCS5Padding(EncSrc, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	crypted := make([]byte, len(EncSrc))
	blockMode.CryptBlocks(crypted, EncSrc)
	return crypted, nil
}

// AesCBCDecrypt .
func AesCBCDecrypt(crypted, key, iv []byte) ([]byte, error) {
	if len(key) == 0 || len(iv) == 0 {
		err := fmt.Errorf("AesCBCDecrypt, key or iv null")
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) < blockSize {
		err = fmt.Errorf("AesCBCDecrypt, iv length error")
		return nil, err
	}
	if len(crypted)%blockSize != 0 {
		return nil, fmt.Errorf("AesCBCDecrypt, data length: %d, 不是 %d的整数倍", len(crypted), blockSize)
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	_t := int(origData[len(origData)-1])
	if _t > 16 {
		err = fmt.Errorf("AesCBCDecrypt, PKCS5Padding, Value: %d", _t)
		return nil, err
	}
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// AesCFBEncrypt .
func AesCFBEncrypt(EncSrc, key, iv []byte) ([]byte, error) {
	if len(key) == 0 || len(iv) == 0 {
		err := fmt.Errorf("AesCFBEncrypt, key or iv null")
		return nil, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//
	blockSize := c.BlockSize()
	if len(iv) < blockSize {
		err = fmt.Errorf("AesCFBEncrypt, iv length error")
		return nil, err
	}
	//cfb := cipher.NewCFBEncrypter(c, AesIV)
	cfb := cipher.NewCFBEncrypter(c, iv[:blockSize])
	crypted := make([]byte, len(EncSrc))
	cfb.XORKeyStream(crypted, EncSrc)
	//
	return crypted, nil
}

// AesCFBDecrypt .
func AesCFBDecrypt(DecSrc, key, iv []byte) ([]byte, error) {
	if len(key) == 0 || len(iv) == 0 {
		err := fmt.Errorf("AesCFBDecrypt, key or iv null")
		return nil, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//
	blockSize := c.BlockSize()
	if len(iv) < blockSize {
		err = fmt.Errorf("AesCFBDecrypt, iv length error")
		return nil, err
	}
	//cfbdec := cipher.NewCFBDecrypter(c, AesIV)
	cfbdec := cipher.NewCFBDecrypter(c, iv[:blockSize])
	origData := make([]byte, len(DecSrc))
	cfbdec.XORKeyStream(origData, DecSrc)
	//
	return origData, nil
}
