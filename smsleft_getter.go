package bouyguessms

import (
	"regexp"
	"strconv"
	"github.com/pkg/errors"
)

type smsLeftGetter interface {
	Get() (SmsLeft, error)
}

type httpSmsLeftGetter struct {
	client httpClient
}

func (getter *httpSmsLeftGetter) Get() (SmsLeft, error) {
	body, err := getter.client.Get("https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/sendSMS.phtml")
	if err != nil {
		return NoSmsLeft, err
	}

	regex := regexp.MustCompile("Il vous reste <strong>(\\d*) SMS gratuit")
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return NoSmsLeft, errors.New("unable to read SMS left from body")
	}

	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return NoSmsLeft, errors.Wrapf(err, "unable to convert %v to int", matches[1])
	}

	return SmsLeft(count), nil
}
