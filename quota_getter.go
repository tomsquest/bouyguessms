package bouyguessms

import (
	"errors"
	"regexp"
	"strconv"
)

type QuotaGetter interface {
	Get() (Quota, error)
}

type quotaGetter struct {
	client httpClient
}

func (getter *quotaGetter) Get() (Quota, error) {
	body, err := getter.client.Get("https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/sendSMS.phtml")
	if err != nil {
		return ExceededQuota, err
	}

	regex := regexp.MustCompile("Il vous reste <strong>(\\d*) SMS gratuit")
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ExceededQuota, errors.New("unable to read quota from body")
	}

	quota, err := strconv.Atoi(matches[1])
	if err != nil {
		return ExceededQuota, err
	}

	return Quota(quota), nil
}
