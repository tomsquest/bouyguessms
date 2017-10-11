package bouyguessms

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSmsSender_SendSms_ok(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "Votre message a bien été envoyé au numéro. 2 messages gratuits"},
		errors: []error{nil, nil},
	}
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	quota, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.NoError(t, err)
	require.EqualValues(t, 2, quota)

	msgPosted := false
	for _, val := range client.postValues {
		if val.Get("fieldMessage") == "msg" &&
			val.Get("fieldMsisdn") == "0601020304" {
			msgPosted = true
		}
	}
	require.True(t, msgPosted, "Message not posted")
}

func TestSmsSender_SendSms_unable_to_get_quota(t *testing.T) {
	client := &fakeClient{}
	sender := smsSender{client, NewFakeQuotaGetter(0, errors.New("unable to get quota"))}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.EqualError(t, err, "unable to get quota")
}

func TestSmsSender_SendSms_quota_exceeded(t *testing.T) {
	client := &fakeClient{}
	sender := smsSender{client, NewFakeQuotaGetter(int(ExceededQuota), nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.EqualError(t, err, "quota exceeded")
}

func TestSmsSender_SendSms_quota_too_low_for_specified_amount_of_phone_numbers(t *testing.T) {
	client := &fakeClient{}
	sender := smsSender{client, NewFakeQuotaGetter(1, nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{
		PhoneNumber("0601010101"),
		PhoneNumber("0602020202"),
	})

	require.EqualError(t, err,
		"too many phone numbers compared to quota left (2 phone numbers, 1 SMS left)")
}

func TestSmsSender_SendSms_error_during_compose(t *testing.T) {
	client := newFakeClient("", errors.New("OMG"))
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.EqualError(t, err, "unable to compose sms: OMG")
}

func TestSmsSender_SendSms_no_validation_msg(t *testing.T) {
	client := newFakeClient("something else than expected validation token", nil)
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.Contains(t, err.Error(), "validation of message failed. Body: something else")
}

func TestSmsSender_SendSms_confirm_noSmsLeftAfter(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >1 SMS gratuits", "Votre message a bien été envoyé. 0 message gratuit"},
		errors: []error{nil, nil},
	}
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	quota, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.NoError(t, err)
	require.EqualValues(t, 0, quota)
}

func TestSmsSender_SendSms_error_during_confirmation(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "ignore"},
		errors: []error{nil, errors.New("an error")},
	}
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.EqualError(t, err, "an error")
}

func TestSmsSender_SendSms_failure_during_confirmation(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "something else than confirm token"},
		errors: []error{nil, nil},
	}
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.Contains(t, err.Error(), "unable to confirm message sending")
}

func TestSmsSender_SendSms_cannot_read_final_quota(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "Votre message a bien été envoyé. NO QUOTA string"},
		errors: []error{nil, nil},
	}
	sender := smsSender{client, NewFakeQuotaGetter(99, nil)}

	_, err := sender.SendSms("msg", PhoneNumbers{PhoneNumber("0601020304")})

	require.Contains(t, err.Error(), "unable to read SMS left")
}
