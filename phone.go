package bouyguessms

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type PhoneNumber string

func (phone PhoneNumber) isValid() bool {
	regex := regexp.MustCompile("^0[67]\\d{8}$")
	return regex.MatchString(string(phone))
}

type PhoneNumbers []PhoneNumber

func ParsePhones(raw string) (PhoneNumbers, error) {
	rawPhones := strings.Split(raw, ";")
	if len(rawPhones) > 5 {
		return nil, errors.New("too many phone numbers given (5 is the max)")
	}

	phones := []PhoneNumber{}
	for _, rawPhone := range rawPhones {
		phoneNumber := PhoneNumber(strings.TrimSpace(rawPhone))
		if !phoneNumber.isValid() {
			return nil, errors.Errorf("invalid phone number %s", phoneNumber)
		}
		phones = append(phones, phoneNumber)
	}

	return phones, nil
}

func (phones PhoneNumbers) String() string {
	strs := []string{}
	for _, phone := range phones {
		strs = append(strs, string(phone))
	}

	return strings.Join(strs, ";")
}
