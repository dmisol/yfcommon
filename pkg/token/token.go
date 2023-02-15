package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dmisol/yfcommon/pkg/model"
	"github.com/golang-jwt/jwt/v4"
)

var secret []byte

func Sign(req *model.TokenReq) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	if _, err = time.LoadLocation(req.Timezone); err != nil {
		log.Println("LoadLocation", err)
		return
	}

	claims["authorized"] = true

	claims["nbf"] = req.Since
	claims["exp"] = req.Until

	claims["tz"] = req.Timezone
	claims["addr"] = req.Addr
	claims["devid"] = req.DeviceId

	//log.Println("signing", claims)
	tokenString, err = token.SignedString(secret)
	//log.Println("token is", tokenString)
	return
}

func Decode(raw string) (devid string, addr string, t0 time.Time, t1 time.Time, err error) {

	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: jwt parsing")
		}
		return secret, nil
	})

	if err != nil {
		if token == nil {
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			l, _ := time.LoadLocation(claims["tz"].(string))

			t0 = time.Unix(int64(claims["nbf"].(float64)), 0).In(l)
			t1 = time.Unix(int64(claims["exp"].(float64)), 0).In(l)

			err = fmt.Errorf("the Toker is disabled now.\nThe valid period is\nsince %v\nuntil %v", t0, t1)
			return
		}
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	l, _ := time.LoadLocation(claims["tz"].(string))
	t0 = time.Unix(int64(claims["nbf"].(float64)), 0).In(l)
	t1 = time.Unix(int64(claims["exp"].(float64)), 0).In(l)

	if ok && token.Valid {
		devid = claims["devid"].(string)
		addr = claims["addr"].(string)
		return
	}
	err = errors.New(fmt.Sprint("token claims failed or invalid", ok, token.Valid))
	return
}
