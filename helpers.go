package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func countRetailerName(points int, retailer string) int {
	count := 0
	// loop through every character and check if it is a letter or number
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}

	points += count

	return points
}

func roundDollarAmount(points int, total string) (int, error) {
	// make sure that total is valid dollar amount
	if !strings.Contains(total, ".") {
		return 0, errors.New("total is not a valid decimal value")
	}

	// split the string on the period and check if the second value in the array is 00
	splitTotal := strings.Split(total, ".")
	if splitTotal[1] == "00" {
		points += 50
	}

	return points, nil
}

func multipleOfPointTwoFive(points int, total string) (int, error) {
	// convert string to float
	floatTotal, err := convertStringToFloat(total)
	if err != nil {
		return 0, err
	}

	// round value to nearest after two decimal places
	floatTotal = math.Round(floatTotal*100) / 100

	res := math.Mod(floatTotal, 0.25)

	// round value up after two decimal places
	res = math.Ceil(res*100) / 100
	if res == 0 {
		points += 25
	}

	return points, nil
}

func everyTwoItems(points int, items []struct {
	ShortDescription string "json:\"shortDescription\""
	Price            string "json:\"price\""
}) int {
	// divide length by two
	itemsLength := len(items) / 2
	points += itemsLength * 5

	return points
}

func trimmedTotal(points int, items []struct {
	ShortDescription string "json:\"shortDescription\""
	Price            string "json:\"price\""
}) (int, error) {
	for _, val := range items {
		trimmed := strings.TrimSpace(val.ShortDescription)
		if len(trimmed)%3 == 0 {

			// convert string to float
			floatTotal, err := convertStringToFloat(val.Price)
			if err != nil {
				return 0, err
			}

			// multiply by 0.2
			floatTotal = floatTotal * .2

			// round value up to nearest int and add to points
			points += int(math.Ceil(floatTotal))
		}
	}

	return points, nil
}

func oddPurchaseDate(points int, date string) (int, error) {
	// make sure purchase date is valid
	if !strings.Contains(date, "-") {
		return 0, errors.New("purchaseDate is not a valid date value")
	}

	splitDate := strings.Split(date, "-")
	day := splitDate[2]
	// get the last value
	lastVal := day[len(day)-1:]
	// convert to int
	intLastVal, err := strconv.Atoi(lastVal)
	if err != nil {
		return 0, err
	}
	if int(intLastVal)%2 != 0 {
		points += 6
	}

	return points, nil
}

func purchaseTime(points int, purchaseTime string) (int, error) {
	// make sure purchase time is valid
	if !strings.Contains(purchaseTime, ":") {
		return 0, errors.New("purchaseTime is not a valid time value")
	}

	date, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return 0, err
	}

	// time of purchase must be after 14:00 and before 16:00
	if date.Hour() >= 14 && date.Hour() < 16 {
		points += 10
	}

	return points, nil
}

func convertStringToFloat(s string) (float64, error) {
	floatTotal, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0.0, err
	}

	return floatTotal, nil
}