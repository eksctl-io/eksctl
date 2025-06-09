package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionManagement tests session registration and management
func TestSessionManagement(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Create a test session
	sessionID := "test-session"
	notificationChan := make(chan mcp.JSONRPCNotification, 10)
	session := &testSession{
		id:                 sessionID,
		notificationChan:   notificationChan,
		isInitialized:      false,
		initializationFunc: func() { /* no-op */ },
	}

	// Register the session
	err = mcpServer.RegisterSession(context.Background(), session)
	assert.NoError(t, err)

	// Try to register the same session again (should fail)
	err = mcpServer.RegisterSession(context.Background(), session)
	assert.Error(t, err)

	// Initialize the session
	session.Initialize()
	assert.True(t, session.Initialized())

	// Send a notification to the session
	err = mcpServer.SendNotificationToSpecificClient(sessionID, "test-notification", map[string]any{
		"data": "test-data",
	})
	assert.NoError(t, err)

	// Check that the notification was received
	select {
	case notification := <-notificationChan:
		assert.Equal(t, "test-notification", notification.Method)
		assert.Equal(t, "test-data", notification.Params.AdditionalFields["data"])
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected notification not received")
	}

	// Unregister the session
	mcpServer.UnregisterSession(context.Background(), sessionID)

	// Try to send a notification to the unregistered session (should fail)
	err = mcpServer.SendNotificationToSpecificClient(sessionID, "test-notification", nil)
	assert.Error(t, err)
}

// TestSessionContext tests session context management
func TestSessionContext(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Create a test session
	sessionID := "test-session"
	notificationChan := make(chan mcp.JSONRPCNotification, 10)
	session := &testSession{
		id:                 sessionID,
		notificationChan:   notificationChan,
		isInitialized:      true,
		initializationFunc: func() { /* no-op */ },
	}

	// Register the session
	err = mcpServer.RegisterSession(context.Background(), session)
	assert.NoError(t, err)

	// Create a context with the session
	ctx := mcpServer.WithContext(context.Background(), session)

	// Get the session from the context
	retrievedSession := server.ClientSessionFromContext(ctx)
	assert.NotNil(t, retrievedSession)
	assert.Equal(t, sessionID, retrievedSession.SessionID())
}

// TestSessionNotifications tests sending notifications to sessions
func TestSessionNotifications(t *testing.T) {
	// Create a simple root command for testing
	rootCmd := &cobra.Command{
		Use:   "eksctl",
		Short: "The official CLI for Amazon EKS",
	}

	// Create the MCP server
	mcpServer, err := newEksctlMCPServer(rootCmd, "test-version")
	require.NoError(t, err)
	require.NotNil(t, mcpServer)

	// Create multiple test sessions
	session1 := &testSession{
		id:                 "session-1",
		notificationChan:   make(chan mcp.JSONRPCNotification, 10),
		isInitialized:      true,
		initializationFunc: func() { /* no-op */ },
	}
	session2 := &testSession{
		id:                 "session-2",
		notificationChan:   make(chan mcp.JSONRPCNotification, 10),
		isInitialized:      true,
		initializationFunc: func() { /* no-op */ },
	}
	session3 := &testSession{
		id:                 "session-3",
		notificationChan:   make(chan mcp.JSONRPCNotification, 10),
		isInitialized:      false, // Not initialized
		initializationFunc: func() { /* no-op */ },
	}

	// Register the sessions
	err = mcpServer.RegisterSession(context.Background(), session1)
	assert.NoError(t, err)
	err = mcpServer.RegisterSession(context.Background(), session2)
	assert.NoError(t, err)
	err = mcpServer.RegisterSession(context.Background(), session3)
	assert.NoError(t, err)

	// Send notification to all clients
	mcpServer.SendNotificationToAllClients("broadcast", map[string]any{
		"data": "broadcast-data",
	})

	// Check that initialized sessions received the notification
	for _, s := range []*testSession{session1, session2} {
		select {
		case notification := <-s.notificationChan:
			assert.Equal(t, "broadcast", notification.Method)
			assert.Equal(t, "broadcast-data", notification.Params.AdditionalFields["data"])
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected notification not received by session %s", s.id)
		}
	}

	// Check that uninitialized session did not receive the notification
	select {
	case notification := <-session3.notificationChan:
		t.Errorf("Unexpected notification received by uninitialized session: %v", notification)
	case <-time.After(100 * time.Millisecond):
		// Expected, no notification for uninitialized session
	}

	// Send notification to specific client
	err = mcpServer.SendNotificationToSpecificClient("session-1", "specific", map[string]any{
		"data": "specific-data",
	})
	assert.NoError(t, err)

	// Check that only session1 received the notification
	select {
	case notification := <-session1.notificationChan:
		assert.Equal(t, "specific", notification.Method)
		assert.Equal(t, "specific-data", notification.Params.AdditionalFields["data"])
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected notification not received by session1")
	}

	// Check that session2 did not receive the notification
	select {
	case notification := <-session2.notificationChan:
		t.Errorf("Unexpected notification received by session2: %v", notification)
	case <-time.After(100 * time.Millisecond):
		// Expected, no notification for session2
	}
}

// testSession is a simple implementation of the ClientSession interface for testing
type testSession struct {
	id                 string
	notificationChan   chan mcp.JSONRPCNotification
	isInitialized      bool
	initializationFunc func()
}

func (s *testSession) SessionID() string {
	return s.id
}

func (s *testSession) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return s.notificationChan
}

func (s *testSession) Initialize() {
	s.isInitialized = true
	s.initializationFunc()
}

func (s *testSession) Initialized() bool {
	return s.isInitialized
}
