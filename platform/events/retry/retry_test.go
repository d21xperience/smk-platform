package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPolicyDelayForAttempt(t *testing.T) {
	policy := Policy{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		BackoffType:  BackoffExponential,
		Jitter:       false,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 100 * time.Millisecond}, // 100ms * 2^0
		{1, 200 * time.Millisecond}, // 100ms * 2^1
		{2, 400 * time.Millisecond}, // 100ms * 2^2
		{3, 800 * time.Millisecond}, // 100ms * 2^3
	}

	for _, tt := range tests {
		got := policy.DelayForAttempt(tt.attempt)
		if got != tt.expected {
			t.Errorf("Attempt %d: expected %v, got %v", tt.attempt, tt.expected, got)
		}
	}
}

func TestPolicyMaxDelayCap(t *testing.T) {
	policy := Policy{
		MaxRetries:   10,
		InitialDelay: 1 * time.Second,
		MaxDelay:     5 * time.Second,
		BackoffType:  BackoffExponential,
		Jitter:       false,
	}

	// Attempt 10 akan menghasilkan 1024s, tapi di-cap ke 5s
	delay := policy.DelayForAttempt(10)
	if delay != 5*time.Second {
		t.Errorf("Expected delay capped at 5s, got %v", delay)
	}
}

func TestPolicyLinearBackoff(t *testing.T) {
	policy := Policy{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		BackoffType:  BackoffLinear,
		Jitter:       false,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 100 * time.Millisecond}, // 100ms * 1
		{1, 200 * time.Millisecond}, // 100ms * 2
		{2, 300 * time.Millisecond}, // 100ms * 3
	}

	for _, tt := range tests {
		got := policy.DelayForAttempt(tt.attempt)
		if got != tt.expected {
			t.Errorf("Attempt %d: expected %v, got %v", tt.attempt, tt.expected, got)
		}
	}
}

func TestExecutorSuccess(t *testing.T) {
	executor := NewExecutor(Policy{
		MaxRetries:   3,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		BackoffType:  BackoffConstant,
	})

	var attempts int
	result := executor.Execute(context.Background(), func() error {
		attempts++
		return nil
	})

	if !result.Success {
		t.Error("Expected success")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
	if result.Attempts != 1 {
		t.Errorf("Expected result.Attempts=1, got %d", result.Attempts)
	}
}

func TestExecutorRetryThenSuccess(t *testing.T) {
	executor := NewExecutor(Policy{
		MaxRetries:   3,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		BackoffType:  BackoffConstant,
	})

	var attempts int
	result := executor.Execute(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if !result.Success {
		t.Error("Expected success after retry")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestExecutorMaxRetriesExceeded(t *testing.T) {
	executor := NewExecutor(Policy{
		MaxRetries:   2,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		BackoffType:  BackoffConstant,
	})

	var attempts int
	result := executor.Execute(context.Background(), func() error {
		attempts++
		return errors.New("permanent error")
	})

	if result.Success {
		t.Error("Expected failure")
	}
	if attempts != 3 { // 1 initial + 2 retries
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
	if result.LastError == nil {
		t.Error("Expected last error to be set")
	}
}

func TestExecutorContextCancellation(t *testing.T) {
	executor := NewExecutor(Policy{
		MaxRetries:   5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		BackoffType:  BackoffConstant,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result := executor.Execute(ctx, func() error {
		return errors.New("error")
	})

	if result.Success {
		t.Error("Expected failure due to context cancellation")
	}
}

func TestExecutorWithCallback(t *testing.T) {
	executor := NewExecutor(Policy{
		MaxRetries:   2,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		BackoffType:  BackoffConstant,
	})

	var callbackCalls int
	var attempts int

	result := executor.ExecuteWithCallback(
		context.Background(),
		func() error {
			attempts++
			if attempts <= 2 {
				return errors.New("error")
			}
			return nil
		},
		func(attempt int, err error, delay time.Duration) {
			callbackCalls++
		},
	)

	if !result.Success {
		t.Error("Expected success")
	}
	if callbackCalls != 2 {
		t.Errorf("Expected 2 callback calls, got %d", callbackCalls)
	}
}

func TestPolicyValidate(t *testing.T) {
	tests := []struct {
		name    string
		policy  Policy
		wantErr bool
	}{
		{"valid", Policy{MaxRetries: 3, InitialDelay: time.Second, MaxDelay: 30 * time.Second}, false},
		{"negative max retries", Policy{MaxRetries: -1}, true},
		{"negative initial delay", Policy{MaxRetries: 3, InitialDelay: -time.Second}, true},
		{"max delay < initial", Policy{MaxRetries: 3, InitialDelay: 10 * time.Second, MaxDelay: time.Second}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
