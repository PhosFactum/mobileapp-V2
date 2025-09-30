package usecases

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type PatientUsecase struct {
	repo          interfaces.PatientRepository
	manual        interfaces.ManualRepository
	contactRepo   interfaces.ContactInfoRepository
	personalRepo  interfaces.PersonalInfoRepository
	FilterBuilder interfaces.FilterBuilderService
}

func NewPatientUsecase(repo interfaces.PatientRepository, contactRepo interfaces.ContactInfoRepository, personalRepo interfaces.PersonalInfoRepository, s interfaces.Service) interfaces.PatientUsecase {
	return &PatientUsecase{
		repo:          repo,
		contactRepo:   contactRepo,
		personalRepo:  personalRepo,
		FilterBuilder: s}
}

// CreatePatient - создание пациента
func (u *PatientUsecase) CreatePatient(patientData *models.CreatePatientData, group_id uint) (*entities.Patient, *errors.AppError) {
	// // 1. Валидация входных данных
	// if err := u.validateCreatePatientData(patientData); err != nil {
	//     return nil, errors.NewValidationError(op, err)
	// }

	// // 2. Проверяем существование обязательных сущностей
	// if err := u.validateRequiredEntities(patientData); err != nil {
	//     return nil, errors.NewValidationError(op, err)
	// }

	// 3. Создаем пациента через репозиторий
	patient, err := u.repo.CreatePatient(patientData, group_id)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to create patient", err, true)
	}

	return patient, nil
}

// GetPatientsByGroup — основной метод
func (u *PatientUsecase) GetPatientsByGroup(groupID uint) ([]models.PatientResponse, *errors.AppError) {
	patients, err := u.repo.GetPatientsByGroup(groupID)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get patients",
			err,
			true,
		)
	}

	// Преобразуем сущности → DTO
	var response []models.PatientResponse
	for _, p := range patients {
		resp, err := u.buildPatientResponse(p)
		if err != nil {
			return nil, err
		}
		response = append(response, resp)
	}

	return response, nil
}

func (u *PatientUsecase) buildPatientResponse(p entities.Patient) (models.PatientResponse, *errors.AppError) {
	examType, err := u.mapExaminationType(p.ExaminationTypeID)
	if err != nil {
		return models.PatientResponse{}, err
	}
	examView, err := u.mapExaminationView(p.ExaminationViewID)
	if err != nil {
		return models.PatientResponse{}, err
	}
	vaccines, err := u.mapVaccines(p.Vaccines, p.VaccineRefusals, p.VaccineWithdrawals, p.Titers)
	if err != nil {
		return models.PatientResponse{}, err
	}

	personalInfo, err := u.mapPersonalInfo(p.PersonalInfo)
	if err != nil {
		return models.PatientResponse{}, err
	}
	return models.PatientResponse{
		ID:             p.ID,
		FullName:       p.FullName,
		BirthDate:      p.BirthDate,
		Age:            u.calculateAge(p.BirthDate),
		IsMale:         p.IsMale,
		Position:       p.Position,
		Division:       p.Division,
		PatientGroupID: p.PatientGroupID,

		// ✅ Просто копируем строки
		ExaminationType: examType,
		ExaminationView: examView,

		Vaccines:        vaccines,
		Receptions:      u.mapReceptions(p.Receptions),
		AnalysisOrder:   u.mapAnalysisOrder(p.AnalysisOrder),
		HarmPoint:       u.mapHarmPoint(p.HarmPoint),
		PersonalInfo:    personalInfo,
		ContactInfo:     u.mapContactInfo(p.ContactInfo),
		Flg:             u.mapFlg(p.Flg),
		Statistics:      u.mapStatistics(p.Statistics),
		Specializations: u.mapSpecializations(p.Specializations),
	}, nil
}

func (u *PatientUsecase) mapHarmPoint(point *entities.HarmPoint) *models.HarmPointResponse {
	return &models.HarmPointResponse{
		ID:    point.ID,
		Value: point.Value,
	}
}

func (u *PatientUsecase) mapExaminationType(id uint) (string, *errors.AppError) {
	if id == 0 {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"NotFound ExamType",
			errors.ErrEmptyData,
			true,
		)
	}

	val, err := u.manual.GetManualValueByTypeAndID(id, entities.RefTypePatientExaminationType)
	if err != nil {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get ExaminationType",
			err,
			false,
		)
	}

	return val, nil
}

