package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateClassCode(id uint) string {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(1<<16 - 1)
	return fmt.Sprintf("%04X%02X", random, id)
}
