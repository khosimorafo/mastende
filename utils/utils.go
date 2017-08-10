package utils

import (
	"time"
	"fmt"
	"errors"
	"math/rand"
)

//Utilities
func DateFormatter(date string) (string, time.Time, error)  {

	layout := "2006-01-02"

	t, err := time.Parse(layout, date)

	ret_t := t.Format(layout)

	if err != nil {
		fmt.Println(err)
		return "", t, errors.New("Date submitted is invalid. ")
	}


	return ret_t, t, nil
}

func DateGetNow() (string, time.Time, error)  {

	layout := "2006-01-02"

	t := time.Now()

	ret_t := t.Format(layout)

	return ret_t, t, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var numberRunes = []rune("0123456789")

func RandNumberRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}
