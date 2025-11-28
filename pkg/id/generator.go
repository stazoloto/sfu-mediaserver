package id

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

func GenerateIDByLength(length int) (string, error) {
	if length < 1 {
		return "", errors.New("length must be greater than zero")
	}

	// Максимальное число для заданного диапазона
	maxNum := new(big.Int)
	maxNum.Exp(big.NewInt(10), big.NewInt(int64(length)), nil)

	//	Случайное число в диапазоне [0, maxNum)
	randomNum, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, randomNum), nil
}

func GenerateIDWithPrefix(prefix string, length int) (string, error) {
	numID, err := GenerateIDByLength(length)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", prefix, numID), nil
}
