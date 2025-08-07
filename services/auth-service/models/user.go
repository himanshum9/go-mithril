package models

type User struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
}

func NewUser(userID, email, role string) *User {
    return &User{
        UserID: userID,
        Email:  email,
        Role:   role,
    }
}