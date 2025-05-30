package types

type Comment struct {
	ID      string `json:"id"`
	Recipe  Recipe `json:"recipe"`
	User    User   `json:"user"`
	Comment string `json:"comment"`
}
