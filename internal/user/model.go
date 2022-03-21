package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email"`
	Username     string `json:"usename" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
}

type CreateUserDTO struct {
	Email        string `json:"email"`
	Username     string `json:"usename"`
	PasswordHash string `json:"password"`
}
