package bouyguessms

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuotaGetter_Get_ok(t *testing.T) {
	getter := httpQuotaGetter{
		newFakeClient("Il vous reste <strong>1 SMS gratuits", nil),
	}

	quota, err := getter.Get()

	require.NoError(t, err)
	require.EqualValues(t, Quota(1), quota)
}

func TestQuotaGetter_Get_client_error(t *testing.T) {
	getter := httpQuotaGetter{
		newFakeClient("some error", errors.New("an error")),
	}

	quota, err := getter.Get()

	require.EqualError(t, err, "an error")
	require.True(t, quota.IsExceeded())
}

func TestQuotaGetter_Get_pattern_not_found(t *testing.T) {
	getter := httpQuotaGetter{
		newFakeClient("not quota pattern", nil),
	}

	quota, err := getter.Get()

	require.EqualError(t, err, "unable to read quota from body")
	require.True(t, quota.IsExceeded())
}
