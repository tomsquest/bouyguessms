package bouyguessms

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginner_Login_ok(t *testing.T) {
	loginner := &loginner{&fakeClient{
		errors: []error{nil, nil},
		bodies: []string{`<form method="post" action="/cas/login;jsessionid=123456.TC_PRD_AB">
		       <input type="hidden" name="lt" value="LT-1234-saucissionalail"/>`, "Success"},
	}}

	err := loginner.Login("login", "pass")

	require.NoError(t, err)
}

func TestLoginner_getTokens_err(t *testing.T) {
	loginner := &loginner{newFakeClient("", errors.New("an error"))}

	tokens, err := loginner.getTokens()

	require.Nil(t, tokens)
	require.EqualError(t, err, "an error")
}

func TestLoginner_getTokens_extracts_jsessionid_and_lt(t *testing.T) {
	expectedJsessionid := "123456.TC_PRD_AB"
	expectedLT := "LT-1234-saucissionalail"

	loginner := &loginner{newFakeClient(
		`<form method="post" action="/cas/login;jsessionid=123456.TC_PRD_AB">
		       <input type="hidden" name="lt" value="LT-1234-saucissionalail"/>`, nil)}

	tokens, err := loginner.getTokens()

	require.NoError(t, err)
	require.Equal(t, expectedJsessionid, tokens.jsessionid)
	require.Equal(t, expectedLT, tokens.lt)
}

func TestLoginner_getTokens_no_jsessionId(t *testing.T) {
	loginner := &loginner{newFakeClient(`<nojessionid />
		       <input type="hidden" name="lt" value="LT-1234-saucissionalail"/>`, nil)}

	tokens, err := loginner.getTokens()

	require.Nil(t, tokens)
	require.EqualError(t, err, "jessionid not found")
}

func TestLoginner_getTokens_no_LT(t *testing.T) {
	loginner := &loginner{newFakeClient(`<form method="post" action="/cas/login;jsessionid=123456.TC_PRD_AB">
        	       <nolt/>`, nil)}

	tokens, err := loginner.getTokens()

	require.Nil(t, tokens)
	require.EqualError(t, err, "lt token not found")
}

func TestLoginner_postLogin_ok(t *testing.T) {
	tokens := &tokens{
		jsessionid: "123",
		lt:         "456",
	}

	loginner := &loginner{newFakeClient("Success", nil)}

	err := loginner.postLogin("login", "pass", tokens)

	require.NoError(t, err)
}

func TestLoginner_postLogin_invalid_credentials(t *testing.T) {
	client := newFakeClient("Votre identifiant ou votre mot de passe est incorrect", nil)
	err := (&loginner{client}).postLogin("login", "pass", &tokens{})

	require.EqualError(t, err, "invalid credentials")
}

func TestLoginner_postLogin_remote_error(t *testing.T) {
	client := newFakeClient("", errors.New("an error"))
	err := (&loginner{client}).postLogin("login", "pass", &tokens{})

	require.EqualError(t, err, "an error")
}
