package walstorage_test

// import (
// 	"context"
// 	"encoding/json"
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// 	walstorage "github.com/10Narratives/distgo-db/internal/storages/wal"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestStorage_LogEntry(t *testing.T) {
// 	dir := t.TempDir()
// 	filePath := filepath.Join(dir, "wal.log")

// 	storage, err := walstorage.New(filePath)
// 	require.NoError(t, err)

// 	entry := walmodels.WALEntry{
// 		ID:        "entry1",
// 		Target:    "test",
// 		Type:      commonmodels.MutationTypeCreate,
// 		Timestamp: time.Now(),
// 		NewValue:  `{"key":"value"}`,
// 	}

// 	err = storage.LogEntry(context.Background(), entry)
// 	require.NoError(t, err)

// 	// Verify the entry was written to the file
// 	file, err := os.Open(filePath)
// 	require.NoError(t, err)
// 	defer file.Close()

// 	scanner := walstorage.NewLineScanner(file)
// 	var storedEntry walmodels.WALEntry
// 	for scanner.Scan() {
// 		line := scanner.Bytes()
// 		err := json.Unmarshal(line, &storedEntry)
// 		require.NoError(t, err)
// 		break // Only one entry expected
// 	}

// 	assert.Equal(t, entry.ID, storedEntry.ID)
// 	assert.Equal(t, entry.Target, storedEntry.Target)
// 	assert.Equal(t, entry.Type, storedEntry.Type)
// 	assert.Equal(t, entry.NewValue, storedEntry.NewValue)
// 	assert.WithinDuration(t, entry.Timestamp, storedEntry.Timestamp, time.Second)
// }

// func TestStorage_Entries(t *testing.T) {
// 	dir := t.TempDir()
// 	filePath := filepath.Join(dir, "wal.log")

// 	storage, err := walstorage.New(filePath)
// 	require.NoError(t, err)

// 	entries := []walmodels.WALEntry{
// 		{
// 			ID:        "entry1",
// 			Target:    "test",
// 			Type:      commonmodels.MutationTypeCreate,
// 			Timestamp: time.Now().Add(-2 * time.Hour),
// 			NewValue:  `{"key":"value1"}`,
// 		},
// 		{
// 			ID:        "entry2",
// 			Target:    "test",
// 			Type:      commonmodels.MutationTypeUpdate,
// 			Timestamp: time.Now().Add(-1 * time.Hour),
// 			NewValue:  `{"key":"value2"}`,
// 		},
// 		{
// 			ID:        "entry3",
// 			Target:    "test",
// 			Type:      commonmodels.MutationTypeDelete,
// 			Timestamp: time.Now(),
// 			NewValue:  `{"key":"value3"}`,
// 		},
// 	}

// 	for _, entry := range entries {
// 		err := storage.LogEntry(context.Background(), entry)
// 		require.NoError(t, err)
// 	}

// 	t.Run("retrieve all entries", func(t *testing.T) {
// 		retrieved, _, err := storage.Entries(context.Background(), 10, "", time.Time{}, time.Now())
// 		require.NoError(t, err)
// 		assert.Len(t, retrieved, 3)
// 		assert.Equal(t, entries[0].ID, retrieved[0].ID)
// 		assert.Equal(t, entries[1].ID, retrieved[1].ID)
// 		assert.Equal(t, entries[2].ID, retrieved[2].ID)
// 	})

// 	t.Run("retrieve with size limit", func(t *testing.T) {
// 		retrieved, _, err := storage.Entries(context.Background(), 2, "", time.Time{}, time.Now())
// 		require.NoError(t, err)
// 		assert.Len(t, retrieved, 2)
// 		assert.Equal(t, entries[0].ID, retrieved[0].ID)
// 		assert.Equal(t, entries[1].ID, retrieved[1].ID)
// 	})

