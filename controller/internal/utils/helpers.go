package utils

import (
	"math/rand"
)

func FlipCoin() bool {
	return rand.Intn(2) == 0 // Возвращает true или false с равной вероятностью
}
