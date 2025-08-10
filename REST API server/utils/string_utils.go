package utils

import (
	"net/mail"
	"sync"
)

type StringUtilsStruct struct {
}

var doOnceForStringUtils sync.Once

var stringUtilsSingleton *StringUtilsStruct = nil

// StringUtils is a factory method that acts as a static member
func StringUtils() *StringUtilsStruct {
	doOnceForStringUtils.Do(func() {
		stringUtilsSingleton = &StringUtilsStruct{}
	})
	return stringUtilsSingleton
}

// IsValidEmail test for a valid email format
func (t *StringUtilsStruct) IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