// 	t.Run("retrieve with timestamp filter", func(t *testing.T) {
// 		from := time.Now().Add(-90 * time.Minute)
// 		to := time.Now()
// 		retrieved, _, err := storage.Entries(context.Background(), 10, "", from, to)
// 		require.NoError(t, err)
// 		assert.Len(t, retrieved, 2)
// 		assert.Equal(t, entries[1].ID, retrieved[0].ID)
// 		assert.Equal(t, entries[2].ID, retrieved[1].ID)
// 	})

// 	t.Run("retrieve with token", func(t *testing.T) {
// 		retrieved, _, err := storage.Entries(context.Background(), 10, "entry2", time.Time{}, time.Now())
// 		require.NoError(t, err)
// 		assert.Len(t, retrieved, 1)
// 		assert.Equal(t, entries[2].ID, retrieved[0].ID)
// 	})
// }

// // func TestStorage_Truncate(t *testing.T) {
// // 	dir := t.TempDir()
// // 	filePath := filepath.Join(dir, "wal.log")

// // 	storage, err := walstorage.New(filePath)
// // 	require.NoError(t, err)

// // 	entries := []walmodels.WALEntry{
// // 		{
// // 			ID:        "entry1",
// // 			Target:    "test",
// // 			Type:      commonmodels.MutationTypeCreate,
// // 			Timestamp: time.Now().Add(-2 * time.Hour),
// // 			NewValue:  `{"key":"value1"}`,
// // 		},
// // 		{
// // 			ID:        "entry2",
// // 			Target:    "test",
// // 			Type:      commonmodels.MutationTypeUpdate,
// // 			Timestamp: time.Now().Add(-1 * time.Hour),
// // 			NewValue:  `{"key":"value2"}`,
// // 		},
// // 		{
// // 			ID:        "entry3",
// // 			Target:    "test",
// // 			Type:      commonmodels.MutationTypeDelete,
// // 			Timestamp: time.Now(),
// // 			NewValue:  `{"key":"value3"}`,
// // 		},
// // 	}

// // 	for _, entry := range entries {
// // 		err := storage.LogEntry(context.Background(), entry)
// // 		require.NoError(t, err)
// // 	}

// // 	before := time.Now().Add(-100 * time.Minute)
// // 	err = storage.Truncate(context.Background(), before)
// // 	require.NoError(t, err)

// // 	retrieved, _, err := storage.Entries(context.Background(), 10, "", time.Time{}, time.Now())
// // 	require.NoError(t, err)

// // 	assert.Len(t, retrieved, 1)
// // 	//assert.Equal(t, entries[2].ID, retrieved[0].ID)
// // }

// // func TestStorage_Errors(t *testing.T) {
// // 	dir := t.TempDir()
// // 	filePath := filepath.Join(dir, "wal.log")

// // 	storage, err := walstorage.New(filePath)
// // 	require.NoError(t, err)

// // 	t.Run("invalid token", func(t *testing.T) {
// // 		_, _, err := storage.Entries(context.Background(), 10, "nonexistent", time.Time{}, time.Now())
// // 		require.Error(t, err)
// // 		assert.Contains(t, err.Error(), "invalid token")
// // 	})

// // 	t.Run("file read error", func(t *testing.T) {
// // 		// Simulate a corrupted WAL file
// // 		err := os.WriteFile(filePath, []byte("corrupted"), 0644)
// // 		require.NoError(t, err)

// // 		_, _, err = storage.Entries(context.Background(), 10, "", time.Time{}, time.Now())
// // 		require.Error(t, err)
// // 		assert.Contains(t, err.Error(), "failed to deserialize WAL entry")
// // 	})

// // 	t.Run("truncate error", func(t *testing.T) {
// // 		// Make the file read-only
// // 		err := os.Chmod(filePath, 0444)
// // 		require.NoError(t, err)

// // 		err = storage.Truncate(context.Background(), time.Now())
// // 		require.Error(t, err)
// // 		assert.Contains(t, err.Error(), "failed to create new WAL file during truncation")
// // 	})
// // }
