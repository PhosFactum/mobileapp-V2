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
	OrganizationUseCase
	PatientGroupUseCase
}

type PatientGroupUseCase interface {
	GetPatientGroupsByCodeOrOrgTitle(search string, page, perPage int) (*models.FilterResponse[[]models.PatientGroupShortResponse], error)
	GetPatientGroupsByOrganizationID(orgID uint, page, perPage int) (*models.FilterResponse[[]models.PatientGroupShortResponse], error)
}

type OrganizationUseCase interface {
	GetAllOrganizations(doctorID uint, page, perPage int) (*models.FilterResponse[[]models.OrganizationShortResponse], error)
}

type ReceptionUsecase interface {
	// 	GetHospitalReceptionsByPatientID(patientId uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError)
	// 	UpdateReceptionHospital(id uint, input *models.UpdateReceptionHospitalRequest) (models.ReceptionHospitalResponse, *errors.AppError)
	// 	GetHospitalReceptionsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError)
	// 	GetHospitalPatientsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]entities.Patient], *errors.AppError)
	// 	GetReceptionHospitalByID(hospID uint) (models.ReceptionFullResponse, error)
	// 	UpdateReceptionHospitalStatus(id uint, newStatus string) (entities.Reception, error)
}

type ContactInfoUsecase interface {
	CreateContactInfo(input *models.CreateContactInfoRequest) (entities.ContactInfo, *errors.AppError)
	GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, *errors.AppError)
}

type DoctorUsecase interface {
	// CreateDoctor(doctor *models.CreateDoctorRequest) (entities.Doctor, *errors.AppError)
	GetDoctorByID(doctorId uint) (entities.Doctor, *errors.AppError)
	UpdateDoctor(doctor *models.UpdateDoctorRequest) (entities.Doctor, *errors.AppError)
	DeleteDoctor(doctorId uint) *errors.AppError
}

type PatientUsecase interface {
	CreatePatient(patientData *models.CreatePatientData) (*entities.Patient, error)
	GetPatientByID(id uint) (entities.Patient, *errors.AppError)
	UpdatePatient(input *models.UpdatePatientRequest) (entities.Patient, *errors.AppError)
	DeletePatient(id uint) *errors.AppError
	GetPatientsByGroup(page, perPage int, group_id uint) (models.FilterResponse[[]entities.Patient], *errors.AppError)
}

type PersonalInfoUsecase interface{}

type AuthUsecase interface {
	LoginDoctor(ctx context.Context, phone, password string) (uint, string, *errors.AppError)
	LogoutDoctor(ctx context.Context, token string) *errors.AppError
}

type ConsentUsecase interface {
	SaveConsent(patientID uint, signature []byte) *errors.AppError
	GetSignature(patientID uint) ([]byte, *errors.AppError)
	GetConsentByPatientID(patientID uint) (*entities.ConsentSignature, *errors.AppError)
}
