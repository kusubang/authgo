package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-app/pkg/usr"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// const SECRET = "secret"
const SUB = "token-subject"

type CustomClaim struct {
	// Sub   string `json:"sub"`
	Id    string `json:"id"`
	Email string `json:"email"`
}

type JwtCustomClaims struct {
	CustomClaim
	jwt.RegisteredClaims
}

type Auth struct {
	UserService usr.UserService
	secret      string
}

func New(uSvc usr.UserService) *Auth {
	return &Auth{uSvc, "secret"}
}

func (a *Auth) Login(id, email, pw string) (string, string, error) {

	_, err := a.UserService.Find(id, pw)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	c := CustomClaim{
		id, email,
	}
	return a.generateTokens(c)
}

func (a *Auth) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := a.IsValid(refreshToken)
	if err != nil {
		return "", "", err
	}

	jsonData, _ := json.Marshal(claims)
	var st CustomClaim
	fmt.Println(jsonData)

	json.Unmarshal(jsonData, &st)
	return a.generateTokens(st)
}

func (a *Auth) IsValid(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// 클라이언트로 받은 토큰이 HMAC 알고리즘이 맞는지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 발급할때의 토큰제목과 일치하는지 체크
	// if claims["sub"].(string) != SUB {
	// 	return nil, errors.New("invalid subject")
	// }

	return claims, nil
	// return token.Claims, nil
}

// generate tokens
func (a *Auth) generateTokens(claims CustomClaim) (string, string, error) {
	// access 토큰 생성: 유효기간 20분
	accessToken, err := a.createToken(
		claims,
		time.Now().Add(time.Minute*20),
	)
	if err != nil {
		return "", "", err
	}
	// refresh 토큰 생성: 유효기간 24시간
	refreshToken, err := a.createToken(
		claims,
		time.Now().Add(time.Hour*24),
	)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (a *Auth) createToken(data CustomClaim, expire time.Time) (string, error) {
	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := token.Claims.(jwt.MapClaims)
	// for key, val := range data {
	// 	claims[key] = val
	// }
	// claims["exp"] = expire.Unix()

	claims := &JwtCustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// // Generate encoded token and send it as response.
	// t, err := token.SignedString([]byte("secret"))
	// if err != nil {
	// 	return err
	// }

	encToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return encToken, nil
}
