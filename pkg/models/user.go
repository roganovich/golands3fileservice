package models

import(
	"time"
	jwt "github.com/golang-jwt/jwt"
)

type User struct {
	ID    		int64    				`json:"id"`
	Name  		string 				`json:"name"`						// ФИО
	Email 		string 				`json:"email"`						// Email
	Password 	string 				`json:"password"`					// Пароль
	Status 		int 				`json:"status"`						// Статус
	CreatedAt   time.Time       	`json:"created_at"`              	// Дата создания
	UpdatedAt   time.Time       	`json:"updated_at"`             	// Дата последнего обновления
	DeletedAt 	*time.Time 	    	`json:"deleted_at, omitempty"`		// Дата удаления
}

type UserView struct {
	ID    		int64    			`json:"id"`
	Name  		string 				`json:"name"`						// ФИО
	Email 		string 				`json:"email"`						// Email
	Role 		Role 				`json:"role"`						// Роль
	Status 		int 				`json:"status"`						// Статус
	CreatedAt   time.Time       	`json:"created_at"`              	// Дата создания
}

type CreateUserRequest struct {
	ID    		int64
	Name 		string 				`json:"name" validate:"required,min=3,max=128"`
	Email 		string 				`json:"email" validate:"required,email"`
	Password	string 				`json:"password" validate:"required,min=4,max=128"`
}

type UpdateUserRequest struct {
	ID    		int64    				`json:"id"`
	Name 		string 				`json:"name" validate:"required,min=3,max=128"`
	Email 		string 				`json:"email" validate:"required,email"`
	Password	string 				`json:"password" validate:"required,min=4,max=128"`
}


type LoginUserRequest struct {
	Email 		string 				`json:"email" validate:"required,email"`
	Password	string 				`json:"password" validate:"required,min=4,max=128"`
}

// Claims содержит информацию, которую мы хотим включить в токен
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}