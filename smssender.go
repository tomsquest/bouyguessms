package bouyguessms

import (
	"github.com/pkg/errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type SmsSender interface {
	SendSms(msg string, to string) error
}

type smsSender struct {
	client      HttpClient
	quotaGetter QuotaGetter
}

func (sender *smsSender) SendSms(msg Msg, phoneNumbers PhoneNumbers) (Quota, error) {
	err := sender.checkSmsLeft(phoneNumbers)
	if err != nil {
		return ExceededQuota, err
	}

	quota, err := sender.composeMessage(msg, phoneNumbers)
	if err != nil {
		return quota, err
	}

	return quota, sender.confirmMessage()
}

func (sender *smsSender) checkSmsLeft(phonenumbers PhoneNumbers) error {
	quota, err := sender.quotaGetter.Get()
	if err != nil {
		return err
	}

	if quota.IsExceeded() {
		return errors.New("quota exceeded")
	}

	if quota.Remaining() < len(phonenumbers) {
		return errors.Errorf("too many phone numbers compared to quota left (%d phone numbers, %d SMS left)", len(phonenumbers), quota)
	}

	return err
}

func (sender *smsSender) composeMessage(msg Msg, phoneNumbers PhoneNumbers) (Quota, error) {
	msgForm := make(url.Values)
	msgForm.Add("fieldMsisdn", phoneNumbers.String())
	msgForm.Add("fieldMessage", msg.String())
	msgForm.Add("Verif.x", "51")
	msgForm.Add("Verif.y", "16")

	sendMessageUrl := "https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/confirmSendSMS.phtml"
	body, err := sender.client.PostForm(sendMessageUrl, msgForm)
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to compose sms")
	}

	if !strings.Contains(body, ">Validation<") {
		return ExceededQuota, errors.Errorf("validation of message failed. Body: %s", body)
	}

	regex := regexp.MustCompile(">(\\d*) SMS gratuit")
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ExceededQuota, errors.New("unable to read SMS left")
	}

	quota, err := strconv.Atoi(matches[1])
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to convert quota to int")
	}

	return Quota(quota), nil
}

func (sender *smsSender) confirmMessage() error {
	body, err := sender.client.Get("https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/resultSendSMS.phtml")
	if err != nil {
		return err
	}

	if !strings.Contains(body, "Votre message a bien été envoyé au numéro") {
		return errors.Errorf("unable to confirm message sending (last step results in unattented message). Body: %s", body)
	}

	return nil
}