func (u *PatientUsecase) mapExaminationView(id uint) (string, *errors.AppError) {
	if id == 0 {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"NotFound ExamView",
			errors.ErrEmptyData,
			true,
		)
	}

	val, err := u.manual.GetManualValueByTypeAndID(id, entities.RefTypePatientExaminationView)
	if err != nil {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get ExaminationType",
			err,
			false,
		)
	}

	return val, nil
}

func (u *PatientUsecase) mapVaccineTitle(id uint) (string, *errors.AppError) {
	if id == 0 {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"NotFound DocType",
			errors.ErrEmptyData,
			true,
		)
	}

	val, err := u.manual.GetManualValueByTypeAndID(id, entities.RefTypeVaccineTitle)
	if err != nil {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get ExaminationType",
			err,
			false,
		)
	}

	return val, nil
}

// mapVaccineToResponse маппит вакцинацию
func (u *PatientUsecase) mapVaccineToResponse(v entities.Vaccine) (models.VaccineAllResponse, *errors.AppError) {
	title, err := u.mapVaccineTitle(v.ID)
	if err != nil {
		return models.VaccineAllResponse{}, err
	}
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "vaccination",
		Title:          title,
		TiterAmountStr: nil,
	}, nil
}

// mapVaccineRefusalToResponse маппит отказ
func (u *PatientUsecase) mapVaccineRefusalToResponse(v entities.VaccineRefusal) (models.VaccineAllResponse, *errors.AppError) {
	title, err := u.mapVaccineTitle(v.ID)
	if err != nil {
		return models.VaccineAllResponse{}, err
	}
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "vaccination",
		Title:          title,
		TiterAmountStr: nil,
	}, nil
}

// mapVaccineWithdrawalToResponse маппит отвод
func (u *PatientUsecase) mapVaccineWithdrawalToResponse(v entities.VaccineWithdrawal) (models.VaccineAllResponse, *errors.AppError) {
	title, err := u.mapVaccineTitle(v.ID)
	if err != nil {
		return models.VaccineAllResponse{}, err
	}
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "vaccination",
		Title:          title,
		TiterAmountStr: nil,
	}, nil
}

// mapTitrToResponse маппит титрование
func (u *PatientUsecase) mapTitrToResponse(v entities.Titr) (models.VaccineAllResponse, *errors.AppError) {
	title, err := u.mapVaccineTitle(v.ID)
	if err != nil {
		return models.VaccineAllResponse{}, err
	}
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "vaccination",
		Title:          title,
		TiterAmountStr: nil,
	}, nil
}

func (u *PatientUsecase) mapVaccines(
	vaccines []entities.Vaccine,
	refusals []entities.VaccineRefusal,
	withdrawals []entities.VaccineWithdrawal,
	titers []entities.Titr,
) ([]models.VaccineAllResponse, *errors.AppError) {

	var result []models.VaccineAllResponse

	// Маппинг вакцинаций
	for _, v := range vaccines {
		resp, err := u.mapVaccineToResponse(v)
		if err != nil {
			return nil, err
		}
		result = append(result, resp)
	}

	// Маппинг отказов
	for _, r := range refusals {
		resp, err := u.mapVaccineRefusalToResponse(r)
		if err != nil {
			return nil, err
		}
		result = append(result, resp)
	}

	// Маппинг отводов
	for _, w := range withdrawals {
		resp, err := u.mapVaccineWithdrawalToResponse(w)
		if err != nil {
			return nil, err
		}
		result = append(result, resp)
	}

	// Маппинг титров
	for _, t := range titers {
		resp, err := u.mapTitrToResponse(t)
		if err != nil {
			return nil, err
		}
		result = append(result, resp)
	}

	return result, nil
}

func (u *PatientUsecase) mapReceptions(receptions []entities.Reception) []models.ReceptionResponse {
	if receptions == nil {
		return nil
	}
	var result []models.ReceptionResponse
	for _, r := range receptions {
		result = append(result, models.ReceptionResponse{
			ID:               r.ID,
			IsCompleted:      r.IsCompleted,
			SpecializationID: r.SpecializationID,
			Specialization:   u.mapSpecialization(r.Specialization),
			Data:             r.Data,
		})
	}
	return result
}

