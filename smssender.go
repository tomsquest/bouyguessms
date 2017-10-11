package bouyguessms

import (
	"github.com/pkg/errors"
	"log"
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

	err = sender.composeMessage(msg, phoneNumbers)
	if err != nil {
		return ExceededQuota, err
	}

	return sender.confirmMessage()
}

func (sender *smsSender) checkSmsLeft(phonenumbers PhoneNumbers) error {
	log.Println("Checking quota left...")
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

	log.Printf("Quota left before sending: %d\n", quota)
	return err
}

func (sender *smsSender) composeMessage(msg Msg, phoneNumbers PhoneNumbers) error {
	log.Printf("Composing message to %s\n", phoneNumbers)

	msgForm := make(url.Values)
	msgForm.Add("fieldMsisdn", phoneNumbers.String())
	msgForm.Add("fieldMessage", msg.String())
	msgForm.Add("Verif.x", "51")
	msgForm.Add("Verif.y", "16")

	sendMessageUrl := "https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/confirmSendSMS.phtml"
	body, err := sender.client.PostForm(sendMessageUrl, msgForm)
	if err != nil {
		return errors.Wrap(err, "unable to compose sms")
	}

	if !strings.Contains(body, ">Validation<") {
		return errors.Errorf("validation of message failed. Body: %s", body)
	}

	log.Println("Message composed")
	return nil
}

func (sender *smsSender) confirmMessage() (Quota, error) {
	log.Println("Confirming message...")
	body, err := sender.client.Get("https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/resultSendSMS.phtml")
	if err != nil {
		return ExceededQuota, err
	}

	if !strings.Contains(body, "Votre message a bien été envoyé") {
		return ExceededQuota, errors.Errorf("unable to confirm message sending (last step results in unattended message). Body: %s", body)
	}

	regex := regexp.MustCompile("(\\d*) messages gratuits")
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ExceededQuota, errors.Errorf("unable to read SMS left. Body: %s", body)
	}

	quota, err := strconv.Atoi(matches[1])
	if err != nil {
		return ExceededQuota, errors.Wrap(err, "unable to convert quota to int")
	}

	log.Printf("Message confirmed. Quota left after sending: %d", quota)
	return Quota(quota), nil
}
