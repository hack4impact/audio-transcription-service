package transcription

import (
	"errors"
	"net/smtp" // mock
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	username = "test@email.com"
	password = "123456"
	host     = "smtp.gmail.com"
	port     = 25
	addr     = host + ":" + string(port)
	to       = []string{"to@email.com"}
	subject  = "subject"
	body     = "body"
)

func TestSendEmail(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the mock package
	smtp.MOCK().SetController(ctrl)

	gomock.InOrder(
		smtp.EXPECT().PlainAuth("", username, password, "smtp.gmail.com").Times(1),
		smtp.EXPECT().SendMail(addr, gomock.Any(), username, to, gomock.Any()).Times(1),
	)

	err := SendEmail(username, password, host, port, to, subject, body)
	assert.NoError(err)
}

func TestSendEmailReturnsError(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the mock package
	smtp.MOCK().SetController(ctrl)

	gomock.InOrder(
		smtp.EXPECT().PlainAuth("", username, password, "smtp.gmail.com"),
		smtp.EXPECT().SendMail(addr, gomock.Any(), username, to, gomock.Any()).Return(errors.New("Bad!")),
	)

	err := SendEmail(username, password, host, port, to, subject, body)
	assert.Error(err)
}
