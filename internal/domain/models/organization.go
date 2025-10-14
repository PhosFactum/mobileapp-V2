package models

type OrganizationShortResponse struct {
	ID      uint            `json:"id"`
	Title   string          `json:"title"`
	Manager ManagerResponse `json:"manager"`
}

// Manager представляет информацию о мэнэджере организации
type ManagerResponse struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
