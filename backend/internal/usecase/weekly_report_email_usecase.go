package usecase

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// defaultEmailLanguage はユーザーが言語を設定していないときに使う既定の言語。
const defaultEmailLanguage = "ja"

// weeklyReportTmpl はウィークリーレポートメールのテンプレート。パースは 1 度だけ行う。
var weeklyReportTmpl = template.Must(template.New("weekly_report").Parse(weeklyReportTemplate))

// SendWeeklyReportUseCase は 1 ユーザーへウィークリーレポートメールを送信する。
type SendWeeklyReportUseCase struct {
	sender repository.EmailSender
	appURL string
}

// NewSendWeeklyReportUseCase は SendWeeklyReportUseCase を生成する。
func NewSendWeeklyReportUseCase(sender repository.EmailSender, appURL string) *SendWeeklyReportUseCase {
	return &SendWeeklyReportUseCase{sender: sender, appURL: appURL}
}

// Execute はレポートを HTML に整形してメールを送信する。
func (uc *SendWeeklyReportUseCase) Execute(ctx context.Context, user *model.User, report *model.ActivityReport) error {
	if user.Email == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "メールアドレスが空です", nil)
	}
	if report == nil {
		return domain.NewError(domain.ErrCodeBadRequest, "レポートがありません", nil)
	}

	lang := user.EmailLanguage
	if lang == "" {
		lang = defaultEmailLanguage
	}

	html, err := RenderWeeklyReportHTML(user, report, lang, uc.appURL)
	if err != nil {
		return domain.NewError(domain.ErrCodeInternal, "メールテンプレートのレンダリングに失敗", err)
	}

	subject := fmt.Sprintf("[DevSync] %s", getEmailTexts(lang)["subject"])
	return uc.sender.Send(ctx, user.Email, subject, html)
}

// SendAllWeeklyReportsUseCase は配信対象の全ユーザーへウィークリーレポートメールを送信する。
type SendAllWeeklyReportsUseCase struct {
	users   repository.UserRepository
	reports repository.WeeklyActivityReportReader
	send    *SendWeeklyReportUseCase
}

// NewSendAllWeeklyReportsUseCase は SendAllWeeklyReportsUseCase を生成する。
func NewSendAllWeeklyReportsUseCase(
	users repository.UserRepository,
	reports repository.WeeklyActivityReportReader,
	send *SendWeeklyReportUseCase,
) *SendAllWeeklyReportsUseCase {
	return &SendAllWeeklyReportsUseCase{users: users, reports: reports, send: send}
}

// Execute は配信を有効にしている全ユーザーへ送信する。
// 1 ユーザーの失敗で全体を止めず、ログに残して次のユーザーへ進む。
func (uc *SendAllWeeklyReportsUseCase) Execute(ctx context.Context) error {
	users, err := uc.users.FindAll(ctx)
	if err != nil {
		return domain.NewError(domain.ErrCodeDatabase, "ユーザー一覧の取得に失敗", err)
	}

	for i := range users {
		user := &users[i]
		if !user.EmailWeeklyReport {
			continue
		}

		report, err := uc.reports.GetWeeklyReport(ctx, user.ID)
		if err != nil {
			log.Printf("ウィークリーレポート生成失敗 (userID=%d): %v", user.ID, err)
			continue
		}

		if err := uc.send.Execute(ctx, user, report); err != nil {
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

// RenderWeeklyReportHTML はレポートデータから HTML メール本文をレンダリングする。
// テンプレートと言語テキストだけに依存する純粋な処理。
func RenderWeeklyReportHTML(user *model.User, report *model.ActivityReport, lang, appURL string) (string, error) {
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
		AppURL:             appURL,
		Texts:              texts,
	}

	var buf bytes.Buffer
	if err := weeklyReportTmpl.Execute(&buf, data); err != nil {
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
