package main

import (
	"encoding/json"
	"time"
)

type CustomTime time.Time

var _ json.Unmarshaler = &CustomTime{}

func (ct *CustomTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}
