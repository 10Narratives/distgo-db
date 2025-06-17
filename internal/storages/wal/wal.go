package walstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	commonsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/common"
	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
)

const (
	fileMode = 0644
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Storage struct {
	mu       sync.Mutex
	filePath string
}

var (
	_ walsrv.WALStorage    = &Storage{}
	_ commonsrv.WALStorage = &Storage{}
)

// New creates a new file-based WAL storage.
func New(filePath string) (*Storage, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, fileMode)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}
	file.Close()

	return &Storage{
		filePath: filePath,
	}, nil
}

func (s *Storage) LogEntry(ctx context.Context, entry walmodels.WALEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to serialize WAL entry: %w", err)
	}

	file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY, fileMode)
	if err != nil {
		return fmt.Errorf("failed to open WAL file for writing: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write WAL entry: %w", err)
	}

	return nil
}

// Entries retrieves WAL entries based on the provided parameters.
func (s *Storage) Entries(
	ctx context.Context,
	size int32,
	token string,
	from time.Time,
	to time.Time,
) ([]walmodels.WALEntry, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open WAL file for reading: %w", err)
	}
	defer file.Close()

	var entries []walmodels.WALEntry
	var lastToken string
	scanner := NewLineScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		var entry walmodels.WALEntry

		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, "", fmt.Errorf("failed to deserialize WAL entry: %w", err)
		}

		if !entry.Timestamp.After(from) || !entry.Timestamp.Before(to) {
			continue
		}

		if token != "" {
			if entry.ID != token {
				continue
			}
			token = ""
			continue
		}

		entries = append(entries, entry)
		lastToken = entry.ID

		if int32(len(entries)) >= size {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("error reading WAL file: %w", err)
	}

	return entries, lastToken, nil
}

func (s *Storage) Truncate(ctx context.Context, before time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to open WAL file for truncation: %w", err)
	}
	defer file.Close()

	var entries []walmodels.WALEntry
	scanner := NewLineScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		var entry walmodels.WALEntry

		if err := json.Unmarshal(line, &entry); err != nil {
			return fmt.Errorf("failed to deserialize WAL entry during truncation: %w", err)
		}

		// Retain entries with timestamps equal to or later than `before`
		if !entry.Timestamp.Before(before) {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading WAL file during truncation: %w", err)
	}

	// Overwrite the file with the remaining entries
	newFile, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to create new WAL file during truncation: %w", err)
	}
	defer newFile.Close()

	for _, entry := range entries {
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to serialize WAL entry during truncation: %w", err)
		}

		if _, err := newFile.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("failed to write WAL entry during truncation: %w", err)
		}
	}

	return nil
}

func NewLineScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	return scanner
}
