package entities

type Specialization struct {
	ID      uint      `gorm:"primarykey" json:"id"`
	Title   string    `gorm:"unique;not null" json:"title" example:"Терапевт"`
	Doctor  []Doctor  `gorm:"many2many:doctor_specializations" json:"-"`
	Patient []Patient `gorm:"many2many:patients_specializations; default:null;" json:"-"`
}
