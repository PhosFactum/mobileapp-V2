package usecases

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type DoctorUsecase struct {
	repo interfaces.DoctorRepository
}

func NewDoctorUsecase(repo interfaces.DoctorRepository) interfaces.DoctorUsecase {
	return &DoctorUsecase{repo: repo}
}

// func (u *DoctorUsecase) CreateDoctor(doctor *models.CreateDoctorRequest) (entities.Doctor, *errors.AppError) {

// 	log.Println("before hash Pass  for Create Doctor")
// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return entities.Doctor{}, errors.NewAppError(400, "error create doctor", err, true)
// 	}
// 	log.Println("hash Pass  for Create Doctor")
// 	log.Println("")
// 	createDoctor := entities.Doctor{
// 		FullName:         doctor.FullName,
// 		Phone:            doctor.Phone,
// 		PasswordHash:     string(passwordHash),
// 		SpecializationID: doctor.SpecializationID,
// 	}

// 	createdDoctorID, errAp := u.repo.CreateDoctor(createDoctor)
// 	if errAp != nil {
// 		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to create doctor", err, true)
// 	}
// 	log.Println("Create Doctor in usace")

// 	createdDoctor, errAp := u.repo.GetDoctorByID(createdDoctorID)
// 	if errAp != nil {
// 		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
// 	}
// 	log.Println("Create Doctor in usace")
// 	return createdDoctor, nil
// }

func (u *DoctorUsecase) GetDoctorByID(id uint) (*models.DoctorResponse, *errors.AppError) {
	// 1. Получаем доктора из репозитория (с предзагрузкой специализаций)
	doc, err := u.repo.GetDoctorByID(id)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get doctor",
			err,
			true,
		)
	}

	// 2. Маппим специализации
	specializations := make([]models.SpecializationResponse, len(doc.Specializations))
	for i, spec := range doc.Specializations {
		specializations[i] = models.SpecializationResponse{
			ID:    spec.ID,
			Title: spec.Title,
		}
	}

	// 3. Создаём и возвращаем модель ответа
	response := &models.DoctorResponse{
		ID:              doc.ID,
		FullName:        doc.FullName,
		Specializations: specializations,
	}

	return response, nil
}

func (u *DoctorUsecase) UpdateDoctor(input *models.UpdateDoctorRequest) (entities.Doctor, *errors.AppError) {

	updateMap := map[string]interface{}{
		"full_name":         input.FullName,
		"login":             input.Phone,
		"password":          input.PasswordHash,
		"specialization_id": input.SpecializationID,
		"updated_at":        time.Now(),
	}

	updatedDoctorID, err := u.repo.UpdateDoctor(input.ID, updateMap)
	if err != nil {
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to update doctor", err, true)
	}
	updatedDoctor, err := u.repo.GetDoctorByID(updatedDoctorID)
	if err != nil {
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
	}

	return updatedDoctor, nil
}

func (u *DoctorUsecase) DeleteDoctor(id uint) *errors.AppError {
	err := u.repo.DeleteDoctor(id)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete doctor", err, true)
	}
	return nil
}
