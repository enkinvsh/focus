package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Language  string `json:"language_code"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "tma ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing authorization"})
			return
		}

		initData := auth[4:]
		user, err := validateInitData(initData)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func validateInitData(initData string) (*TelegramUser, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, errors.New("invalid init data format")
	}

	hash := values.Get("hash")
	if hash == "" {
		return nil, errors.New("missing hash")
	}

	var keys []string
	for k := range values {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, k+"="+values.Get(k))
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	botToken := os.Getenv("BOT_TOKEN")
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))

	h := hmac.New(sha256.New, secretKey.Sum(nil))
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != hash {
		return nil, errors.New("invalid signature")
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, errors.New("missing user data")
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, errors.New("invalid user data")
	}

	return &user, nil
}

func GetUser(c *gin.Context) *TelegramUser {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*TelegramUser)
}
