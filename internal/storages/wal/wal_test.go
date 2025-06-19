package walstorage_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	walstorage "github.com/10Narratives/distgo-db/internal/storages/wal"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_Append(t *testing.T) {
	t.Parallel()

	tempFile, err := os.CreateTemp("", "wal-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storage, err := walstorage.New(tempFile.Name())
	require.NoError(t, err)

	entry := walmodels.WALEntry{
		ID:        uuid.New(),
		Timestamp: time.Now().UTC(),
		Mutation:  commonmodels.MutationTypeCreate,
		Payload:   json.RawMessage(`{"key":"value"}`),
		Entity:    1,
	}

	err = storage.Append(context.Background(), entry)
	require.NoError(t, err)

	entries, _, err := storage.Entries(context.Background(), 10, "", time.Time{}, time.Now().Add(time.Second))
	require.NoError(t, err)
	require.Len(t, entries, 1)

	assert.Equal(t, entry.ID, entries[0].ID)
	assert.Equal(t, entry.Timestamp.Unix(), entries[0].Timestamp.Unix())
	assert.Equal(t, entry.Mutation, entries[0].Mutation)
	assert.JSONEq(t, string(entry.Payload), string(entries[0].Payload))
	assert.Equal(t, entry.Entity, entries[0].Entity)
}

func TestStorage_Entries(t *testing.T) {
	t.Parallel()

	tempFile, err := os.CreateTemp("", "wal-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storage, err := walstorage.New(tempFile.Name())
	require.NoError(t, err)

	now := time.Now().UTC()
	entries := []walmodels.WALEntry{
		{
			ID:        uuid.New(),
			Timestamp: now.Add(-2 * time.Hour),
			Mutation:  commonmodels.MutationTypeCreate,
			Payload:   json.RawMessage(`{"key":"value1"}`),
			Entity:    1,
		},
		{
			ID:        uuid.New(),
			Timestamp: now.Add(-1 * time.Hour),
			Mutation:  commonmodels.MutationTypeUpdate,
			Payload:   json.RawMessage(`{"key":"value2"}`),
			Entity:    1,
		},
		{
			ID:        uuid.New(),
			Timestamp: now,
			Mutation:  commonmodels.MutationTypeDelete,
			Payload:   json.RawMessage(`{"key":"value3"}`),
			Entity:    1,
		},
	}

	for _, entry := range entries {
		err := storage.Append(context.Background(), entry)
		require.NoError(t, err)
	}

	from := now.Add(-3 * time.Hour)
	to := now.Add(-1*time.Hour + time.Minute)
	result, _, err := storage.Entries(context.Background(), 10, "", from, to)
	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, entries[0].ID, result[0].ID)
	assert.Equal(t, entries[1].ID, result[1].ID)

	result, nextToken, err := storage.Entries(context.Background(), 2, "", time.Time{}, time.Now().Add(time.Second))
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.NotEmpty(t, nextToken)

	result, _, err = storage.Entries(context.Background(), 2, nextToken, time.Time{}, time.Now().Add(time.Second))
	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestStorage_Truncate(t *testing.T) {
	t.Parallel()

	tempFile, err := os.CreateTemp("", "wal-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storage, err := walstorage.New(tempFile.Name())
	require.NoError(t, err)

	now := time.Now().UTC()
	entries := []walmodels.WALEntry{
		{
			ID:        uuid.New(),
			Timestamp: now.Add(-3 * time.Hour),
			Mutation:  commonmodels.MutationTypeCreate,
			Payload:   json.RawMessage(`{"key":"value1"}`),
			Entity:    1,
		},
		{
			ID:        uuid.New(),
			Timestamp: now.Add(-2 * time.Hour),
			Mutation:  commonmodels.MutationTypeUpdate,
			Payload:   json.RawMessage(`{"key":"value2"}`),
			Entity:    1,
		},
		{
			ID:        uuid.New(),
			Timestamp: now,
			Mutation:  commonmodels.MutationTypeDelete,
			Payload:   json.RawMessage(`{"key":"value3"}`),
			Entity:    1,
		},
	}

	for _, entry := range entries {
		err := storage.Append(context.Background(), entry)
		require.NoError(t, err)
	}

	before := now.Add(-2 * time.Hour)
	err = storage.Truncate(context.Background(), before)
	require.NoError(t, err)

	result, _, err := storage.Entries(context.Background(), 10, "", time.Time{}, time.Now().Add(time.Second))
	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, entries[1].ID, result[0].ID)
	assert.Equal(t, entries[2].ID, result[1].ID)
}

func TestStorage_Errors(t *testing.T) {
	t.Parallel()

	_, err := walstorage.New("")
	assert.Error(t, err)

	tempFile, err := os.CreateTemp("", "wal-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storage, err := walstorage.New(tempFile.Name())
	require.NoError(t, err)

	_, _, err = storage.Entries(context.Background(), 10, "", time.Time{}, time.Now().Add(time.Second))
	assert.NoError(t, err)

	err = storage.Truncate(context.Background(), time.Time{})
	assert.NoError(t, err)
}
