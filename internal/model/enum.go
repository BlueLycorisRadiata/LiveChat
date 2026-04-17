package model

type ConversationType string
type ParticipantRole string
type MessageType string
type MessageRole string

const (
	ConversationTypePrivate ConversationType = "private"
	ConversationTypeGroup   ConversationType = "group"
	ConversationTypeAI      ConversationType = "ai"
)

const (
	ParticipantRoleOwner  ParticipantRole = "owner"
	ParticipantRoleAdmin  ParticipantRole = "admin"
	ParticipantRoleMember ParticipantRole = "member"
)

func (r ParticipantRole) IsValid() bool {
	switch r {
	case ParticipantRoleOwner, ParticipantRoleAdmin, ParticipantRoleMember:
		return true
	default:
		return false
	}
}

func (r ParticipantRole) CanManageMembers() bool {
	return r == ParticipantRoleOwner || r == ParticipantRoleAdmin
}

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeFile   MessageType = "file"
	MessageTypeSystem MessageType = "system"
)

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

func (t MessageType) IsValid() bool {
	switch t {
	case MessageTypeText, MessageTypeImage, MessageTypeFile, MessageTypeSystem:
		return true
	default:
		return false
	}
}

func (t ConversationType) IsValid() bool {
	switch t {
	case ConversationTypePrivate, ConversationTypeGroup, ConversationTypeAI:
		return true
	default:
		return false
	}
}

func (r MessageRole) IsValid() bool {
	switch r {
	case RoleUser, RoleAssistant, RoleSystem:
		return true
	default:
		return false
	}
}
