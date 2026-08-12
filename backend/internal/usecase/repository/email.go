package repository

import "context"

// EmailSender はメール送信に対する、usecase 側が要求する契約。
// 実装は adapter/external に置く。
type EmailSender interface {
	// Send は HTML メールを 1 通送信する。
	Send(ctx context.Context, to, subject, htmlBody string) error
}
