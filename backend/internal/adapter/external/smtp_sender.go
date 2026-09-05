package external

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"github.com/norman6464/devsync/backend/internal/infra/config"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// smtpEmailSender は [repository.EmailSender] の SMTP 実装。
type smtpEmailSender struct {
	cfg *config.Config
}

// NewSMTPEmailSender は EmailSender の SMTP 実装を返す。
func NewSMTPEmailSender(cfg *config.Config) repository.EmailSender {
	return &smtpEmailSender{cfg: cfg}
}

var _ repository.EmailSender = (*smtpEmailSender)(nil)

// Send は SMTP サーバーを通じて HTML メールを送信する。
// smtp.SendMail が ctx を受け取らないため、ctx はキャンセルの伝播には使えない。
func (s *smtpEmailSender) Send(ctx context.Context, to, subject, htmlBody string) error {
	from := s.cfg.EmailFrom
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)

	// MIMEヘッダーを構築
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	// SMTP認証（ユーザー名が空の場合は認証なし）
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	}

	return smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes())
}
