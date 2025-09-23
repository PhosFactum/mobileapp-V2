package interfaces

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
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
	FLGRepository
}

type PatientGroupRepository interface {
	GetPatientGroupsByCodeOrOrgTitle(search string, page, perPage int) ([]entities.PatientGroup, int64, error)
	GetPatientGroupsWithPatientsByOrganizationID(orgID uint, page, perPage int) ([]entities.PatientGroup, int64, error)
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
	// 	CreateReceptionHospital(reception entities.Reception) error
	// 	UpdateReception(id uint, updateMap map[string]interface{}) (uint, error)
	// 	DeleteReception(id uint) error
	// 	GetReceptionByID(id uint) (entities.Reception, error)
	// 	GetReceptionByPatientID(patientID uint) ([]entities.Reception, error)
}

// updated to match the new structured
type PatientRepository interface {
	CreatePatient(patient entities.Patient) (uint, error)
	UpdatePatient(id uint, updateMap map[string]interface{}) (uint, error)
	DeletePatient(id uint) error
	GetPatientByID(id uint) (entities.Patient, error)
	GetAllPatients(page, count int, queryFilter string, queryOrder string, filterParams []interface{}) ([]entities.Patient, int64, error)
	GetPatientsByFullName(name string) ([]entities.Patient, error)
	GetPatientByIDWithTx(tx *gorm.DB, id uint) (*entities.Patient, error)
	UpdatePatientWithTx(tx *gorm.DB, id uint, updateMap map[string]interface{}) (uint, error)
	GetPatientsByGroup(page, perPage int, group_id uint) ([]entities.Patient, int64, error)
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

type FLGRepository interface {
	CreateFLG(flg entities.FLG) (uint, error)
	UpdateFLG(id uint, updateMap map[string]interface{}) (*entities.FLG, error)
}
