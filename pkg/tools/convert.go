package tools

import (
	"fmt"
	"time"
)

func ToRFC3339(dataString string, timezone string) (string, error) {
	// timezone "Europe/Moscow"
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("load location error %s", err)
	}

	// Из строки получаем объект Time в нужном формате и локейшене
	newOrderTime, err := time.ParseInLocation(time.RFC3339, dataString, loc)
	if err != nil {
		return "", fmt.Errorf("error conver time %s", err)
	}

	// Форматируем в: "2020-12-10T15:15:45+03:00"
	dataString = newOrderTime.In(loc).Format("2006-01-2T15:04:05Z07:00")
	return dataString, nil
}
