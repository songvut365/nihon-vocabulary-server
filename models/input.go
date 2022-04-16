package models

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ResetPasswordInput struct {
	Email string `json:"email"`
}

type UpdateUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"LastName"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type VocabularyInput struct {
	Japanese Japanese `json:"japanese"`
	Thai     []string `json:"thai"`
	English  []string `json:"english"`
	Examples []string `json:"examples"`
	Image    string   `json:"image"`
	Voice    string   `json:"vocie"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags"`
}	

