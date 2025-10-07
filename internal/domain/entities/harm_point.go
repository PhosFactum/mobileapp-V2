package entities

type HarmPoint struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`

	ReceptionTemplates []ReceptionTemplate `gorm:"many2many:harm_point_reception_templates;"`
	Analyses           []Analysis          `gorm:"many2many:harm_point_analyses;"`
}
