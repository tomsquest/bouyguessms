package bouyguessms

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuota_isExceeded(t *testing.T) {
	require.True(t, Quota(-1).IsExceeded())
	require.True(t, Quota(0).IsExceeded())
	require.False(t, Quota(1).IsExceeded())
}

func TestQuota_newExceeded_isExceedded(t *testing.T) {
	require.True(t, ExceededQuota.IsExceeded())
}
