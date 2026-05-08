package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var ErrAlreadyRunning = errors.New("another instance is already running")

type Lock struct {
	f *os.File
}

func Acquire() (*Lock, error) {
	path := filepath.Join(os.TempDir(), "ai-check-guardrails.lock")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("opening lock file: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, ErrAlreadyRunning
		}
		return nil, fmt.Errorf("acquiring lock: %w", err)
	}
	return &Lock{f: f}, nil
}

func (l *Lock) Release() {
	if l.f != nil {
		syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
		l.f.Close()
	}
}
