package bouyguessms

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsg_String_small(t *testing.T) {
	msg := message("it fits")

	require.Equal(t, msg.String(), "it fits")
}

func TestMsg_String_truncateTooLong(t *testing.T) {
	msg := message("more than max size. more than max size. more than max size. more than max size. more than max size. more than max size. more than max size. more than max size. MORE")

	require.Equal(t, msg.String(), "more than max size. more than max size. more than max size. more than max size. more than max size. more than max size. more than max size. more than max size. ")
}
