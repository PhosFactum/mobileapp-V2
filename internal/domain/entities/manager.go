package entities

// Doctor представляет информацию о враче
type Manager struct {
	ID uint `gorm:"primarykey" json:"id"`

	FullName string `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	Phone    string `gorm:"unique;not null" json:"phone" example:"+79991234567"`
}
