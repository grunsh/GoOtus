package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsonIter "github.com/json-iterator/go"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	user := User{}
	res := make(DomainStat)

	for scanner.Scan() {
		if err := jsonIter.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if !strings.Contains(user.Email, "."+domain) || !strings.Contains(user.Email, "@") {
			continue
		}

		matchDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		res[matchDomain]++
	}

	return res, nil
}
