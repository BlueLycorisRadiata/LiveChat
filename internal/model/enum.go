package model

type ConversationType string
type ParticipantRole string
type MessageType string

const (
	ConversationTypePrivate ConversationType = "private"
	ConversationTypeGroup   ConversationType = "group"
)

const (
	ParticipantRoleOwner  ParticipantRole = "owner"
	ParticipantRoleMember ParticipantRole = "member"
)

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeFile   MessageType = "file"
	MessageTypeSystem MessageType = "system"
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
	case ConversationTypePrivate, ConversationTypeGroup:
		return true
	default:
		return false
	}
}
