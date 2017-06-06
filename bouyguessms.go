package bouyguessms

import "github.com/pkg/errors"

func GetQuota(login, pass string) (Quota, error) {
	client, err := NewHttpClient()
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to create httpClient")
	}

	loginner := loginner{client}
	if err = loginner.Login(login, pass); err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to login")
	}

	quotaGetter := &quotaGetter{client}
	return quotaGetter.Get()
}

func SendSms(login, pass string, msg string, to string) (Quota, error) {
	client, err := NewHttpClient()
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to create httpClient")
	}

	loginner := loginner{client}
	if err = loginner.Login(login, pass); err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to login")
	}

	phoneNumbers, err := ParsePhones(to)
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to parse `to` field")
	}

	smsSender := &smsSender{client, &quotaGetter{client}}
	return smsSender.SendSms(Msg(msg), phoneNumbers)
}
