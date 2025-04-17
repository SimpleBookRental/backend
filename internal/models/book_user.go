package models

// BookCreateRequest represents the request to create a book
// @Description Book create with user request
type BookCreateRequest struct {
	Title       string `json:"title" example:"The Great Gatsby"`         // Book title
	Author      string `json:"author" example:"F. Scott Fitzgerald"`     // Book author
	ISBN        string `json:"isbn" example:"978-1234567890"`            // Book ISBN
	Description string `json:"description" example:"A classic novel"`    // Book description
}

// BookTransferRequest represents the request to transfer a book
// @Description Book transfer request
type BookTransferRequest struct {
	FromUserID string `json:"from_user_id" binding:"required,uuid" example:"b1a2c3d4-e5f6-7890-abcd-1234567890ab"` // From user ID
	ToUserID   string `json:"to_user_id" binding:"required,uuid" example:"c2b3a4d5-e6f7-8901-bcda-2345678901bc"`   // To user ID
}
