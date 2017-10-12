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
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	smsLeft, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.NoError(t, err)
	require.EqualValues(t, 2, smsLeft)

	msgPosted := false
	for _, val := range client.postValues {
		if val.Get("fieldMessage") == "msg" &&
			val.Get("fieldMsisdn") == "0601020304" {
			msgPosted = true
		}
	}
	require.True(t, msgPosted, "Message not posted")
}

func TestSmsSender_SendSms_unable_to_get_sms_left(t *testing.T) {
	client := &fakeClient{}
	sender := httpSmsSender{client, newFakeSmsGetter(0, errors.New("unable to get smsLeft"))}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.EqualError(t, err, "unable to get smsLeft")
}

func TestSmsSender_SendSms_sms_exceeded(t *testing.T) {
	client := &fakeClient{}
	sender := httpSmsSender{client, newFakeSmsGetter(int(NoSmsLeft), nil)}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.EqualError(t, err, "No SMS left")
}

func TestSmsSender_SendSms_sms_left_too_low_for_specified_amount_of_phone_numbers(t *testing.T) {
	client := &fakeClient{}
	sender := httpSmsSender{client, newFakeSmsGetter(1, nil)}

	_, err := sender.SendSms("msg", phoneNumbers{
		phoneNumber("0601010101"),
		phoneNumber("0602020202"),
	})

	require.EqualError(t, err,
		"too many phone numbers compared to SMS left (2 phone numbers, 1 SMS left)")
}

func TestSmsSender_SendSms_error_during_compose(t *testing.T) {
	client := newFakeClient("", errors.New("OMG"))
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.EqualError(t, err, "unable to compose sms: OMG")
}

func TestSmsSender_SendSms_no_validation_msg(t *testing.T) {
	client := newFakeClient("something else than expected validation token", nil)
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.Contains(t, err.Error(), "validation of message failed. Body: something else")
}

func TestSmsSender_SendSms_confirm_noSmsLeftAfter(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >1 SMS gratuits", "Votre message a bien été envoyé. 0 message gratuit"},
		errors: []error{nil, nil},
	}
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	smsLeft, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.NoError(t, err)
	require.EqualValues(t, 0, smsLeft)
}

func TestSmsSender_SendSms_error_during_confirmation(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "ignore"},
		errors: []error{nil, errors.New("an error")},
	}
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.EqualError(t, err, "an error")
}

func TestSmsSender_SendSms_failure_during_confirmation(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "something else than confirm token"},
		errors: []error{nil, nil},
	}
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.Contains(t, err.Error(), "unable to confirm message sending")
}

func TestSmsSender_SendSms_cannot_read_final_count_of_sms(t *testing.T) {
	client := &fakeClient{
		bodies: []string{">Validation< ... >3 SMS gratuits", "Votre message a bien été envoyé. NO SMS LEFT string"},
		errors: []error{nil, nil},
	}
	sender := httpSmsSender{client, newFakeSmsGetter(99, nil)}

	_, err := sender.SendSms("msg", phoneNumbers{phoneNumber("0601020304")})

	require.Contains(t, err.Error(), "unable to read SMS left")
}