func (u *PatientUsecase) mapDocumentType(id uint) (string, *errors.AppError) {
	if id == 0 {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"NotFound DocType",
			errors.ErrEmptyData,
			true,
		)
	}

	val, err := u.manual.GetManualValueByTypeAndID(id, entities.RefTypePersonalDocumentType)
	if err != nil {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get ExaminationType",
			err,
			false,
		)
	}

	return val, nil
}

func (u *PatientUsecase) mapPersonalInfo(pi *entities.PersonalInfo) (*models.PersonalInfoResponse, *errors.AppError) {
	if pi == nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"NotFound PersonalInfo",
			errors.ErrEmptyData,
			true,
		)
	}
	docType, err := u.mapDocumentType(pi.ID)
	if err != nil {
		return nil, err
	}
	return &models.PersonalInfoResponse{
		ID:           pi.ID,
		DocNumber:    pi.DocNumber,
		DocSeries:    pi.DocSeries,
		SNILS:        pi.SNILS,
		OMS:          pi.OMS,
		DocumentType: docType,
	}, nil
}

func (u *PatientUsecase) mapContactInfo(ci *entities.ContactInfo) *models.ContactInfoResponse {
	if ci == nil {
		return nil
	}
	return &models.ContactInfoResponse{
		ID:      ci.ID,
		Phone:   ci.Phone,
		Email:   ci.Email,
		Address: ci.Address,
	}
}

func (u *PatientUsecase) mapFlg(flg *entities.Flg) *models.FlgResponse {
	if flg == nil {
		return nil
	}
	return &models.FlgResponse{
		ID:           flg.ID,
		IsCompleted:  flg.IsCompleted,
		Organization: flg.Organization,
		Number:       flg.Number,
		Result:       flg.Result,
	}
}

func (u *PatientUsecase) mapAnalysis(a *entities.Analysis) *models.AnalysisResponse {
	if a == nil {
		return nil
	}
	return &models.AnalysisResponse{
		ID:    a.ID,
		Code:  a.Code,
		Title: a.Title,
		Price: a.Price,
	}
}

func (u *PatientUsecase) mapAnalysisOrderItems(items []entities.AnalysisOrderItem) []models.AnalysisOrderItemResponse {
	if items == nil {
		return nil
	}
	var result []models.AnalysisOrderItemResponse
	for _, item := range items {
		result = append(result, models.AnalysisOrderItemResponse{
			ID:          item.ID,
			AnalysisID:  item.AnalysisID,
			Analysis:    u.mapAnalysis(item.Analysis),
			IsCompleted: item.IsCompleted,
		})
	}
	return result
}

func (u *PatientUsecase) mapAnalysisOrder(ao *entities.AnalysisOrder) *models.AnalysisOrderResponse {
	if ao == nil {
		return nil
	}
	// Вычисляем сумму
	total := uint(0)
	for _, item := range ao.OrderItems {
		total += item.PriceAtAssignment
	}
	return &models.AnalysisOrderResponse{
		ID:          ao.ID,
		OrderNumber: ao.OrderNumber,
		TotalAmount: total,
		OrderItems:  u.mapAnalysisOrderItems(ao.OrderItems),
	}
}

func (u *PatientUsecase) mapStatistics(stat *entities.PatientStatistics) *models.PatientStatisticsResponse {
	if stat == nil {
		return nil
	}
	return &models.PatientStatisticsResponse{
		ID:                     stat.ID,
		TotalReceptions:        stat.TotalReceptions,
		CompletedReceptions:    stat.CompletedReceptions,
		TotalAnalysisOrders:    stat.TotalAnalysisOrders,
		CompletedAnalysisItems: stat.CompletedAnalysisItems,
	}
}

func (u *PatientUsecase) mapSpecializations(specs []entities.Specialization) []models.SpecializationResponse {
	if specs == nil {
		return nil
	}
	var result []models.SpecializationResponse
	for _, s := range specs {
		result = append(result, models.SpecializationResponse{
			ID:    s.ID,
			Title: s.Title,
		})
	}
	return result
}

func (u *PatientUsecase) mapSpecialization(spec *entities.Specialization) *models.SpecializationResponse {
	if spec == nil {
		return nil
	}
	return &models.SpecializationResponse{
		ID:    spec.ID,
		Title: spec.Title,
	}
}

func (u *PatientUsecase) calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	birthdayThisYear := time.Date(now.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.Local)
	if now.Before(birthdayThisYear) {
		age--
	}
	return age
}
