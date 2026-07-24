package response

import "vault/src/features/chat/domain/entities"

type ConversationSummaryResponse struct {
	OtherUserID        string              `json:"other_user_id"`
	OtherUserName      string              `json:"other_user_name"`
	OtherUserAvatarURL string              `json:"other_user_avatar_url"`
	LastMessage        ChatMessageResponse `json:"last_message"`
	UnreadCount        int                 `json:"unread_count"`
}

func FromConversationSummary(s entities.ConversationSummary) ConversationSummaryResponse {
	return ConversationSummaryResponse{
		OtherUserID:        s.OtherUserID,
		OtherUserName:      s.OtherUserName,
		OtherUserAvatarURL: s.OtherUserAvatarURL,
		LastMessage:        FromEntity(s.LastMessage),
		UnreadCount:        s.UnreadCount,
	}
}

func FromConversationSummaries(list []entities.ConversationSummary) []ConversationSummaryResponse {
	out := make([]ConversationSummaryResponse, 0, len(list))
	for _, s := range list {
		out = append(out, FromConversationSummary(s))
	}
	return out
}
