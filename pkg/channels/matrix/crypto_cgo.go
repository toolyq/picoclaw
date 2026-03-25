//go:build cgo

package matrix

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"go.mau.fi/util/dbutil"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"

	"github.com/sipeed/picoclaw/pkg/logger"
)

func (c *MatrixChannel) initCrypto(ctx context.Context) error {
	logger.InfoC("matrix", "Initializing crypto helper")

	// Ensure the crypto database directory exists
	if err := os.MkdirAll(c.cryptoDbPath, 0o700); err != nil {
		return fmt.Errorf("create crypto database directory: %w", err)
	}

	// Create database with sqlite driver (modernc.org/sqlite)
	dbPath := filepath.Join(c.cryptoDbPath, dbName)
	connStr := "file:" + dbPath + "?_foreign_keys=on"

	db, err := sql.Open(sqliteDriver, connStr)
	if err != nil {
		return fmt.Errorf("open crypto database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Execute PRAGMA statements
	// This is equivalent to the "sqlite3-fk-wal" dialect used by cryptohelper
	pragmaStmts := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
	}
	for _, pragma := range pragmaStmts {
		if _, err = db.ExecContext(ctx, pragma); err != nil {
			_ = db.Close()
			return fmt.Errorf("execute %s: %w", pragma, err)
		}
	}

	// Wrap with dbutil for dialect support
	wrappedDB, err := dbutil.NewWithDB(db, sqliteDriver)
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("wrap database: %w", err)
	}

	ch, err := cryptohelper.NewCryptoHelper(c.client, []byte(c.config.CryptoPassphrase), wrappedDB)
	if err != nil {
		return fmt.Errorf("create crypto helper: %w", err)
	}

	if c.client.DeviceID == "" {
		resp, whoamiErr := c.client.Whoami(ctx)
		if whoamiErr != nil {
			_ = db.Close()
			return fmt.Errorf("get device ID via whoami: %w", whoamiErr)
		}
		c.client.DeviceID = resp.DeviceID
	}

	if err = ch.Init(ctx); err != nil {
		ch.Close()
		return fmt.Errorf("init crypto helper: %w", err)
	}

	c.client.Crypto = ch
	c.cryptoHelper = ch

	logger.InfoC("matrix", "Crypto helper initialized successfully")
	return nil
}

func (c *MatrixChannel) decryptEvent(ctx context.Context, evt *event.Event) (*event.MessageEventContent, bool) {
	ch, ok := c.cryptoHelper.(*cryptohelper.CryptoHelper)
	if !ok || ch == nil {
		logger.DebugCF("matrix", "Received encrypted message but crypto is not enabled", map[string]any{
			"room_id": evt.RoomID.String(),
		})
		return nil, false
	}

	decrypted, err := ch.Decrypt(ctx, evt)
	if err != nil {
		logger.WarnCF("matrix", "Failed to decrypt message", map[string]any{
			"room_id": evt.RoomID.String(),
			"error":   err.Error(),
		})
		return nil, false
	}

	if decrypted.Type != event.EventMessage {
		logger.DebugCF("matrix", "Decrypted event is not a message event", map[string]any{
			"room_id": evt.RoomID.String(),
			"type":    decrypted.Type.String(),
		})
		return nil, false
	}

	return decrypted.Content.AsMessage(), true
}

func (c *MatrixChannel) closeCrypto() {
	if ch, ok := c.cryptoHelper.(*cryptohelper.CryptoHelper); ok && ch != nil {
		ch.Close()
		c.cryptoHelper = nil
		c.client.Crypto = nil
	}
}
