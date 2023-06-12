package sessionModel

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/handymesh/hyshAuthService/db/redis"
)

const (
	ACCESS_TOKEN_DURATION  = time.Hour * 1
	REFRESH_TOKEN_DURATION = time.Hour * 24 * 30 // 30 days
	RECOVERY_LINK_DURATION = time.Hour * 1
)

var (
	log        = logrus.New()
	signingKey = []byte("your-secret-key")
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func NewAccessToken(timeDuration int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(ACCESS_TOKEN_DURATION).Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(timeDuration time.Duration) (string, error) {
	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	err = redis.Redis.Set(refreshToken.String(), "true", REFRESH_TOKEN_DURATION).Err()
	if err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func NewRecoveryLink(value string) (string, error) {
	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	err = redis.Redis.Set(refreshToken.String(), value, RECOVERY_LINK_DURATION).Err()
	if err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func Delete(token string) error {
	err := redis.Redis.Del(token).Err()
	if err != nil {
		return err
	}
	return nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	return token, err
}

func CheckRefreshToken(token string) (bool, error) {
	value := redis.Redis.Get(token)
	if value.Err() != nil {
		return false, value.Err()
	}

	status, err := value.Result()
	if err != nil && status != "true" {
		return false, err
	}

	return true, nil
}

func GetValueByKey(token string) (string, error) {
	value := redis.Redis.Get(token)
	if value.Err() != nil {
		return "", value.Err()
	}

	status, err := value.Result()
	if err != nil && status != "true" {
		return "", err
	}

	return status, nil
}
