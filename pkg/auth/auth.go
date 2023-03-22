package auth

import (
	"errors"
	"fmt"
	"go-app/pkg/usr"
	"time"

	"github.com/golang-jwt/jwt"
)

const SECRET = "secret"
const SUB = "token-subject"

type Auth struct {
	UserService usr.UserService
}

func New(uSvc usr.UserService) *Auth {
	return &Auth{uSvc}
}

type ClaimData map[string]interface{}

func (a *Auth) Login(id, email, pw string) (string, string, error) {

	_, err := a.UserService.Find(id, pw)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	claims := ClaimData{
		"sub":   SUB,
		"id":    id,
		"email": email,
	}
	return generateTokens(claims)
}

func (a *Auth) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := a.IsValid(refreshToken)
	cl := ClaimData(claims)
	if err != nil {
		return "", "", nil
	}
	return generateTokens(cl)
}

func (a *Auth) IsValid(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// 클라이언트로 받은 토큰이 HMAC 알고리즘이 맞는지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 발급할때의 토큰제목과 일치하는지 체크
	if claims["sub"].(string) != SUB {
		return nil, errors.New("invalid subject")
	}

	return claims, nil
}

// generate tokens
func generateTokens(claims ClaimData) (string, string, error) {
	// access 토큰 생성: 유효기간 20분
	accessToken, err := createToken(
		claims,
		time.Now().Add(time.Minute*20),
	)
	if err != nil {
		return "", "", err
	}
	// refresh 토큰 생성: 유효기간 24시간
	refreshToken, err := createToken(
		claims,
		time.Now().Add(time.Hour*24),
	)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func createToken(data ClaimData, expire time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range data {
		claims[key] = val
	}
	claims["exp"] = expire.Unix()

	encToken, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return "", err
	}
	return encToken, nil
}
