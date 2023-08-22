package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

func ParseTime(input string) time.Time {
	layout := "2006-01-02 15:04:05"
	loc, err := time.LoadLocation("Asia/Jakarta") // contoh lokasi
	if err != nil {
		fmt.Println(err)
	}
	parseResult, _ := time.ParseInLocation(layout, input, loc)
	return parseResult
}

func EncryptString(text string) (string, error) {
	plaintext := []byte(text)
	key := []byte(viper.GetString("SECRET.SECRET_KEY_JWT"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(viper.GetString("SECRET.SECRET_KEY_JWT")))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	fmt.Println("cipper" + string(ciphertext))
	return string(ciphertext), nil
}

func ValidateUUID(encrypt string, uuid string) bool {
	decripted, err := DecryptString(encrypt)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return decripted == uuid
}
func TimeInLocal(local string) (time.Time, time.Time) {
	t := time.Now()

	location, err := time.LoadLocation(local)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("Location : ", location, " Time : ", t.In(location)) // America/New_York
	resp := t.In(location)

	return t, resp
}
