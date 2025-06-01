package models

type Todo struct {
	BaseModel
	Title       string `json:"title" gorm:"not null" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed" gorm:"default:false"`
}

func (Todo) TableName() string {
	return "todos"
}
