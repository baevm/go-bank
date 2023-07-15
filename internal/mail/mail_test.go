package mail

import (
	"go-bank/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	cfg, err := config.Load("../..")
	require.NoError(t, err)

	sender := NewEmailSender(cfg.EMAIL_NAME, cfg.EMAIL_ADDRESS, cfg.EMAIL_PASSWORD)

	err = sender.SendEmail("Test email", "Test email content", []string{"user@test.com"}, nil, nil, nil)

	require.NoError(t, err)
}
