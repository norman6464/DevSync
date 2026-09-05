package notify

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sentMessage は fakeSender が受け取った 1 件。
type sentMessage struct {
	userID  uint
	payload []byte
}

type fakeSender struct{ sent []sentMessage }

func (f *fakeSender) SendToUser(userID uint, message []byte) {
	f.sent = append(f.sent, sentMessage{userID: userID, payload: message})
}

// fakeCreator は保存側のスタブ。保存に成功すると ID を採番する。
type fakeCreator struct {
	err     error
	created []*model.Notification
	nextID  uint
}

func (f *fakeCreator) Create(_ context.Context, n *model.Notification) error {
	if f.err != nil {
		return f.err
	}
	f.nextID++
	n.ID = f.nextID
	f.created = append(f.created, n)
	return nil
}

type fakeFollowerNotifier struct {
	err       error
	batch     []*model.Notification
	followers []uint
}

func (f *fakeFollowerNotifier) FindFollowerIDs(context.Context, uint) ([]uint, error) {
	return f.followers, nil
}

func (f *fakeFollowerNotifier) CreateBatch(_ context.Context, notifications []*model.Notification) error {
	if f.err != nil {
		return f.err
	}
	f.batch = notifications
	return nil
}

func decode(t *testing.T, payload []byte) map[string]any {
	t.Helper()
	var got map[string]any
	require.NoError(t, json.Unmarshal(payload, &got))
	return got
}

func TestBroadcastingCreator_Create(t *testing.T) {
	t.Run("保存に成功したら受信者へ配信する", func(t *testing.T) {
		inner := &fakeCreator{}
		sender := &fakeSender{}
		postID := uint(7)
		notification := &model.Notification{
			UserID:    3,
			Type:      model.NotificationTypeFollow,
			ActorID:   9,
			PostID:    &postID,
			CreatedAt: time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC),
		}

		require.NoError(t, NewBroadcastingCreator(inner, sender).Create(context.Background(), notification))

		require.Len(t, sender.sent, 1)
		assert.Equal(t, uint(3), sender.sent[0].userID, "受信者本人にだけ送る")

		got := decode(t, sender.sent[0].payload)
		assert.Equal(t, "notification", got["type"])
		assert.Equal(t, float64(1), got["id"], "保存で採番された ID を載せる")
		assert.Equal(t, string(model.NotificationTypeFollow), got["notification_type"])
		assert.Equal(t, float64(9), got["actor_id"])
		assert.Equal(t, float64(7), got["post_id"])
	})

	t.Run("保存に失敗したら配信しない", func(t *testing.T) {
		inner := &fakeCreator{err: errors.New("db error")}
		sender := &fakeSender{}

		err := NewBroadcastingCreator(inner, sender).Create(context.Background(), &model.Notification{UserID: 3})

		assert.Error(t, err)
		assert.Empty(t, sender.sent, "保存できていない通知は配信しない")
	})

	t.Run("関連が無ければ省略する", func(t *testing.T) {
		inner := &fakeCreator{}
		sender := &fakeSender{}

		require.NoError(t, NewBroadcastingCreator(inner, sender).Create(
			context.Background(), &model.Notification{UserID: 1, Type: model.NotificationTypeBadge, ActorID: 1},
		))

		got := decode(t, sender.sent[0].payload)
		assert.NotContains(t, got, "post_id")
		assert.NotContains(t, got, "question_id")
	})
}

func TestBroadcastingFollowerNotifier_CreateBatch(t *testing.T) {
	t.Run("受信者ごとに配信する", func(t *testing.T) {
		inner := &fakeFollowerNotifier{}
		sender := &fakeSender{}
		notifications := []*model.Notification{
			{ID: 1, UserID: 10, Type: model.NotificationTypePost, ActorID: 5},
			{ID: 2, UserID: 11, Type: model.NotificationTypePost, ActorID: 5},
		}

		require.NoError(t, NewBroadcastingFollowerNotifier(inner, sender).CreateBatch(context.Background(), notifications))

		require.Len(t, sender.sent, 2)
		assert.Equal(t, uint(10), sender.sent[0].userID)
		assert.Equal(t, uint(11), sender.sent[1].userID)
	})

	t.Run("保存に失敗したら配信しない", func(t *testing.T) {
		inner := &fakeFollowerNotifier{err: errors.New("db error")}
		sender := &fakeSender{}

		err := NewBroadcastingFollowerNotifier(inner, sender).CreateBatch(
			context.Background(), []*model.Notification{{UserID: 10}},
		)

		assert.Error(t, err)
		assert.Empty(t, sender.sent)
	})

	t.Run("空なら何もしない", func(t *testing.T) {
		inner := &fakeFollowerNotifier{}
		sender := &fakeSender{}

		require.NoError(t, NewBroadcastingFollowerNotifier(inner, sender).CreateBatch(context.Background(), nil))

		assert.Empty(t, sender.sent)
	})

	t.Run("フォロワー取得はそのまま委ねる", func(t *testing.T) {
		inner := &fakeFollowerNotifier{followers: []uint{1, 2, 3}}

		got, err := NewBroadcastingFollowerNotifier(inner, &fakeSender{}).FindFollowerIDs(context.Background(), 9)

		require.NoError(t, err)
		assert.Equal(t, []uint{1, 2, 3}, got)
	})
}
