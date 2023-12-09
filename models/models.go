package models

type User struct {
	Name            string   `json:"name"`
	Age             int      `json:"age,gt=0"`
	Gender          *string  `json:"gender,omitempty"`
	Interests       []string `json:"interests"`
	Anonymous       *bool    `json:"anonymous,omitempty"`
	UsersVisibility []string `json:"usersVisibility,omitempty"`
}
