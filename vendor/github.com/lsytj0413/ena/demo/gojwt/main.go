package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var _ = jwt.New

var secret = "xxyy=="

func main() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "demo",                                                // 签发者
		"sub": "user",                                                // 面向的用户
		"aud": "user",                                                // 接收的一方
		"exp": time.Now().Add(time.Minute * time.Duration(1)).Unix(), // 过期时间
		"nbf": "",                                                    // 在什么时间点之前不可用
		"iat": time.Now().Unix(),                                     // 签发时间
		"jti": "1",                                                   // 一次性标识
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Signed error: ", err)
		return
	}
	fmt.Println("Signed token: ", tokenString)

	// time.Sleep(time.Minute * 2)
	parseToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Parse error: ", err)
		return
	}
	claims := parseToken.Claims.(jwt.MapClaims)
	err = claims.Valid()
	if err != nil {
		fmt.Println("MapClaims valid error: ", err)
		return
	}

	for k, v := range claims {
		fmt.Println(k, ": ", v)
	}
}
