package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leonardchinonso/lokate-go/config"
	"github.com/leonardchinonso/lokate-go/models/dao"
)

type tokenCustomClaims struct {
	User *dao.UserDAO `json:"user"`
	jwt.StandardClaims
}

func generateToken(user *dao.UserDAO, jwtSecretKey string, expiresIn int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExpiresIn := unixTime + expiresIn

	// create a claims object
	claims := tokenCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiresIn,
			IssuedAt:  unixTime,
		},
	}

	// create a jwt token object and set the expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token string
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		log.Printf("Error generating single token for userId: %v. Error: %v\n", user.Id, err.Error())
		return "", err
	}

	return tokenString, nil
}

func generateAccessToken(user *dao.UserDAO) (string, error) {
	// get the access token secret key for signing the token
	atExpiresIn, err := strconv.Atoi(config.Map[config.ATExpiresIn])
	if err != nil {
		return "", err
	}

	// get the access token secret key
	atSecretKey := config.Map[config.ATSecretKey]

	return generateToken(user, atSecretKey, int64(atExpiresIn))
}

func generateRefreshToken(user *dao.UserDAO) (string, error) {
	// get the refresh token secret key for signing the token
	rtExpiresIn, err := strconv.Atoi(config.Map[config.RTExpiresIn])
	if err != nil {
		return "", err
	}

	// get the refresh token secret key
	rtSecretKey := config.Map[config.RTSecretKey]

	return generateToken(user, rtSecretKey, int64(rtExpiresIn))
}

func GenerateTokenPair(user *dao.UserDAO) (string, string, error) {
	at, err := generateAccessToken(user)
	if err != nil {
		log.Printf("Error generating access token for uid: %v. Error: %v\n", user.Id, err.Error())
		return "", "", err
	}

	rt, err := generateRefreshToken(user)
	if err != nil {
		log.Printf("Error generating refresh token for uid: %v. Error: %v\n", user.Id, err.Error())
		return "", "", err
	}

	token := &dao.TokenDAO{
		UserId:       user.Id,
		AccessToken:  at,
		RefreshToken: rt,
	}

	if err = token.Upsert(); err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func verifyAccessToken(tokenString, atSecretKey string) (*tokenCustomClaims, error) {
	claims := &tokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(atSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*tokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

func GetUserFromAccessToken(tokenString string) (*dao.UserDAO, error) {
	claims, err := verifyAccessToken(tokenString, config.Map[config.ATSecretKey])

	if err != nil {
		log.Printf("Unable to validate or parse access token. Error: %v\n", err)
		return nil, fmt.Errorf("cannot authenticate user: %v", err)
	}

	return claims.User, nil
}
