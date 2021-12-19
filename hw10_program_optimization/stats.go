package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var ErrEmptyDomain = errors.New("empty domain")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}
	res := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")

		if strings.HasSuffix(email, "."+domain) {
			res[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}
	return res, nil
}
