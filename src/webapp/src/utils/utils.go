package utils

import (
	"fmt"
	"strings"
	"time"
)

type JsonSpecialDateTime struct {
	time.Time
}

func (sd *JsonSpecialDateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if strInput == "0000-00-00 00:00:00" {
		sd.Time = time.Time{}
		return nil
	}
	newTime, err := time.Parse("2006-01-02 15:04:05", strInput)
	if err != nil {
		return err
	}
	sd.Time = newTime
	return nil
}

func (sd *JsonSpecialDateTime) MarshalJSON() ([]byte, error) {
	if sd.Time.IsZero() {
		return []byte("\"0001-01-01 00:00:00\""), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", sd.Time.Format("2006-01-02 15:04:05"))), nil
}
