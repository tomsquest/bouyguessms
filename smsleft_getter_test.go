package bouyguessms

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSmsLeftGetter_Get_ok(t *testing.T) {
	getter := httpSmsLeftGetter{
		newFakeClient("Il vous reste <strong>1 SMS gratuits", nil),
	}

	smsLeft, err := getter.Get()

	require.NoError(t, err)
	require.EqualValues(t, SmsLeft(1), smsLeft)
}

func TestSmsLeftGetter_Get_client_error(t *testing.T) {
	getter := httpSmsLeftGetter{
		newFakeClient("some error", errors.New("an error")),
	}

	smsLeft, err := getter.Get()

	require.EqualError(t, err, "an error")
	require.True(t, smsLeft.IsExceeded())
}

func TestSmsLeftGetter_Get_pattern_not_found(t *testing.T) {
	getter := httpSmsLeftGetter{
		newFakeClient("not smsLeft pattern", nil),
	}

	smsLeft, err := getter.Get()

	require.EqualError(t, err, "unable to read SMS left from body")
	require.True(t, smsLeft.IsExceeded())
}
