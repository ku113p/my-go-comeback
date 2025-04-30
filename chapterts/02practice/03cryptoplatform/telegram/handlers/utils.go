package handlers

import (
	"crypto/platform/models"
	"fmt"
	"strconv"
	"strings"
)

func newNotificationFromString(s string) (*models.Notification, error) {
	words := strings.SplitN(s, " ", 3)
	if len(words) != 3 {
		return nil, fmt.Errorf("invalid format")
	}

	symbol, signString, amountString := words[0], words[1], words[2]
	symbol = strings.ToUpper(symbol)

	sign, err := models.ParseSign(signString)
	if err != nil {
		return nil, fmt.Errorf("invalid sign")
	}

	amountString = strings.ReplaceAll(amountString, ",", ".")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount")
	}

	msg := fmt.Sprintf("price %v %v %v", symbol, sign.String(), amount)
	n := models.NewNotification(symbol, *sign, amount, nil, &msg)

	return n, nil
}
