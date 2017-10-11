package bouyguessms

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type phoneNumber string

func (phone phoneNumber) isValid() bool {
	regex := regexp.MustCompile("^0[67]\\d{8}$")
	return regex.MatchString(string(phone))
}

type phoneNumbers []phoneNumber

func parsePhones(raw string) (phoneNumbers, error) {
	rawPhones := strings.Split(raw, ";")
	if len(rawPhones) > 5 {
		return nil, errors.New("too many phone numbers given (5 is the max)")
	}

	phones := []phoneNumber{}
	for _, rawPhone := range rawPhones {
		phoneNumber := phoneNumber(strings.TrimSpace(rawPhone))
		if !phoneNumber.isValid() {
			return nil, errors.Errorf("invalid phone number %s", phoneNumber)
		}
		phones = append(phones, phoneNumber)
	}

	return phones, nil
}

func (phones phoneNumbers) String() string {
	strs := []string{}
	for _, phone := range phones {
		strs = append(strs, string(phone))
	}

	return strings.Join(strs, ";")
}
