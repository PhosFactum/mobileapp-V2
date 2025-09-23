package models

import "time"

type VaccineResponse struct {
	ID          uint      `json:"id"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `json:"is_completed"`

	// Флаги
	IsRefusal   bool `json:"is_refusal"`
	IsExemption bool `json:"is_exemption"`

	// Поля
	TiterAmount     *int    `json:"titer_amount"`
	MedWithdrawlNum *int    `json:"med_withdrawl_num"`
	Result          *string `json:"result"`

	// Связанные справочники
	Title             *TitleResponse             `json:"title,omitempty"`
	Medication        *MedicationResponse        `json:"medication,omitempty"`
	Dose              *DoseResponse              `json:"dose,omitempty"`
	Number            *NumberResponse            `json:"number,omitempty"`
	CertificateNumber *CertificateNumberResponse `json:"certificate_number,omitempty"`
	BodyPart          *BodyPartResponse          `json:"body_part,omitempty"`
	Method            *MethodResponse            `json:"method,omitempty"`
	Place             *PlaceResponse             `json:"place,omitempty"`
}

type TitleResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type MedicationResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type DoseResponse struct {
	ID    uint    `json:"id"`
	Value float64 `json:"value"`
}

type NumberResponse struct {
	ID    uint `json:"id"`
	Value int  `json:"value"`
}

type CertificateNumberResponse struct {
	ID    uint `json:"id"`
	Value int  `json:"value"`
}

type BodyPartResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type MethodResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type PlaceResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}
