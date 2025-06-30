package structures

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"displayName"`
	ProfilePicture string `json:"profilePicture"`
}

type Reaction struct {
	ID    int    `json:"id"`
	Emoji string `json:"emoji"`
	User  User   `json:"user"`
}

type Message struct {
	ID             int        `json:"id"`
	ConversationID int        `json:"conversation_id"`
	Content        string     `json:"content"`
	IsForwarded    bool       `json:"isForwarded"`
	MediaType      string     `json:"mediaType"`
	Reactions      []Reaction `json:"reactions"`
	Sender         User       `json:"sender"`
	Status         string     `json:"status"`
	Timestamp      string     `json:"timestamp"`
}

type Conversation struct {
	ID      int `json:"id"`
	User1ID int `json:"user1_id"`
	User2ID int `json:"user2_id"`
}

type Group struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Members []User `json:"members"`
}

type SessionRequest struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
}
