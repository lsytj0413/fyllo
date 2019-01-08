package main

import (
	"fmt"
	"net"
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

	names, err := net.LookupIP("map.baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(names)

	arr1 := []int{}
	arr2 := []int{}
	for i := 1; i <= 16; i++ {
		arr1 = append(arr1, i)
	}

	for len(arr1) > 2 {
		for i := 0; i < len(arr1); i++ {
			arr2 = append(arr2, arr1[i])
			i++
		}
		arr1, arr2 = arr2, nil
		// fmt.Println(len(arr1))
	}
	fmt.Println(arr1[1])

	done := false

	go func() {
		done = true
	}()

	for !done {
		fmt.Println("not done!") //not inlined
	}
	fmt.Println("done!")
}
