package interfaces

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

type Repository interface {
	AuthRepository
	DoctorRepository
	PatientRepository
	TxRepository
	ConsentSignatureRepository
	OrganizationRepository
	PatientGroupRepository
	ReceptionRepository
	VaccineRepository
	ManualRepository
}

type PatientGroupRepository interface {
	GetPatientGroupsByCodeOrOrgTitle(search string, page, perPage int) ([]entities.PatientGroup, int64, error)
	GetPatientGroupsByOrganizationID(orgID uint, page, perPage int) ([]entities.PatientGroup, int64, error)
}

type ManualRepository interface {
	GetManualValueByTypeAndID(id uint, ref_type entities.ReferenceType) (string, error)
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
}

// updated to match the new structured
type PatientRepository interface {

	// Работа с зависимыми сущностями
	CreateContactInfo(tx *gorm.DB, contactInfo *entities.ContactInfo) *errors.AppError
	CreatePersonalInfo(tx *gorm.DB, personalInfo *entities.PersonalInfo) *errors.AppError
	CreatePatientStatistics(tx *gorm.DB, stats *entities.PatientStatistics) *errors.AppError

	// Работа с пациентом
	GetPatientsByGroup(group_id uint) ([]entities.Patient, *errors.AppError)
	CreatePatient(tx *gorm.DB, patient *entities.Patient) *errors.AppError
	PreloadPatientWithSpecializations(tx *gorm.DB, patientID uint) (*entities.Patient, *errors.AppError)

	CreateAnalysisOrder(tx *gorm.DB, order *entities.AnalysisOrder) *errors.AppError
	UpdateAnalysisOrder(tx *gorm.DB, order *entities.AnalysisOrder) *errors.AppError

	// Работа со связями
	CacheSpecializations(tx *gorm.DB, patient *entities.Patient, specializations []entities.Specialization) *errors.AppError

	// Работа с приёмами
	CreateReceptions(tx *gorm.DB, receptions []entities.Reception) *errors.AppError

	// Получение шаблонов по вредному фактору
	GetReceptionTemplatesByHarmPointID(tx *gorm.DB, harmPointID uint) ([]entities.ReceptionTemplate, *errors.AppError)
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
