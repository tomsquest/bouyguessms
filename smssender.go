package bouyguessms

import (
	"github.com/pkg/errors"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type smsSender interface {
	SendSms(msg string, to string) error
}

type httpSmsSender struct {
	client        httpClient
	smsLeftGetter smsLeftGetter
}

func (sender *httpSmsSender) SendSms(msg message, phoneNumbers phoneNumbers) (SmsLeft, error) {
	err := sender.checkSmsLeft(phoneNumbers)
	if err != nil {
		return NoSmsLeft, err
	}

	err = sender.composeMessage(msg, phoneNumbers)
	if err != nil {
		return NoSmsLeft, err
	}

	return sender.confirmMessage()
}

func (sender *httpSmsSender) checkSmsLeft(phonenumbers phoneNumbers) error {
	log.Println("Checking SMS left...")
	smsLeft, err := sender.smsLeftGetter.Get()
	if err != nil {
		return err
	}

	if smsLeft.IsExceeded() {
		return errors.New("No SMS left")
	}

	if int(smsLeft) < len(phonenumbers) {
		return errors.Errorf("too many phone numbers compared to SMS left (%d phone numbers, %d SMS left)", len(phonenumbers), smsLeft)
	}

	log.Printf("SMS left before sending: %d\n", smsLeft)
	return err
}

func (sender *httpSmsSender) composeMessage(msg message, phoneNumbers phoneNumbers) error {
	log.Printf("Composing message to %s\n", phoneNumbers)

	msgForm := make(url.Values)
	msgForm.Add("fieldMsisdn", phoneNumbers.String())
	msgForm.Add("fieldMessage", msg.String())
	msgForm.Add("Verif.x", "51")
	msgForm.Add("Verif.y", "16")

	sendMessageURL := "https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/confirmSendSMS.phtml"
	body, err := sender.client.PostForm(sendMessageURL, msgForm)
	if err != nil {
		return errors.Wrap(err, "unable to compose sms")
	}

	if !strings.Contains(body, ">Validation<") {
		return errors.Errorf("validation of message failed. Body: %s", body)
	}

	log.Println("Message composed")
	return nil
}

func (sender *httpSmsSender) confirmMessage() (SmsLeft, error) {
	log.Println("Confirming message...")
	body, err := sender.client.Get("https://www.secure.bbox.bouyguestelecom.fr/services/SMSIHD/resultSendSMS.phtml")
	if err != nil {
		return NoSmsLeft, err
	}

	if !strings.Contains(body, "Votre message a bien été envoyé") {
		return NoSmsLeft, errors.Errorf("unable to confirm message sending (last step results in unattended message). Body: %s", body)
	}

	regex := regexp.MustCompile("(\\d*) message(s)? gratuit(s)?")
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return NoSmsLeft, errors.Errorf("unable to read SMS left. Body: %s", body)
	}

	smsLeft, err := strconv.Atoi(matches[1])
	if err != nil {
		return NoSmsLeft, errors.Wrapf(err, "unable to convert %v to int", matches[1])
	}

	log.Printf("Message confirmed. SMS left after sending: %d", smsLeft)
	return SmsLeft(smsLeft), nil
}
