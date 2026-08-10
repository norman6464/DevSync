package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"

	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	usecaserepo "github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// EmailSenderInterface はメール送信の抽象化インターフェース。
// テスト時にモック実装に差し替え可能にするために使用する。
type EmailSenderInterface interface {
	Send(to, subject, htmlBody string) error
}

// SMTPEmailSender はSMTP経由でメールを送信する本番用実装。
type SMTPEmailSender struct {
	cfg *config.Config
}

// NewSMTPEmailSender は新しいSMTPEmailSenderインスタンスを生成する。
func NewSMTPEmailSender(cfg *config.Config) *SMTPEmailSender {
	return &SMTPEmailSender{cfg: cfg}
}

// Send はSMTPサーバーを通じてHTMLメールを送信する。
func (s *SMTPEmailSender) Send(to, subject, htmlBody string) error {
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

// WeeklyReportEmailService はウィークリーレポートメールの生成・送信を管理する。
// レポートデータの取得、HTMLレンダリング、メール送信のオーケストレーションを担当する。
type WeeklyReportEmailService struct {
	sender   EmailSenderInterface
	reports  usecaserepo.WeeklyActivityReportReader
	userRepo repository.UserRepositoryInterface
	appURL   string
	tmpl     *template.Template
}

// NewWeeklyReportEmailService は新しいWeeklyReportEmailServiceインスタンスを生成する。
func NewWeeklyReportEmailService(sender EmailSenderInterface, reports usecaserepo.WeeklyActivityReportReader, userRepo repository.UserRepositoryInterface) *WeeklyReportEmailService {
	svc := &WeeklyReportEmailService{
		sender:   sender,
		reports:  reports,
		userRepo: userRepo,
		appURL:   "http://localhost:5173",
	}

	// HTMLテンプレートをパース
	svc.tmpl = template.Must(template.New("weekly_report").Parse(weeklyReportTemplate))

	return svc
}

// SetAppURL はメール内で使用するアプリケーションURLを設定する。
func (s *WeeklyReportEmailService) SetAppURL(url string) {
	s.appURL = url
}

// SendWeeklyReport は指定ユーザーにウィークリーレポートメールを送信する。
func (s *WeeklyReportEmailService) SendWeeklyReport(user *model.User, report *model.ActivityReport) error {
	if user.Email == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "メールアドレスが空です", nil)
	}
	if report == nil {
		return domain.NewError(domain.ErrCodeBadRequest, "レポートがありません", nil)
	}

	lang := user.EmailLanguage
	if lang == "" {
		lang = "ja"
	}

	html, err := s.RenderHTML(user, report, lang)
	if err != nil {
		return domain.NewError(domain.ErrCodeInternal, "メールテンプレートのレンダリングに失敗", err)
	}

	texts := getEmailTexts(lang)
	subject := fmt.Sprintf("[DevSync] %s", texts["subject"])

	return s.sender.Send(user.Email, subject, html)
}

// SendAllWeeklyReports は全対象ユーザーにウィークリーレポートメールを一括送信する。
// メール配信が無効なユーザーはスキップし、1ユーザーのエラーで他は止まらない。
func (s *WeeklyReportEmailService) SendAllWeeklyReports() error {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return domain.NewError(domain.ErrCodeDatabase, "ユーザー一覧の取得に失敗", err)
	}

	for _, user := range users {
		// メール配信が無効なユーザーはスキップ
		if !user.EmailWeeklyReport {
			continue
		}

		// 定時バッチからの呼び出しでリクエスト ctx を持たないため、ここが ctx の起点になる
		report, err := s.reports.GetWeeklyReport(context.Background(), user.ID)
		if err != nil {
			log.Printf("ウィークリーレポート生成失敗 (userID=%d): %v", user.ID, err)
			continue
		}

		if err := s.SendWeeklyReport(&user, report); err != nil {
			log.Printf("ウィークリーレポートメール送信失敗 (userID=%d): %v", user.ID, err)
			continue
		}

		log.Printf("ウィークリーレポートメール送信成功 (userID=%d)", user.ID)
	}

	return nil
}

// emailTemplateData はHTMLメールテンプレートに渡すデータ構造体。
type emailTemplateData struct {
	UserName           string
	StartDate          string
	EndDate            string
	TotalContributions int
	PostsCreated       int
	CommentsCreated    int
	LikesReceived      int
	GoalsCompleted     int
	GoalsProgress      int
	NewFollowers       int
	MessagesExchanged  int
	TopLanguages       []model.LanguageActivity
	AppURL             string
	Texts              map[string]string
}

