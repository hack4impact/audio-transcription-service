package transcription

import (
	"net/smtp" // mock
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the mockfmt mock package
	smtp.MOCK().SetController(ctrl)

	// Setup the ext mock package
	ext.MOCK().SetController(ctrl)

}
