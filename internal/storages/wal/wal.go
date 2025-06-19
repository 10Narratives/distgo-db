package walstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
)

type Storage struct {
	mu   sync.Mutex
	file *os.File
}

var (
	_ walsrv.WALStorage = &Storage{}
)

func New(path string) (*Storage, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Storage{
		file: file,
	}, nil
}

func (s *Storage) Append(ctx context.Context, entry walmodels.WALEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if _, err := s.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write entry to file: %w", err)
	}

	return nil
}

func (s *Storage) Entries(ctx context.Context, size int32, token string, from time.Time, to time.Time) ([]walmodels.WALEntry, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.file.Seek(0, 0); err != nil {
		return nil, "", fmt.Errorf("failed to seek file: %w", err)
	}

	var entries []walmodels.WALEntry
	var nextToken string
	var lineCount int32 = 0

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var entry walmodels.WALEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue // Skip invalid entries
		}

		if !entry.Timestamp.After(from) || !entry.Timestamp.Before(to) {
			continue
		}

		if token != "" && nextToken == "" {
			nextToken = entry.ID.String()
			continue
		}

		entries = append(entries, entry)
		lineCount++

		if size > 0 && lineCount >= size {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(entries) > 0 {
		lastEntry := entries[len(entries)-1]
		nextToken = lastEntry.ID.String()
	}

	return entries, nextToken, nil
}
func (s *Storage) Truncate(ctx context.Context, before time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tempFile, err := os.CreateTemp("", "wal-temp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := s.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var entry walmodels.WALEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}

		if entry.Timestamp.After(before) || entry.Timestamp.Equal(before) {
			if _, err := tempFile.Write(append(line, '\n')); err != nil {
				return fmt.Errorf("failed to write to temp file: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	if err := os.Rename(tempFile.Name(), s.file.Name()); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}

	newFile, err := os.OpenFile(s.file.Name(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to reopen file: %w", err)
	}
	s.file = newFile

	return nil
}
