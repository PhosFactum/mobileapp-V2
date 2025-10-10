package interfaces

import (
	"context"
	"encoding/json"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type Repository interface {
	AuthRepository
	DoctorRepository
	PatientRepository
	ConsentSignatureRepository
	OrganizationRepository
	PatientGroupRepository
	ReceptionRepository
	VaccineRepository
	ManualRepository
	AnalysisRepository
	AnalysisOrderRepository
}

type PatientGroupRepository interface {
	GetPatientGroupsByOrganizationID(orgID uint, search string, page, perPage int) ([]entities.PatientGroup, int64, error)
}

type ManualRepository interface {
	GetManualValueByTypeAndID(ctx context.Context, id uint, refType entities.ReferenceType) (string, error)
	GetManualValuesByType(ctx context.Context, refType entities.ReferenceType) ([]string, error)
	GetAllManuals(ctx context.Context) ([]entities.Manual, error)
}

type OrganizationRepository interface {
	GetAllDoctorOrganizations(doctorID uint, search string, page, perPage int) ([]entities.Organization, int64, error)
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
type ReceptionRepository interface {
	GetReceptionTemplatesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.ReceptionTemplate, error)
	CreateReceptions(ctx context.Context, receptions []entities.Reception) error
	GetReceptionTemplatesByCodes(ctx context.Context, codes []string) ([]entities.ReceptionTemplate, error)
	GetTemplateByReceptionID(ctx context.Context, receptionID uint) (*entities.ReceptionTemplate, error)
	UpdateReceptionData(ctx context.Context, receptionID uint, data json.RawMessage) error
}

type AnalysisRepository interface {
	GetAnalysesByCodes(ctx context.Context, codes []string) ([]entities.Analysis, error)
	GetAnalysesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.Analysis, error)

	GetAnalysisByID(ctx context.Context, id uint) (*entities.Analysis, error)
	GetAllAnalysisIDs(ctx context.Context) ([]uint, error)
}

type AnalysisOrderRepository interface {
	UpdateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error
	CreateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error
	CreateAnalysisItems(ctx context.Context, items []entities.AnalysisOrderItem) error

	GetByID(ctx context.Context, id uint) (*entities.AnalysisOrder, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID uint) ([]entities.AnalysisOrderItem, error)
	CreateOrderItems(ctx context.Context, items []entities.AnalysisOrderItem) error
	UpdateOrderItem(ctx context.Context, item entities.AnalysisOrderItem) error
	DeleteOrderItems(ctx context.Context, itemIDs []uint) error
}

// updated to match the new structured
type PatientRepository interface {
	// Работа с пациентом
	GetPatientsByGroup(ctx context.Context, groupID uint) ([]entities.Patient, *errors.AppError)
	CreateContactInfo(ctx context.Context, contactInfo *entities.ContactInfo) error
	CreatePersonalInfo(ctx context.Context, personalInfo *entities.PersonalInfo) error
	CreatePatient(ctx context.Context, patient *entities.Patient) error
	CacheSpecializations(ctx context.Context, patient *entities.Patient, specializations []entities.Specialization) error
	CreatePatientStatistics(ctx context.Context, stats *entities.PatientStatistics) error
	PreloadPatientWithSpecializations(ctx context.Context, patientID uint) (*entities.Patient, error)
}

// updated to match the new structure
type ContactInfoRepository interface {
	CreateContactInfo(info entities.ContactInfo) (uint, error)
	UpdateContactInfo(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteContactInfo(id uint) error
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
	CreateVaccine(ctx context.Context, vaccine *entities.Vaccine) error
	CreateVaccineRefusal(ctx context.Context, refusal *entities.VaccineRefusal) error
	CreateVaccineWithdrawal(ctx context.Context, withdrawal *entities.VaccineWithdrawal) error
	CreateTitr(ctx context.Context, titr *entities.Titr) error
}

// updated to match the new structure
type PersonalInfoRepository interface {
}
