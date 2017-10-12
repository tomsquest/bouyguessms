package bouyguessms

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSmsLeft_isExceeded(t *testing.T) {
	require.True(t, SmsLeft(-1).IsExceeded())
	require.True(t, SmsLeft(0).IsExceeded())
	require.False(t, SmsLeft(1).IsExceeded())
}

func TestSmsLeft_newExceeded_isExceedded(t *testing.T) {
	require.True(t, NoSmsLeft.IsExceeded())
}
