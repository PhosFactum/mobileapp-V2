package interfaces

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type Usecases interface {
	DoctorUsecase
	PatientUsecase
	ReceptionUsecase
	AuthUsecase
	OrganizationUsecase
	PatientGroupUsecase
	ManualUsecase
	VaccineUsecase
	AnalysisOrderUsecase
	FlgUsecase
}

type FlgUsecase interface {
	CreateFlgWithPhoto(ctx context.Context, req *models.CreateFlgRequest) (*models.FlgResponse, *errors.AppError)
}

type AnalysisOrderUsecase interface {
	UpdateAnalysisOrder(ctx context.Context, req *models.UpdateAnalysisOrderRequest) *errors.AppError
}

type ManualUsecase interface {
	GetAllManuals(ctx context.Context) ([]models.ManualResponse, *errors.AppError)
}

type PatientGroupUsecase interface {
	GetPatientGroupsByOrganizationID(orgID uint, search string, page, perPage int) (*models.FilterResponse[[]models.PatientGroupShortResponse], *errors.AppError)
}

type OrganizationUsecase interface {
	GetAllDoctorOrganizations(doctorID uint, search string, page, perPage int) (*models.FilterResponse[[]models.OrganizationShortResponse], *errors.AppError)
}

type ReceptionUsecase interface {
	UpdateReceptionData(ctx context.Context, req *models.UpdateReceptionDataRequest) *errors.AppError
}

type ContactInfoUsecase interface {
	GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, *errors.AppError)
}

type DoctorUsecase interface {
	// CreateDoctor(doctor *models.CreateDoctorRequest) (entities.Doctor, *errors.AppError)
	GetDoctorByID(id uint) (*models.DoctorResponse, *errors.AppError)
	UpdateDoctor(doctor *models.UpdateDoctorRequest) (entities.Doctor, *errors.AppError)
	DeleteDoctor(doctorId uint) *errors.AppError
}

type PatientUsecase interface {
	CreatePatient(ctx context.Context, req models.CreatePatientRequest) (*models.PatientResponse, *errors.AppError)
	GetPatientsByGroup(ctx context.Context, groupID uint) ([]models.PatientResponse, *errors.AppError)
}

type PersonalInfoUsecase interface{}

type AuthUsecase interface {
	LoginDoctor(ctx context.Context, req models.DoctorLoginRequest) (uint, string, *errors.AppError)
	LogoutDoctor(ctx context.Context, token string) *errors.AppError
}

type ConsentUsecase interface {
	SaveConsent(patientID uint, signature []byte) *errors.AppError
	GetSignature(patientID uint) ([]byte, *errors.AppError)
	GetConsentByPatientID(patientID uint) (*entities.ConsentSignature, *errors.AppError)
}

type VaccineUsecase interface {
	CreateVaccine(ctx context.Context, req *models.CreateVaccineRequest) (*entities.Vaccine, *errors.AppError)
	CreateVaccineRefusal(ctx context.Context, req *models.CreateVaccineRefusalRequest) (*entities.VaccineRefusal, *errors.AppError)
	CreateVaccineWithdrawal(ctx context.Context, req *models.CreateVaccineWithdrawalRequest) (*entities.VaccineWithdrawal, *errors.AppError)
	CreateTitr(ctx context.Context, req *models.CreateTitrRequest) (*entities.Titr, *errors.AppError)
}
