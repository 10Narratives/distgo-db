package walstore

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"

	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/wal"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
)

type Storage struct {
	file *os.File
	mu   sync.Mutex
}

var _ documentsrv.WALStorage = &Storage{}

func New(path string) (*Storage, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot open WAL file: %w", err)
	}
	return &Storage{file: file}, nil
}

func (s *Storage) Write(entry walmodels.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.file.Write(append(entry, '\n'))
	return err
}

func (s *Storage) Replay(handler func(walmodels.Entry) error) error {
	data, err := os.ReadFile(s.file.Name())
	if err != nil {
		return fmt.Errorf("failed to read WAL file: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if err := handler(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error during WAL replay: %w", err)
	}

	return nil
}
