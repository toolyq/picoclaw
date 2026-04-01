//go:build !cgo

package matrix

import (
	"context"

	"github.com/sipeed/picoclaw/pkg/logger"
	"maunium.net/go/mautrix/event"
)

func (c *MatrixChannel) initCrypto(ctx context.Context) error {
	logger.WarnC("matrix", "Encryption requested but this build does not support it (CGO disabled)")
	return nil
}

func (c *MatrixChannel) decryptEvent(ctx context.Context, evt *event.Event) (*event.MessageEventContent, bool) {
	logger.DebugCF("matrix", "Received encrypted message but crypto is not supported in this build", map[string]any{
		"room_id": evt.RoomID.String(),
	})
	return nil, false
}

func (c *MatrixChannel) closeCrypto() {
	// Nothing to do
}
