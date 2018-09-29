package models

// Message struct
type Message struct {
	ID int `json:"id" storm:"id,increment"`
	Model
	Type   string `json:"type"`
	Text   string `json:"text"`
	SRC    string `json:"src,omitempty"`
	From   User   `json:"from"`
	FromId int    `json:"from_id"`
}

// Message request struct for Echo to bind to (not used for image messages)
type MessageRequest struct {
	Type  string
	Text  string
	Token string
}
