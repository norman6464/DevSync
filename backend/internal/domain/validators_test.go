package domain

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"有効なメール", "user@example.com", false},
		{"有効なメール（数字付き）", "user123@example.com", false},
		{"有効なメール（ドット付き）", "user.name@example.co.jp", false},
		{"無効なメール（@なし）", "userexample.com", true},
		{"無効なメール（ドメインなし）", "user@", true},
		{"無効なメール（空）", "", true},
		{"無効なメール（スペースのみ）", "   ", true},
		{"無効なメール（長すぎる）", strings.Repeat("a", 250) + "@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"有効なパスワード", "password123", false},
		{"有効なパスワード（記号付き）", "pass@word!", false},
		{"有効なパスワード（最小長）", "passw0rd!", false},
		{"無効なパスワード（短すぎる）", "pass1!", true},
		{"無効なパスワード（数字記号なし）", "password", true},
		{"無効なパスワード（空）", "", true},
		{"無効なパスワード（長すぎる）", strings.Repeat("a", 130), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有効なユーザー名", "user_name", false},
		{"有効なユーザー名（ハイフン）", "user-name", false},
		{"有効なユーザー名（数字）", "user123", false},
		{"無効なユーザー名（短すぎる）", "u", true},
		{"無効なユーザー名（記号）", "user@name", true},
		{"無効なユーザー名（スペース）", "user name", true},
		{"無効なユーザー名（空）", "", true},
		{"無効なユーザー名（長すぎる）", strings.Repeat("a", 35), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"有効なタイトル", "My Title", false},
		{"有効なタイトル（長い）", strings.Repeat("a", 200), false},
		{"有効なタイトル（マルチバイト200文字）", strings.Repeat("あ", 200), false},
		{"無効なタイトル（空）", "", true},
		{"無効なタイトル（スペースのみ）", "   ", true},
		{"無効なタイトル（長すぎる）", strings.Repeat("a", 201), true},
		{"無効なタイトル（マルチバイト201文字）", strings.Repeat("あ", 201), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"有効なコンテンツ", "Some content", false},
		{"有効なコンテンツ（長い）", strings.Repeat("a", 10000), false},
		{"有効なコンテンツ（マルチバイト10000文字）", strings.Repeat("あ", 10000), false},
		{"無効なコンテンツ（空）", "", true},
		{"無効なコンテンツ（スペースのみ）", "   ", true},
		{"無効なコンテンツ（長すぎる）", strings.Repeat("a", 10001), true},
		{"無効なコンテンツ（マルチバイト10001文字）", strings.Repeat("あ", 10001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContent(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"有効なURL（https）", "https://example.com", false},
		{"有効なURL（http）", "http://example.com/path", false},
		{"有効なURL（クエリ付き）", "https://example.com/search?q=test", false},
		{"有効なURL（空・オプショナル）", "", false},
		{"無効なURL（httpなし）", "example.com", true},
		{"無効なURL（ftpスキーム）", "ftp://example.com", true},
		{"無効なURL（javascriptスキーム）", "javascript:alert(1)", true},
		{"無効なURL（ホスト名なし）", "https://", true},
		{"無効なURL（長すぎる）", "https://" + strings.Repeat("a", 2050), true},
		{"SSRF: localhost", "https://localhost/admin", true},
		{"SSRF: 127.0.0.1", "http://127.0.0.1:8080", true},
		{"SSRF: プライベートIP(10.x)", "http://10.0.0.1", true},
		{"SSRF: プライベートIP(192.168.x)", "http://192.168.1.1", true},
		{"SSRF: プライベートIP(172.16.x)", "http://172.16.0.1", true},
		{"SSRF: 0.0.0.0", "http://0.0.0.0", true},
		{"SSRF: IPv6ループバック", "http://[::1]/path", true},
		{"SSRF: IPv6 unspecified", "http://[::]/path", true},
		{"SSRF: AWSメタデータIP", "http://169.254.169.254/latest/meta-data/", true},
		{"SSRF: GCPメタデータ", "http://metadata.google.internal/computeMetadata/v1/", true},
		{"SSRF: GCPメタデータ（末尾ドット）", "http://metadata.google.internal./computeMetadata/v1/", true},
		{"SSRF: EC2 IPv6メタデータ", "http://[fd00:ec2::254]/latest/meta-data/", true},
		{"SSRF: instance-data", "http://instance-data/latest/meta-data/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsBlockedHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		blocked bool
	}{
		{"localhost", "localhost", true},
		{"0.0.0.0", "0.0.0.0", true},
		{"[::1]", "[::1]", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"プライベートIP(10.x)", "10.0.0.1", true},
		{"プライベートIP(192.168.x)", "192.168.1.1", true},
		{"プライベートIP(172.16.x)", "172.16.0.1", true},
		{"リンクローカル", "169.254.169.254", true},
		{"AWSメタデータIP", "169.254.169.254", true},
		{"GCPメタデータホスト", "metadata.google.internal", true},
		{"GCPメタデータホスト（末尾ドット）", "metadata.google.internal.", true},
		{"EC2 IPv6メタデータ", "fd00:ec2::254", true},
		{"instance-data", "instance-data", true},
		{"大文字小文字混在（LOCALHOST）", "LOCALHOST", true},
		{"大文字小文字混在（Metadata）", "Metadata.Google.Internal", true},
		{"外部ホスト", "example.com", false},
		{"外部IP", "8.8.8.8", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBlockedHost(tt.host)
			assert.Equal(t, tt.blocked, result)
		})
	}
}

func TestIsBlockedIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		blocked bool
	}{
		{"ループバック(127.0.0.1)", "127.0.0.1", true},
		{"ループバック(::1)", "::1", true},
		{"プライベート(10.x)", "10.0.0.1", true},
		{"プライベート(192.168.x)", "192.168.1.1", true},
		{"プライベート(172.16.x)", "172.16.0.1", true},
		{"リンクローカルユニキャスト", "169.254.1.1", true},
		{"未指定(0.0.0.0)", "0.0.0.0", true},
		{"未指定(::)", "::", true},
		{"外部IP(8.8.8.8)", "8.8.8.8", false},
		{"外部IP(1.1.1.1)", "1.1.1.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			assert.NotNil(t, ip, "IPのパースに失敗: %s", tt.ip)
			result := isBlockedIP(ip)
			assert.Equal(t, tt.blocked, result)
		})
	}
}

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		wantErr bool
	}{
		{"有効なタグ", []string{"tag1", "tag2"}, false},
		{"有効なタグ（空配列）", []string{}, false},
		{"無効なタグ（多すぎる）", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}, true},
		{"無効なタグ（空文字）", []string{"tag1", ""}, true},
		{"無効なタグ（長すぎる）", []string{strings.Repeat("a", 31)}, true},
		{"有効なタグ（マルチバイト30文字）", []string{strings.Repeat("あ", 30)}, false},
		{"無効なタグ（マルチバイト31文字）", []string{strings.Repeat("あ", 31)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
		wantErr    bool
	}{
		{"有効なページネーション", 10, 0, 10, 0, false},
		{"limit最大値", 100, 0, 100, 0, false},
		{"limit超過（正規化）", 200, 0, 100, 0, false},
		{"limit未指定（デフォルト）", 0, 0, 10, 0, false},
		{"offset有効", 10, 20, 10, 20, false},
		{"offset負数", 10, -1, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, offset, err := ValidatePagination(tt.limit, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLimit, limit)
				assert.Equal(t, tt.wantOffset, offset)
			}
		})
	}
}

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		min       int
		max       int
		fieldName string
		wantErr   bool
	}{
		{"有効な文字列", "test", 1, 10, "テスト", false},
		{"最小長ちょうど", "a", 1, 10, "テスト", false},
		{"最大長ちょうど", strings.Repeat("a", 10), 1, 10, "テスト", false},
		{"マルチバイト最大長ちょうど", strings.Repeat("あ", 10), 1, 10, "テスト", false},
		{"無効（短すぎる）", "", 1, 10, "テスト", true},
		{"無効（長すぎる）", strings.Repeat("a", 11), 1, 10, "テスト", true},
		{"無効（マルチバイト長すぎる）", strings.Repeat("あ", 11), 1, 10, "テスト", true},
		{"最大長制限なし", strings.Repeat("a", 1000), 1, 0, "テスト", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringLength(tt.s, tt.min, tt.max, tt.fieldName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateExternalUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有効（英数字）", "testuser123", false},
		{"有効（ハイフン）", "test-user", false},
		{"有効（アンダースコア）", "test_user", false},
		{"有効（ドット）", "test.user", false},
		{"有効（1文字）", "a", false},
		{"有効（最大長）", strings.Repeat("a", 50), false},
		{"無効（空）", "", true},
		{"無効（スペースのみ）", "   ", true},
		{"無効（長すぎる）", strings.Repeat("a", 51), true},
		{"無効（スラッシュ）", "user/path", true},
		{"無効（アンパサンド）", "user&param=1", true},
		{"無効（イコール）", "user=value", true},
		{"無効（クエスチョン）", "user?q=1", true},
		{"無効（スペース）", "user name", true},
		{"無効（日本語）", "ユーザー", true},
		{"無効（パーセントエンコード）", "user%2F", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExternalUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLanguageCode(t *testing.T) {
	tests := []struct {
		name    string
		lang    string
		wantErr bool
	}{
		{"有効（日本語）", "ja", false},
		{"有効（英語）", "en", false},
		{"有効（韓国語）", "ko", false},
		{"有効（中国語）", "zh", false},
		{"有効（スペイン語）", "es", false},
		{"有効（フランス語）", "fr", false},
		{"無効（空文字）", "", true},
		{"無効（不正なコード）", "xx", true},
		{"無効（長い文字列）", "japanese", true},
		{"無効（インジェクション試行）", "ja&key=val", true},
		{"無効（SQLインジェクション試行）", "ja' OR '1'='1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLanguageCode(tt.lang)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEnum(t *testing.T) {
	allowedValues := []string{"option1", "option2", "option3"}

	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{"有効な値", "option1", "オプション", false},
		{"有効な値2", "option3", "オプション", false},
		{"空文字（オプショナル）", "", "オプション", false},
		{"無効な値", "invalid", "オプション", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnum(tt.value, allowedValues, tt.fieldName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
