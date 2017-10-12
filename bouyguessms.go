package bouyguessms

import "github.com/pkg/errors"

// GetSmsLeft fetches the remaining number of SMS
func GetSmsLeft(login, pass string) (SmsLeft, error) {
	client, err := newHTTPClient()
	if err != nil {
		return NoSmsLeft, errors.Wrap(err, "unable to create httpClient")
	}

	loginner := httpLoginner{client}
	if err = loginner.Login(login, pass); err != nil {
		return NoSmsLeft, errors.Wrap(err, "unable to login")
	}

	smsLeftGetter := &httpSmsLeftGetter{client}
	return smsLeftGetter.Get()
}

// SendSms sends msg to the specified recipients
func SendSms(login, pass string, msg string, to string) (SmsLeft, error) {
	client, err := newHTTPClient()
	if err != nil {
		return NoSmsLeft, errors.Wrap(err, "unable to create httpClient")
	}

	loginner := httpLoginner{client}
	if err = loginner.Login(login, pass); err != nil {
		return NoSmsLeft, errors.Wrap(err, "unable to login")
	}

	phoneNumbers, err := parsePhones(to)
	if err != nil {
		return NoSmsLeft, errors.Wrap(err, "unable to parse `to` field")
	}

	smsSender := &httpSmsSender{client, &httpSmsLeftGetter{client}}
	return smsSender.SendSms(message(msg), phoneNumbers)
}
