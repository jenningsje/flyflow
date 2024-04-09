package models

type Model struct {
	BaseModel
	UserId            uint   `json:"user_id"`
	ModelName         string `json:"model_name"`
	InternalModelName string `json:"-"`
	APIUrl            string `json:"-"`
	APIKey            string `json:"-"`
	Format            string `json:"format"`
}