package bouyguessms

import "github.com/pkg/errors"

// GetQuota fetches the remaining number of SMS Left
func GetQuota(login, pass string) (Quota, error) {
	client, err := newHTTPClient()
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to create httpClient")
	}

	loginner := defaultLoginner{client}
	if err = loginner.Login(login, pass); err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to login")
	}

	quotaGetter := &httpQuotaGetter{client}
	return quotaGetter.Get()
}

// SendSms sends msg to the specified recipients
func SendSms(login, pass string, msg string, to string) (Quota, error) {
	client, err := newHTTPClient()
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to create httpClient")
	}

	loginner := defaultLoginner{client}
	if err = loginner.Login(login, pass); err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to login")
	}

	phoneNumbers, err := parsePhones(to)
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to parse `to` field")
	}

	smsSender := &httpSmsSender{client, &httpQuotaGetter{client}}
	return smsSender.SendSms(message(msg), phoneNumbers)
}
