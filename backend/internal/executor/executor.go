package executor

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/config"
	"github.com/dojo-harvester/backend/internal/logger"
)

type Process struct {
	ID        string
	Cmd       *exec.Cmd
	Cancel    context.CancelFunc
	StartTime time.Time
}

type Executor struct {
	cfg       *config.ScriptsConfig
	processes map[string]*Process
	mu        sync.RWMutex
}

func NewExecutor(cfg *config.Config) *Executor {
	return &Executor{
		cfg:       &cfg.Scripts,
		processes: make(map[string]*Process),
	}
}

// StartScript executes a given script and tracks its lifecycle
func (e *Executor) StartScript(processID string, scriptPath string, args []string) error {
	e.mu.Lock()
	if _, exists := e.processes[processID]; exists {
		e.mu.Unlock()
		return fmt.Errorf("process %s is already running", processID)
	}
	e.mu.Unlock()

	// Default timeout for any script
	timeout := time.Duration(e.cfg.TimeoutSecs) * time.Second
	if timeout <= 0 {
		timeout = 1 * time.Hour
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	
	cmd := exec.CommandContext(ctx, scriptPath, args...)
	
	// Add process to tracker
	proc := &Process{
		ID:        processID,
		Cmd:       cmd,
		Cancel:    cancel,
		StartTime: time.Now(),
	}

	e.mu.Lock()
	e.processes[processID] = proc
	e.mu.Unlock()

	logger.Get().Info("Starting script", zap.String("id", processID), zap.String("script", scriptPath), zap.Strings("args", args))

	if err := cmd.Start(); err != nil {
		cancel()
		e.RemoveProcess(processID)
		return fmt.Errorf("failed to start script: %w", err)
	}

	// Run wait in background to clean up when finished
	go func() {
		err := cmd.Wait()
		e.RemoveProcess(processID)

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				logger.Get().Error("Script timed out", zap.String("id", processID))
			} else {
				logger.Get().Error("Script failed", zap.String("id", processID), zap.Error(err))
			}
		} else {
			logger.Get().Info("Script completed successfully", zap.String("id", processID))
		}
	}()

	return nil
}

// StopScript cancels the context of a running script to stop it
func (e *Executor) StopScript(processID string) error {
	e.mu.RLock()
	proc, exists := e.processes[processID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("process %s not found", processID)
	}

	logger.Get().Info("Stopping script manually", zap.String("id", processID))
	proc.Cancel()
	
	return nil
}

// RemoveProcess cleans up the tracking map
func (e *Executor) RemoveProcess(processID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.processes, processID)
}

// GetProcessStatus returns if a process is still running
func (e *Executor) GetProcessStatus(processID string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, exists := e.processes[processID]
	return exists
}
