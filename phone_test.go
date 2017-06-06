package bouyguessms

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPhoneNumber_Valid(t *testing.T) {
	require.True(t, PhoneNumber("0601020304").isValid())
	require.True(t, PhoneNumber("0701020304").isValid())
}

func TestPhoneNumber_Invalid(t *testing.T) {
	require.False(t, PhoneNumber("060102030").isValid())
	require.False(t, PhoneNumber("06010203040").isValid())

	require.False(t, PhoneNumber("0801020304").isValid())
	require.False(t, PhoneNumber("1601020304").isValid())

	require.False(t, PhoneNumber("0601020a04").isValid())
	require.False(t, PhoneNumber("0601020B04").isValid())
}

func TestParsePhones_single(t *testing.T) {
	phones, err := ParsePhones("0601010101")

	require.NoError(t, err)
	require.Len(t, phones, 1)
	require.Contains(t, phones, PhoneNumber("0601010101"))
}

func TestParsePhones_many(t *testing.T) {
	to, err := ParsePhones("0601010101;0602020202;0603030303;0604040404;0605050505")

	require.NoError(t, err)
	require.Len(t, to, 5)
}

func TestParsePhones_empty(t *testing.T) {
	phones, err := ParsePhones("")

	require.EqualError(t, err, "invalid phone number ")
	require.Nil(t, phones)
}

func TestParsePhones_tooMany(t *testing.T) {
	phones, err := ParsePhones("0601010101;0602020202;0603030303;0604040404;0605050505;0606060606")

	require.EqualError(t, err, "too many phone numbers given (5 is the max)")
	require.Nil(t, phones)
}
