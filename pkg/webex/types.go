package webex

import "time"

type Webhook struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	TargetUrl string    `json:"targetUrl"`
	Resource  string    `json:"resource"`
	Event     string    `json:"event"`
	Filter    string    `json:"filter,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	Items     []Webhook `json:"items"`
}

type MessageWebhookBody struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	TargetUrl string    `json:"targetUrl,omitempty"`
	Resource  string    `json:"resource,omitempty"`
	Event     string    `json:"event,omitempty"`
	Filter    string    `json:"filter,omitempty"`
	OrgId     string    `json:"orgId,omitempty"`
	CreatedBy string    `json:"createdBy,omitempty"`
	AppId     string    `json:"appId,omitempty"`
	OwnedBy   string    `json:"ownedBy,omitempty"`
	Status    string    `json:"status,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	ActorId   string    `json:"actorId,omitempty"`
	Data      struct {
		Id          string    `json:"id,omitempty"`
		RoomId      string    `json:"roomId,omitempty"`
		RoomType    string    `json:"roomType,omitempty"`
		PersonId    string    `json:"personId,omitempty"`
		PersonEmail string    `json:"personEmail,omitempty"`
		Created     time.Time `json:"created,omitempty"`
	} `json:"data"`
}

type Message struct {
	Id       string `json:"id,omitempty"`
	RoomId   string `json:"roomId,omitempty"`
	RoomType string `json:"roomType,omitempty"`
	Text     string `json:"text,omitempty"`
	Markdown string `json:"markdown,omitempty"`

	// Use ToPersonEmail for posts. PersonEmail (no to) is returned in gets.
	// I don't make the rules. Thanks Webex <3
	ToPersonEmail string `json:"toPersonEmail,omitempty"`
	PersonEmail   string `json:"personEmail,omitempty"`
	PersonId      string `json:"personId,omitempty"`

	// This does not go to the request; it is used for logging in the server message requests,
	// specifically for rooms so the name can be identified
	RecipientType string
	RecipientName string
}

type ListRoomsResp struct {
	Items []Room `json:"items"`
}

type Room struct {
	Id           string    `json:"id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	IsLocked     bool      `json:"isLocked"`
	LastActivity time.Time `json:"lastActivity"`
	CreatorId    string    `json:"creatorId"`
	Created      time.Time `json:"created"`
	OwnerId      string    `json:"ownerId"`
	IsPublic     bool      `json:"isPublic"`
	IsReadOnly   bool      `json:"isReadOnly"`
}

type ListPeopleResp struct {
	NotFoundIds any      `json:"notFoundIds"`
	Items       []Person `json:"items"`
}

type Person struct {
	ID           string   `json:"id"`
	Emails       []string `json:"emails"`
	PhoneNumbers []struct {
		Type    string `json:"type"`
		Value   string `json:"value"`
		Primary bool   `json:"primary"`
	} `json:"phoneNumbers"`
	DisplayName  string    `json:"displayName"`
	NickName     string    `json:"nickName"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Avatar       string    `json:"avatar"`
	OrgID        string    `json:"orgId"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	LastActivity time.Time `json:"lastActivity"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
}
