package models

type User struct {
	FirstName string `json:"firstname" validate:"required,min=2,max=100"`
	LastName  string `json:"lastname" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"email,required"`
	Phone     string `json:"phone" validate:"required"`
	UserName  string `json:"username" validate:"reuired,min=6"`
	Password  string `json:"password" validate:"required,min=6"`
	// Role 	  string `json:"role" validate:"required, eq=ADMIN|USER"`
}

type Notes struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
