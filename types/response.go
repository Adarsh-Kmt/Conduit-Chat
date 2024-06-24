package types

type OutgoingChatRequestListResponse struct {
	OutgoingChatRequestList []ChatRequest `json:"outgoingChatRequestList" bson:"OutgoingChatRequestList"`
}

type IncomingChatRequestListResponse struct {
	IncomingChatRequestList []ChatRequest `json:"incomingChatRequestList" bson:"IncomingChatRequestList"`
}