// RenderHTML はレポートデータからHTMLメール本文をレンダリングする。
func (s *WeeklyReportEmailService) RenderHTML(user *model.User, report *model.ActivityReport, lang string) (string, error) {
	texts := getEmailTexts(lang)

	data := emailTemplateData{
		UserName:           user.Name,
		StartDate:          report.StartDate.Format("2006/01/02"),
		EndDate:            report.EndDate.Format("2006/01/02"),
		TotalContributions: report.TotalContributions,
		PostsCreated:       report.PostsCreated,
		CommentsCreated:    report.CommentsCreated,
		LikesReceived:      report.LikesReceived,
		GoalsCompleted:     report.GoalsCompleted,
		GoalsProgress:      report.GoalsProgress,
		NewFollowers:       report.NewFollowers,
		MessagesExchanged:  report.MessagesExchanged,
		TopLanguages:       report.TopLanguages,
		AppURL:             s.appURL,
		Texts:              texts,
	}

	var buf bytes.Buffer
	if err := s.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getEmailTexts は指定言語のメールテキストマップを返す。
func getEmailTexts(lang string) map[string]string {
	texts, ok := emailTextsMap[lang]
	if !ok {
		texts = emailTextsMap["ja"]
	}
	return texts
}

// emailTextsMap はメールテンプレート用の多言語テキスト。
var emailTextsMap = map[string]map[string]string{
	"ja": {
		"subject":            "ウィークリーアクティビティレポート",
		"greeting":           "さん、今週のアクティビティレポートです",
		"period":             "集計期間",
		"summary":            "今週のサマリー",
		"totalContributions": "コントリビューション",
		"postsCreated":       "投稿",
		"commentsCreated":    "コメント",
		"likesReceived":      "いいね",
		"goalsCompleted":     "目標達成",
		"goalsProgress":      "目標進捗",
		"newFollowers":       "新規フォロワー",
		"messagesExchanged":  "メッセージ",
		"topLanguages":       "トップ言語",
		"viewFullReport":     "レポートを詳しく見る",
		"unsubscribe":        "配信停止はこちら",
		"footer":             "このメールはDevSyncから自動送信されています。",
	},
	"en": {
		"subject":            "Weekly Activity Report",
		"greeting":           ", here's your weekly activity report",
		"period":             "Period",
		"summary":            "This Week's Summary",
		"totalContributions": "Contributions",
		"postsCreated":       "Posts",
		"commentsCreated":    "Comments",
		"likesReceived":      "Likes",
		"goalsCompleted":     "Goals Completed",
		"goalsProgress":      "Goals Progress",
		"newFollowers":       "New Followers",
		"messagesExchanged":  "Messages",
		"topLanguages":       "Top Languages",
		"viewFullReport":     "View Full Report",
		"unsubscribe":        "Unsubscribe",
		"footer":             "This email was sent automatically by DevSync.",
	},
	"ko": {
		"subject":            "주간 활동 보고서",
		"greeting":           "님, 이번 주 활동 보고서입니다",
		"period":             "기간",
		"summary":            "이번 주 요약",
		"totalContributions": "기여",
		"postsCreated":       "게시물",
		"commentsCreated":    "댓글",
		"likesReceived":      "좋아요",
		"goalsCompleted":     "목표 달성",
		"goalsProgress":      "목표 진행률",
		"newFollowers":       "새 팔로워",
		"messagesExchanged":  "메시지",
		"topLanguages":       "인기 언어",
		"viewFullReport":     "전체 보고서 보기",
		"unsubscribe":        "구독 취소",
		"footer":             "이 이메일은 DevSync에서 자동으로 발송되었습니다.",
	},
	"zh-CN": {
		"subject":            "每周活动报告",
		"greeting":           "，这是您的每周活动报告",
		"period":             "统计期间",
		"summary":            "本周摘要",
		"totalContributions": "贡献",
		"postsCreated":       "帖子",
		"commentsCreated":    "评论",
		"likesReceived":      "点赞",
		"goalsCompleted":     "目标完成",
		"goalsProgress":      "目标进度",
		"newFollowers":       "新粉丝",
		"messagesExchanged":  "消息",
		"topLanguages":       "热门语言",
		"viewFullReport":     "查看完整报告",
		"unsubscribe":        "取消订阅",
		"footer":             "此邮件由DevSync自动发送。",
	},
	"zh-TW": {
		"subject":            "每週活動報告",
		"greeting":           "，這是您的每週活動報告",
		"period":             "統計期間",
		"summary":            "本週摘要",
		"totalContributions": "貢獻",
		"postsCreated":       "帖子",
		"commentsCreated":    "評論",
		"likesReceived":      "按讚",
		"goalsCompleted":     "目標完成",
		"goalsProgress":      "目標進度",
		"newFollowers":       "新粉絲",
		"messagesExchanged":  "訊息",
		"topLanguages":       "熱門語言",
		"viewFullReport":     "查看完整報告",
		"unsubscribe":        "取消訂閱",
		"footer":             "此郵件由DevSync自動發送。",
	},
	"es": {
		"subject":            "Informe de Actividad Semanal",
		"greeting":           ", aquí está tu informe semanal",
		"period":             "Período",
		"summary":            "Resumen de Esta Semana",
		"totalContributions": "Contribuciones",
		"postsCreated":       "Publicaciones",
		"commentsCreated":    "Comentarios",
		"likesReceived":      "Me gusta",
		"goalsCompleted":     "Metas Completadas",
		"goalsProgress":      "Progreso de Metas",
		"newFollowers":       "Nuevos Seguidores",
		"messagesExchanged":  "Mensajes",
		"topLanguages":       "Lenguajes Principales",
		"viewFullReport":     "Ver Informe Completo",
		"unsubscribe":        "Cancelar suscripción",
		"footer":             "Este correo fue enviado automáticamente por DevSync.",
	},
	"fr": {
		"subject":            "Rapport d'Activité Hebdomadaire",
		"greeting":           ", voici votre rapport d'activité",
		"period":             "Période",
		"summary":            "Résumé de la Semaine",
		"totalContributions": "Contributions",
		"postsCreated":       "Publications",
		"commentsCreated":    "Commentaires",
		"likesReceived":      "J'aime",
		"goalsCompleted":     "Objectifs Atteints",
		"goalsProgress":      "Progression des Objectifs",
		"newFollowers":       "Nouveaux Abonnés",
		"messagesExchanged":  "Messages",
		"topLanguages":       "Langages Principaux",
		"viewFullReport":     "Voir le Rapport Complet",
		"unsubscribe":        "Se désabonner",
		"footer":             "Cet email a été envoyé automatiquement par DevSync.",
	},
	"de": {
		"subject":            "Wöchentlicher Aktivitätsbericht",
		"greeting":           ", hier ist dein Aktivitätsbericht",
		"period":             "Zeitraum",
		"summary":            "Zusammenfassung der Woche",
		"totalContributions": "Beiträge",
		"postsCreated":       "Beiträge",
		"commentsCreated":    "Kommentare",
		"likesReceived":      "Likes",
		"goalsCompleted":     "Ziele Erreicht",
		"goalsProgress":      "Zielfortschritt",
		"newFollowers":       "Neue Follower",
		"messagesExchanged":  "Nachrichten",
		"topLanguages":       "Top Sprachen",
		"viewFullReport":     "Vollständigen Bericht Anzeigen",
		"unsubscribe":        "Abbestellen",
		"footer":             "Diese E-Mail wurde automatisch von DevSync gesendet.",
	},
	"pt": {
		"subject":            "Relatório de Atividade Semanal",
		"greeting":           ", aqui está seu relatório semanal",
		"period":             "Período",
		"summary":            "Resumo da Semana",
		"totalContributions": "Contribuições",
		"postsCreated":       "Publicações",
		"commentsCreated":    "Comentários",
		"likesReceived":      "Curtidas",
		"goalsCompleted":     "Metas Concluídas",
		"goalsProgress":      "Progresso das Metas",
		"newFollowers":       "Novos Seguidores",
		"messagesExchanged":  "Mensagens",
		"topLanguages":       "Linguagens Principais",
		"viewFullReport":     "Ver Relatório Completo",
		"unsubscribe":        "Cancelar inscrição",
		"footer":             "Este email foi enviado automaticamente pelo DevSync.",
	},
	"ru": {
		"subject":            "Еженедельный Отчёт об Активности",
		"greeting":           ", ваш еженедельный отчёт",
		"period":             "Период",
		"summary":            "Итоги Недели",
		"totalContributions": "Вклады",
		"postsCreated":       "Публикации",
		"commentsCreated":    "Комментарии",
		"likesReceived":      "Лайки",
		"goalsCompleted":     "Цели Достигнуты",
		"goalsProgress":      "Прогресс Целей",
		"newFollowers":       "Новые Подписчики",
		"messagesExchanged":  "Сообщения",
		"topLanguages":       "Популярные Языки",
		"viewFullReport":     "Посмотреть Полный Отчёт",
		"unsubscribe":        "Отписаться",
		"footer":             "Это письмо было отправлено автоматически DevSync.",
	},
}

// weeklyReportTemplate はウィークリーレポートメールのHTMLテンプレート。
var weeklyReportTemplate = strings.TrimSpace(`
<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DevSync Weekly Report</title>
</head>
<body style="margin:0;padding:0;background-color:#111827;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="background-color:#111827;">
<tr><td align="center" style="padding:24px 16px;">
<table role="presentation" width="600" cellpadding="0" cellspacing="0" style="background-color:#1f2937;border-radius:12px;overflow:hidden;">

<!-- Header -->
<tr><td style="background:linear-gradient(135deg,#3b82f6,#8b5cf6);padding:32px 24px;text-align:center;">
<h1 style="color:#ffffff;font-size:24px;margin:0 0 8px;">DevSync</h1>
<p style="color:#e0e7ff;font-size:14px;margin:0;">{{.UserName}}{{.Texts.greeting}}</p>
<p style="color:#c7d2fe;font-size:12px;margin:8px 0 0;">{{.Texts.period}}: {{.StartDate}} 〜 {{.EndDate}}</p>
</td></tr>

<!-- Summary -->
<tr><td style="padding:24px;">
<h2 style="color:#f9fafb;font-size:18px;margin:0 0 16px;">{{.Texts.summary}}</h2>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0">
<tr>
<td width="50%" style="padding:8px;">
<div style="background-color:#374151;border-radius:8px;padding:16px;text-align:center;">
<p style="color:#9ca3af;font-size:11px;margin:0 0 4px;text-transform:uppercase;">{{.Texts.totalContributions}}</p>
<p style="color:#60a5fa;font-size:28px;font-weight:bold;margin:0;">{{.TotalContributions}}</p>
</div>
</td>
<td width="50%" style="padding:8px;">
<div style="background-color:#374151;border-radius:8px;padding:16px;text-align:center;">
<p style="color:#9ca3af;font-size:11px;margin:0 0 4px;text-transform:uppercase;">{{.Texts.postsCreated}}</p>
<p style="color:#34d399;font-size:28px;font-weight:bold;margin:0;">{{.PostsCreated}}</p>
</div>
</td>
</tr>
<tr>
<td width="50%" style="padding:8px;">
<div style="background-color:#374151;border-radius:8px;padding:16px;text-align:center;">
<p style="color:#9ca3af;font-size:11px;margin:0 0 4px;text-transform:uppercase;">{{.Texts.commentsCreated}}</p>
<p style="color:#fbbf24;font-size:28px;font-weight:bold;margin:0;">{{.CommentsCreated}}</p>
</div>
</td>
<td width="50%" style="padding:8px;">
<div style="background-color:#374151;border-radius:8px;padding:16px;text-align:center;">
<p style="color:#9ca3af;font-size:11px;margin:0 0 4px;text-transform:uppercase;">{{.Texts.likesReceived}}</p>
<p style="color:#f87171;font-size:28px;font-weight:bold;margin:0;">{{.LikesReceived}}</p>
</div>
</td>
</tr>
<tr>
<td width="50%" style="padding:8px;">
<div style="background-color:#374151;border-radius:8px;padding:16px;text-align:center;">
<p style="color:#9ca3af;font-size:11px;margin:0 0 4px;text-transform:uppercase;">{{.Texts.goalsCompleted}}</p>
<p style="color:#a78bfa;font-size:28px;font-weight:bold;margin:0;">{{.GoalsCompleted}}</p>
</div>
</td>
<td width="50%" style="padding:8px;">
<div style="background-color:#374151;border-radius:8px;padding:16px;text-align:center;">
<p style="color:#9ca3af;font-size:11px;margin:0 0 4px;text-transform:uppercase;">{{.Texts.newFollowers}}</p>
<p style="color:#fb923c;font-size:28px;font-weight:bold;margin:0;">{{.NewFollowers}}</p>
</div>
</td>
</tr>
</table>
</td></tr>

{{if .TopLanguages}}
<!-- Top Languages -->
<tr><td style="padding:0 24px 24px;">
<h3 style="color:#f9fafb;font-size:16px;margin:0 0 12px;">{{.Texts.topLanguages}}</h3>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0">
{{range .TopLanguages}}
<tr><td style="padding:4px 0;">
<span style="color:#d1d5db;font-size:14px;">{{.Language}}</span>
<span style="color:#6b7280;font-size:12px;float:right;">{{.Repos}} repos</span>
</td></tr>
{{end}}
</table>
</td></tr>
{{end}}

<!-- CTA -->
<tr><td style="padding:0 24px 32px;text-align:center;">
<a href="{{.AppURL}}/reports" style="display:inline-block;background:linear-gradient(135deg,#3b82f6,#8b5cf6);color:#ffffff;text-decoration:none;padding:12px 32px;border-radius:8px;font-size:14px;font-weight:600;">
{{.Texts.viewFullReport}}
</a>
</td></tr>

<!-- Footer -->
<tr><td style="background-color:#111827;padding:24px;text-align:center;border-top:1px solid #374151;">
<p style="color:#6b7280;font-size:12px;margin:0 0 8px;">{{.Texts.footer}}</p>
<a href="{{.AppURL}}/settings" style="color:#60a5fa;font-size:12px;">{{.Texts.unsubscribe}}</a>
</td></tr>

</table>
</td></tr>
</table>
</body>
</html>
`)
