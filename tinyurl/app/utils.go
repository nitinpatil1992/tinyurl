package app

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateShortString() string {

	var alphaNumeric string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	firstIndex := rand.Intn(len(alphaNumeric) - 1)

	uid := uuid.New().String()
	uuidLastPart := strings.Split(uid, "-")[4]

	shortUrl := string(alphaNumeric[firstIndex]) + uuidLastPart
	return shortUrl
}
