package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"net/mail"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	var user User
	result := make(DomainStat)
	d := "." + domain
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("error unmarshalling user: %w", err)
		}

		user.Email = strings.ToLower(user.Email)
		_, err := mail.ParseAddress(user.Email)
		if err == nil && strings.HasSuffix(user.Email, d) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
