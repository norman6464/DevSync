package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testEntity はcheckOwnershipテスト用のダミーエンティティ。
type testEntity struct {
	ID      uint
	OwnerID uint
}

func TestCheckOwnership_Success(t *testing.T) {
	entity := &testEntity{ID: 1, OwnerID: 10}
	finder := func(id uint) (*testEntity, error) {
		return entity, nil
	}
	getOwnerID := func(e *testEntity) uint { return e.OwnerID }

	result, err := checkOwnership(finder, 1, 10, getOwnerID)
	assert.NoError(t, err)
	assert.Equal(t, entity, result)
}

func TestCheckOwnership_Forbidden(t *testing.T) {
	entity := &testEntity{ID: 1, OwnerID: 10}
	finder := func(id uint) (*testEntity, error) {
		return entity, nil
	}
	getOwnerID := func(e *testEntity) uint { return e.OwnerID }

	result, err := checkOwnership(finder, 1, 99, getOwnerID)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCheckOwnership_FinderError(t *testing.T) {
	dbErr := errors.New("record not found")
	finder := func(id uint) (*testEntity, error) {
		return nil, dbErr
	}
	getOwnerID := func(e *testEntity) uint { return e.OwnerID }

	result, err := checkOwnership(finder, 1, 10, getOwnerID)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
}

func TestCheckOwnership_ZeroUserID(t *testing.T) {
	entity := &testEntity{ID: 1, OwnerID: 10}
	finder := func(id uint) (*testEntity, error) {
		return entity, nil
	}
	getOwnerID := func(e *testEntity) uint { return e.OwnerID }

	result, err := checkOwnership(finder, 1, 0, getOwnerID)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCheckOwnership_OwnerIDZero_MatchesZeroUser(t *testing.T) {
	entity := &testEntity{ID: 1, OwnerID: 0}
	finder := func(id uint) (*testEntity, error) {
		return entity, nil
	}
	getOwnerID := func(e *testEntity) uint { return e.OwnerID }

	result, err := checkOwnership(finder, 1, 0, getOwnerID)
	assert.NoError(t, err)
	assert.Equal(t, entity, result)
}

func TestCheckOwnership_FinderPassesCorrectID(t *testing.T) {
	var receivedID uint
	finder := func(id uint) (*testEntity, error) {
		receivedID = id
		return &testEntity{ID: id, OwnerID: 10}, nil
	}
	getOwnerID := func(e *testEntity) uint { return e.OwnerID }

	_, _ = checkOwnership(finder, 42, 10, getOwnerID)
	assert.Equal(t, uint(42), receivedID)
}
