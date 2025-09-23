package interfaces

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	AuthRepository
	DoctorRepository
	PatientRepository
	ContactInfoRepository
	PersonalInfoRepository
	TxRepository
	ConsentSignatureRepository
	OrganizationRepository
	PatientGroupRepository
	ReceptionRepository
	VaccineRepository
}

type PatientGroupRepository interface {
	GetPatientGroupsByCodeOrOrgTitle(search string, page, perPage int) ([]entities.PatientGroup, int64, error)
	GetPatientGroupsByOrganizationID(orgID uint, page, perPage int) ([]entities.PatientGroup, int64, error)
}

type OrganizationRepository interface {
	GetAllOrganizations(doctorID uint, page, perPage int) ([]entities.Organization, int64, error)
}

type TxRepository interface {
	BeginTx() (*gorm.DB, error)
	CommitTx(tx *gorm.DB) error
	RollbackTx(tx *gorm.DB) error
}

// updated to match the new structure
type DoctorRepository interface {
	UpdateDoctor(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteDoctor(id uint) error
	GetDoctorByID(id uint) (entities.Doctor, error)
	GetDoctorName(id uint) (string, error)
	GetDoctorByLogin(login string) (entities.Doctor, error)

	// GetDoctorSpecialization(id uint) (string, error)
	GetDoctorPassHash(id uint) (string, error)
}

// updated to match the new structure
type PersonalInfoRepository interface {
	CreatePersonalInfo(info entities.PersonalInfo) (uint, error)
	UpdatePersonalInfo(id uint, updateMap map[string]interface{}) (uint, error)
	DeletePersonalInfo(id uint) error

	GetPersonalInfoByID(id uint) (entities.PersonalInfo, error)
	GetPersonalInfoByPatientID(patientID uint) (entities.PersonalInfo, error)
	UpdatePersonalInfoByPatientID(id uint, updateMap map[string]interface{}) (uint, error)
	GetPersonalInfoByPatientIDWithTx(tx *gorm.DB, patientID uint) (*entities.PersonalInfo, error)
	UpdatePersonalInfoByPatientIDWithTx(tx *gorm.DB, patientID uint, updateMap map[string]interface{}) (uint, error)
}

// updated to match the new structure
type ReceptionRepository interface {
	GetPatientReceptionsByPatientID(patientID uint) ([]entities.Reception, error)
	// 	CreateReceptionHospital(reception entities.Reception) error
	// 	UpdateReception(id uint, updateMap map[string]interface{}) (uint, error)
	// 	DeleteReception(id uint) error
	// 	GetReceptionByID(id uint) (entities.Reception, error)
	// 	GetReceptionByPatientID(patientID uint) ([]entities.Reception, error)
}

// updated to match the new structured
type PatientRepository interface {
	CreatePatient(patientData *models.CreatePatientData, group_id uint) (*entities.Patient, error)
	GetPatientsByGroup(group_id uint) ([]entities.Patient, error)
}

// updated to match the new structure
type ContactInfoRepository interface {
	CreateContactInfo(info entities.ContactInfo) (uint, error)
	UpdateContactInfo(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteContactInfo(id uint) error

	GetContactInfoByID(id uint) (entities.ContactInfo, error)
	GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, error)
	UpdateContactInfoByPatientID(id uint, updateMap map[string]interface{}) (uint, error)
	CreateContactInfoWithTx(tx *gorm.DB, info entities.ContactInfo) (uint, error)
	GetContactInfoByIDWithTx(tx *gorm.DB, id uint) (*entities.ContactInfo, error)
	UpdateContactInfoByIDWithTx(tx *gorm.DB, id uint, updateMap map[string]interface{}) (uint, error)
}

type AuthRepository interface {
	GetByLogin(ctx context.Context, login string) (*entities.Doctor, error)
}

type ConsentSignatureRepository interface {
	SaveSignature(patientID uint, signature []byte) error
	GetSignature(patientID uint) ([]byte, error)
	GetByPatientID(patientID uint) (*entities.ConsentSignature, error)
}

type VaccineRepository interface {
}

type AnalysisRepository interface {
}
