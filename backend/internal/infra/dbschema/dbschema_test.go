package dbschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsExpectedUserGitHubIndex は pg_indexes が返す定義の突き合わせを確認する。
// 同名で定義の違うインデックスがあると CREATE ... IF NOT EXISTS は中身を確かめずに
// 成功するため、ここでの判定が「気付かないまま一意性が失われる」ことを防ぐ。
func TestIsExpectedUserGitHubIndex(t *testing.T) {
	tests := []struct {
		name string
		def  string
		want bool
	}{
		{
			name: "意図した定義",
			def:  "CREATE UNIQUE INDEX idx_users_git_hub_id_linked ON public.users USING btree (git_hub_id) WHERE (git_hub_id <> 0)",
			want: true,
		},
		{
			name: "ユニークでない",
			def:  "CREATE INDEX idx_users_git_hub_id_linked ON public.users USING btree (git_hub_id) WHERE (git_hub_id <> 0)",
			want: false,
		},
		{
			name: "述語が無い（全行が対象で 2 人目の登録が衝突する）",
			def:  "CREATE UNIQUE INDEX idx_users_git_hub_id_linked ON public.users USING btree (git_hub_id)",
			want: false,
		},
		{
			name: "述語が違う",
			def:  "CREATE UNIQUE INDEX idx_users_git_hub_id_linked ON public.users USING btree (git_hub_id) WHERE (git_hub_id > 0)",
			want: false,
		},
		{
			name: "対象の列が違う",
			def:  "CREATE UNIQUE INDEX idx_users_git_hub_id_linked ON public.users USING btree (email) WHERE (git_hub_id <> 0)",
			want: false,
		},
		{
			name: "空（インデックスが無い）",
			def:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isExpectedUserGitHubIndex(tt.def))
		})
	}
}
