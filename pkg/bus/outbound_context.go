package bus

import "strings"

// NewOutboundContext builds the minimal normalized addressing context required
// to deliver an outbound text message or reply.
func NewOutboundContext(channel, chatID, replyToMessageID string) InboundContext {
	return normalizeInboundContext(InboundContext{
		Channel:          strings.TrimSpace(channel),
		ChatID:           strings.TrimSpace(chatID),
		ReplyToMessageID: strings.TrimSpace(replyToMessageID),
	})
}

// NormalizeOutboundMessage ensures Context is normalized and keeps convenience
// mirrors in sync for runtime consumers.
func NormalizeOutboundMessage(msg OutboundMessage) OutboundMessage {
	msg.Context = normalizeInboundContext(msg.Context)
	msg.Channel = msg.Context.Channel
	msg.ChatID = msg.Context.ChatID
	msg.Scope = cloneOutboundScope(msg.Scope)
	if msg.Context.ReplyToMessageID == "" {
		msg.Context.ReplyToMessageID = strings.TrimSpace(msg.ReplyToMessageID)
	}
	msg.ReplyToMessageID = msg.Context.ReplyToMessageID
	return msg
}

// NormalizeOutboundMediaMessage ensures media outbound messages also carry a
// normalized context while keeping convenience mirrors in sync.
func NormalizeOutboundMediaMessage(msg OutboundMediaMessage) OutboundMediaMessage {
	msg.Context = normalizeInboundContext(msg.Context)
	msg.Channel = msg.Context.Channel
	msg.ChatID = msg.Context.ChatID
	msg.Scope = cloneOutboundScope(msg.Scope)
	return msg
}

func cloneOutboundScope(scope *OutboundScope) *OutboundScope {
	if scope == nil {
		return nil
	}
	cloned := *scope
	if len(scope.Dimensions) > 0 {
		cloned.Dimensions = append([]string(nil), scope.Dimensions...)
	}
	if len(scope.Values) > 0 {
		cloned.Values = make(map[string]string, len(scope.Values))
		for key, value := range scope.Values {
			cloned.Values[key] = value
		}
	}
	return &cloned
}
