package core

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const PASS_LENGTH = 256

type Password [PASS_LENGTH]byte

func init() {
	rand.Seed(time.Now().Unix())
}

func (pwd *Password) String() string {
	return base64.StdEncoding.EncodeToString(pwd[:])
}

func ParsePassword(pwd string) (*Password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(pwd))
	if err != nil || len(bs) != PASS_LENGTH {
		return nil, errors.New("不合法密码")

	}
	password := Password{}
	copy(password[:], bs)
	bs = nil
	return &password, nil
}

func RandPassword() string {
	intArr := rand.Perm(PASS_LENGTH)
	Password := &Password{}
	for i, v := range intArr {
		Password[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return Password.String()
}
