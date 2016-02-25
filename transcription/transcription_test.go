package transcription

import (
	"net/smtp" // mock
	"os"       // mock
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the mock package
	smtp.MOCK().SetController(ctrl)
	os.MOCK().SetController(ctrl)

	mailEmail := "test@email.com"
	mailPassword := "123456"
	from := "from@email.com"
	to := []string{"to@email.com"}
	subject := "subject"
	body := "body"

	gomock.InOrder(
		os.EXPECT().Getenv("MAIL_EMAIL").Return(mailEmail).Times(1),
		os.EXPECT().Getenv("MAIL_PASSWORD").Return(mailPassword).Times(1),
		smtp.EXPECT().PlainAuth("", mailEmail, mailPassword, "smtp.gmail.com").Times(1),
		smtp.EXPECT().SendMail("smtp.gmail.com:25", gomock.Any(), from, to, gomock.Any()).Times(1),
	)

	sendEmail(from, to, subject, body)
}
