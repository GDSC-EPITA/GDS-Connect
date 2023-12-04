package models

type User struct {
	Id  int	`json:"id,gt=0"`
	Name string `json:"name"`
	Age  int    `json:"age,gt=0"`
	Gender    *string  `json:"gender,omitempty"`
	Interests []string `json:"interests"`
}
