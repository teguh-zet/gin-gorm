package models

type Book struct{
	ID uint			`json:"id" gorm:"primaryKey;autoIncrement"`
	Title string	`json:"title" gorm:"not null"`
	Author string	`json:"author"`
}


func(Book) TableName() string{
	return  "books"
}

type CreateBookRequest struct{
	Title string 	`json:"title" binding:"required,min=2,max=100"`
	Author string 	`json:"author" binding:"required,min=2,max=100"`
}

