package bouyguessms

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type loginner interface {
	Login(login, pass string) error
}

type httpLoginner struct {
	client httpClient
}

func (l *httpLoginner) Login(login, pass string) error {
	tokens, err := l.getTokens()
	if err != nil {
		return err
	}

	return l.postLogin(login, pass, tokens)
}

type tokens struct {
	jsessionid string
	lt         string
}

func (l *httpLoginner) getTokens() (*tokens, error) {
	body, err := l.client.Get("https://www.mon-compte.bouyguestelecom.fr/cas/login")
	if err != nil {
		return nil, err
	}

	jsessionid, err := l.extractJsessionid(body)
	if err != nil {
		return nil, err
	}

	lt, err := l.extractLT(body)
	if err != nil {
		return nil, err
	}

	return &tokens{jsessionid, lt}, nil
}

func (l *httpLoginner) extractJsessionid(body string) (string, error) {
	regex := regexp.MustCompile("(?i:jsessionid)=(.+?)\"")
	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", errors.New("jessionid not found")
}

func (l *httpLoginner) extractLT(body string) (string, error) {
	regex := regexp.MustCompile("name=\"lt\" value=\"(.+?)\"")
	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", errors.New("lt token not found")
}

func (l *httpLoginner) postLogin(login, pass string, tokens *tokens) error {
	loginURL := "https://www.mon-compte.bouyguestelecom.fr/cas/login;jsessionid=" + tokens.jsessionid + "?service=https%3A%2F%2Fwww.secure.bbox.bouyguestelecom.fr%2Fservices%2FSMSIHD%2FsendSMS.phtml"

	data := make(url.Values)
	data.Add("username", login)
	data.Add("password", pass)
	data.Add("rememberMe", "true")
	data.Add("_rememberMe", "on")
	data.Add("lt", tokens.lt)
	data.Add("execution", "e1s1")
	data.Add("_eventId", "submit")

	body, err := l.client.PostForm(loginURL, data)
	if err != nil {
		return err
	}

	if strings.Contains(body, "Votre identifiant ou votre mot de passe est incorrect") {
		return errors.New("invalid credentials")
	}

	return nil
}
