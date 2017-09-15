package utils

import (
	"time"
	"fmt"
	"errors"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var numberRunes = []rune("0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandNumberRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}

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

func DateResolver(destination *string, input interface{}) error {
	
	result, _, err := DateFormatter(input.(string))

	if err != nil {

		return err
	}

	*destination = result
	return nil
}


/**
 * If the date is invalid, its is acceptable to present todays date as the answer.
 **/
func DateResolverDefaultOnToday(destination *string, input interface{}) error {

	result, _, err := DateFormatter(input.(string))

	if err != nil {

		result,_, _ := DateGetNow()
		*destination = result

	} else {

		*destination = result
	}

	return nil
}

func StringResolver(destination *string, input interface{}) error {
	
	var result string = ""

	if _, isOK := input.(string); isOK {
		 
		result = input.(string)
	} else {

		*destination = result
		return errors.New("Invalid string input. ")
	}
	
	*destination = result
	return nil
}

func DateGetNow() (string, time.Time, error)  {

	layout := "2006-01-02"

	t := time.Now()

	ret_t := t.Format(layout)

	return ret_t, t, nil
}


