package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"

	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
)

var tokenApiSecret = []byte{0x47, 0x30, 0x30, 0x78, 0x77, 0x30, 0x72, 0x6b, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x63, 0x4c, 0x40, 0x75, 0x64, 0x53, 0x33, 0x43, 0x52, 0x33, 0x54, 0x4b, 0x33, 0x59, 0x74, 0x30, 0x4b, 0x65, 0x4f}
var tokenSigningKy = []byte{0x40, 0x79, 0x61, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x40, 0x75, 0x64, 0x53, 0x33, 0x43, 0x52, 0x33, 0x54, 0x40, 0x50, 0x69, 0x4b, 0x33, 0x59, 0x74, 0x30, 0x4b, 0x65, 0x4f, 0x33, 0x32}

type TokenUtilsStruct struct {
}

var doOnceForTokenUtils sync.Once

var tokenUtilsSingleton *TokenUtilsStruct = nil

// TokenUtils is a factory method that acts as a static member
func TokenUtils() *TokenUtilsStruct {
	doOnceForTokenUtils.Do(func() {
		tokenUtilsSingleton = &TokenUtilsStruct{}
	})
	return tokenUtilsSingleton
}

// region ID Generators methods ----------------------------------------------------------------------------------------

// ID return a long string (10 characters) based on epoch micro-seconds in base 36
func (t *TokenUtilsStruct) ID() string {
	return strconv.FormatUint(uint64(time.Now().UnixMicro()), 36)
}

// ShortID return a short string (6 characters) based on epoch seconds in base 36
func (t *TokenUtilsStruct) ShortID(delta ...int) string {
	value := uint64(time.Now().Unix())
	for _, d := range delta {
		value += uint64(d)
	}
	return strconv.FormatUint(value, 36)
}

// GUID return a long string (36 characters) of a Global Unique Identifier (5 segments with dash separators)
func (t *TokenUtilsStruct) GUID() string {
	return uuid.New().String()
}

// NanoID return a long string on unique identifier based the nano-id generator
func (t *TokenUtilsStruct) NanoID() string {
	if nanoID, err := nanoid.Standard(21); err != nil {
		return t.GUID()
	} else {
		return nanoID()
	}
}

// endregion

// region Access Token parsing helpers ---------------------------------------------------------------------------------

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenData
}

// CreateToken build JWT token from Token Data structure
func (t *TokenUtilsStruct) CreateToken(td *TokenData) (string, error) {
	claims := TokenClaims{}
	claims.SubjectId = td.SubjectId
	claims.SubjectType = td.SubjectType
	claims.Status = td.Status
	claims.ExpiresIn = td.ExpiresIn
	claims.Subject = td.SubjectId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSigningKy)
}

// ParseToken rebuild Token Data structure from JWT token
func (t *TokenUtilsStruct) ParseToken(tokenString string) (*TokenData, error) {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenSigningKy, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and extract the claims
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return &TokenData{
			SubjectId:   claims.SubjectId,
			SubjectType: claims.SubjectType,
			Status:      claims.Status,
			ExpiresIn:   claims.ExpiresIn,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// endregion

// region API Key parsing helpers --------------------------------------------------------------------------------------

// CreateApiKey generate API Key from application name
func (t *TokenUtilsStruct) CreateApiKey(appName string) (string, error) {
	return encrypt(appName)
}

// ParseApiKey extract application name from API key
func (t *TokenUtilsStruct) ParseApiKey(apiKey string) (string, error) {
	return decrypt(apiKey)
}

// endregion

// region PRIVATE SECTION ----------------------------------------------------------------------------------------------

// encrypt string using AES and return base64
func encrypt(value string) (string, error) {

	block, err := aes.NewCipher(tokenApiSecret)
	if err != nil {
		return "", err
	}

	// Generate a new random IV
	cipherText := make([]byte, aes.BlockSize+len(value))
	iv := cipherText[:aes.BlockSize]
	if _, er := io.ReadFull(rand.Reader, iv); er != nil {
		return "", er
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(value))

	return hex.EncodeToString(cipherText), nil
}

// decrypt base64 string using AES
func decrypt(value string) (string, error) {
	cipherTextBytes, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(tokenApiSecret)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", fmt.Errorf("cipher text too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}

// endregion
