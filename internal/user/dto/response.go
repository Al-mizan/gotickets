package dto

type Response struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" `
	Email     string `json:"email"  `
	CreatedAt string `json:"created_at"`
}
