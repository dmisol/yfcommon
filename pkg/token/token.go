package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dmisol/yfcommon/pkg/model"
	//"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v5"
)

var (
	Secret      []byte
	UserTimeout = time.Hour
)

func SignSingle(req *model.TokenReq) (tokenString string, err error) {

	if _, err = time.LoadLocation(req.Timezone); err != nil {
		fmt.Println("LoadLocation", err)
		return
	}

	claims := &model.GuestToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(req.Until, 0)),
			NotBefore: jwt.NewNumericDate(time.Unix(req.Since, 0)),
		},
		Tz:    req.Timezone,
		Addr:  req.Addr,
		DevId: req.DeviceId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	tokenString, err = token.SignedString(Secret)
	return
}

func SignMultiple(addr string, devs map[string]string, since int64, until int64, tz string) (tokenString string, err error) {
	if _, err = time.LoadLocation(tz); err != nil {
		fmt.Println("LoadLocation", err)
		return
	}
	if len(devs) == 0 {
		err = fmt.Errorf("No devices to include")
		return
	}
	if len(devs) == 1 {
		for _, v := range devs {
			tr := &model.TokenReq{
				Since:    since,
				Until:    until,
				Timezone: tz,
				Addr:     addr,
				DeviceId: v,
			}
			return SignSingle(tr)
		}
	}

	claims := &model.GuestToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(until, 0)),
			NotBefore: jwt.NewNumericDate(time.Unix(since, 0)),
		},
		Tz:      tz,
		Addr:    addr,
		Devices: devs,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	tokenString, err = token.SignedString(Secret)
	return
}

// todo: wipe it all!
func DecodeKey(raw string) (devid string, devices map[string]string, addr string, t0 time.Time, t1 time.Time, tid int64, err error) {

	token, err := jwt.ParseWithClaims(raw, &model.GuestToken{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return Secret, nil
	})

	if err == jwt.ErrTokenSignatureInvalid || err == jwt.ErrSignatureInvalid {
		return
	}

	if err != nil {
		if token == nil {
			return
		}
		if strings.Contains(err.Error(), "signature is invalid") {
			return
		}
		claims, ok := token.Claims.(*model.GuestToken)
		if ok {
			l, _ := time.LoadLocation(claims.Tz)

			t0 = time.Unix(claims.NotBefore.Unix(), 0).In(l)
			t1 = time.Unix(claims.ExpiresAt.Unix(), 0).In(l)

			devid = claims.DevId
			devices = claims.Devices
			addr = claims.Addr

			err = fmt.Errorf("the Key is disabled now.\nThe valid period is\nsince %v\nuntil %v", t0, t1)
			return
		}
		return
	}

	claims, ok := token.Claims.(*model.GuestToken)

	l, _ := time.LoadLocation(claims.Tz)

	t0 = time.Unix(claims.NotBefore.Unix(), 0).In(l)
	t1 = time.Unix(claims.ExpiresAt.Unix(), 0).In(l)

	if ok && token.Valid {
		devid = claims.DevId
		devices = claims.Devices
		addr = claims.Addr
		tid = claims.NotifyTid
		return
	}

	err = errors.New(fmt.Sprint("token claims failed or invalid", ok, token.Valid))
	return
}

func SignUserDevice(user string, devid string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true

	now := time.Now()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(UserTimeout).Unix()

	claims["user"] = user
	if len(devid) > 0 {
		claims["devid"] = devid
	}

	tokenString, err = token.SignedString(Secret)
	return
}

func ValidateUserDevice(raw string) (user string, devid string, err error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: jwt parsing")
		}
		return Secret, nil
	})
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user = claims["user"].(string)
		d := claims["devid"]
		if d != nil {
			devid = d.(string)
		}
		return
	}
	err = fmt.Errorf("token invalid")
	return
}
