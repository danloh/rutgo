// def some struct , will be used in db, handler

package model

// User register user
type User struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Pswmd string `json:"pswmd,omitempty"`
	Email string `json:"email"`
}

// Poem item
type Poem struct {
	ID      int      `json:"id,omitempty"`
	Title   string   `json:"title,omitempty"`
	Author  string   `json:"author,omitempty"`
	Content []string `json:"content,omitempty"`
}

// Draw item
type Draw struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	URL    string `json:"url,omitempty"`
}

// Translate item
type Translate struct {
	ID       int      `json:"id,omitempty"`
	Title    string   `json:"title,omitempty"`
	Author   string   `json:"author,omitempty"`
	Language string   `json:"language,omitempty"`
	Content  []string `json:"content,omitempty"`
}

// Audio item
type Audio struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Language string `json:"language,omitempty"`
	URL      string `json:"url,omitempty"`
}

// Bard item
type Bard struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Era      string `json:"era,omitempty"`
	Language string `json:"language,omitempty"`
	Intro    string `json:"intro,omitempty"`
}
