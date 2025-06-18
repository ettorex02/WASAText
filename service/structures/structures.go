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
	ID          int        `json:"id"`
	Content     string     `json:"content"`
	IsForwarded bool       `json:"isForwarded"`
	MediaType   string     `json:"mediaType"`
	Reactions   []Reaction `json:"reactions"`
	Sender      User       `json:"sender"`
	Status      string     `json:"status"`
	Timestamp   string     `json:"timestamp"`
}

type Conversation struct {
	ID           int     `json:"id"`
	LastMessage  Message `json:"lastMessage"`
	Participants []User  `json:"participants"`
}

type Group struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Members []User `json:"members"`
}
