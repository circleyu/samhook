package samhook

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestSendWithOptions_Success(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	msg := createTestMessage()
	err := SendWithOptions(server.URL, msg)
	if err != nil {
		t.Fatalf("SendWithOptions() error = %v", err)
	}
}

func TestSendWithOptions_WithTimeout(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	msg := createTestMessage()
	err := SendWithOptions(server.URL, msg, WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("SendWithOptions() error = %v", err)
	}
}

func TestSendWithOptions_WithClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	msg := createTestMessage()
	err := SendWithOptions(server.URL, msg, WithClient(customClient))
	if err != nil {
		t.Fatalf("SendWithOptions() error = %v", err)
	}
}

func TestSendWithContext_Success(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	ctx := context.Background()
	msg := createTestMessage()
	err := SendWithContext(ctx, server.URL, msg)
	if err != nil {
		t.Fatalf("SendWithContext() error = %v", err)
	}
}

func TestSendWithContext_WithTimeout(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	msg := createTestMessage()
	err := SendWithContext(ctx, server.URL, msg)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestSendWithContext_Cancellation(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	msg := createTestMessage()
	err := SendWithContext(ctx, server.URL, msg)
	if err == nil {
		t.Error("expected cancellation error")
	}
}
