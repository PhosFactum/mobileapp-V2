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
		response = append(response, u.buildPatientResponse(p))
	}

	return response, nil
}

// buildPatientResponse — собирает модель пациента для API
func (u *PatientUsecase) buildPatientResponse(p entities.Patient) models.PatientResponse {
	return models.PatientResponse{
		ID:             p.ID,
		FullName:       p.FullName,
		BirthDate:      p.BirthDate,
		Age:            u.calculateAge(p.BirthDate),
		IsMale:         p.IsMale,
		Position:       p.Position,
		Division:       p.Division,
		PatientGroupID: p.PatientGroupID,

		ExaminationType: u.mapExaminationType(p.ExaminationType),
		ExaminationView: u.mapExaminationView(p.ExaminationView),
		HarmPoint:       u.mapHarmPoint(p.HarmPoint),
		PersonalInfo:    u.mapPersonalInfo(p.PersonalInfo),
		ContactInfo:     u.mapContactInfo(p.ContactInfo),
		Flg:             u.mapFlg(p.Flg),
		AnalysisOrder:   u.mapAnalysisOrder(p.AnalysisOrder),
		Statistics:      u.mapStatistics(p.Statistics),
		Vaccines:        u.mapVaccines(p.Vaccines),
		Receptions:      u.mapReceptions(p.Receptions),
		Specializations: u.mapSpecializations(p.Specializations),
	}
}

// calculateAge — вычисляет возраст на основе даты рождения
func (u *PatientUsecase) calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	birthdayThisYear := time.Date(now.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.Local)
	if now.Before(birthdayThisYear) {
		age--
	}
	return age
}

func (u *PatientUsecase) mapExaminationType(et *entities.ExaminationType) *models.ExaminationTypeResponse {
	if et == nil {
		return nil
	}
	return &models.ExaminationTypeResponse{
		ID:    et.ID,
		Value: et.Value,
	}
}

func (u *PatientUsecase) mapExaminationView(ev *entities.ExaminationView) *models.ExaminationViewResponse {
	if ev == nil {
		return nil
	}
	return &models.ExaminationViewResponse{
		ID:    ev.ID,
		Value: ev.Value,
	}
}

func (u *PatientUsecase) mapHarmPoint(hp *entities.HarmPoint) *models.HarmPointResponse {
	if hp == nil {
		return nil
	}
	return &models.HarmPointResponse{
		ID:    hp.ID,
		Value: hp.Value,
	}
}

func (u *PatientUsecase) mapPersonalInfo(pi *entities.PersonalInfo) *models.PersonalInfoResponse {
	if pi == nil {
		return nil
	}
	return &models.PersonalInfoResponse{
		ID:           pi.ID,
		DocNumber:    pi.DocNumber,
		DocSeries:    pi.DocSeries,
		SNILS:        pi.SNILS,
		OMS:          pi.OMS,
		DocumentType: u.mapDocumentType(pi.DocumentType),
	}
}

func (u *PatientUsecase) mapDocumentType(dt *entities.DocumentType) *models.DocumentTypeResponse {
	if dt == nil {
		return nil
	}
	return &models.DocumentTypeResponse{
		ID:    dt.ID,
		Value: dt.Value,
	}
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

func (u *PatientUsecase) mapVaccines(vaccines []entities.Vaccine) []models.VaccineResponse {
	if vaccines == nil {
		return nil
	}
	var result []models.VaccineResponse
	for _, v := range vaccines {
		result = append(result, models.VaccineResponse{
			ID:              v.ID,
			Date:            v.Date,
			IsCompleted:     v.IsCompleted,
			IsRefusal:       v.IsRefusal,
			IsExemption:     v.IsExemption,
			TiterAmount:     v.TiterAmount,
			MedWithdrawlNum: v.MedWithdrawlNum,
			Result:          v.Result,

			Title:             u.mapTitle(v.Title),
			Medication:        u.mapMedication(v.Medication),
			Dose:              u.mapDose(v.Dose),
			Number:            u.mapNumber(v.Number),
			CertificateNumber: u.mapCertificateNumber(v.CertificateNumber),
			BodyPart:          u.mapBodyPart(v.BodyPart),
			Method:            u.mapMethod(v.Method),
			Place:             u.mapPlace(v.Place),
		})
	}
	return result
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

func (u *PatientUsecase) mapAnalysisOrder(ao *entities.AnalysisOrder) *models.AnalysisOrderResponse {
	if ao == nil {
		return nil
	}
	return &models.AnalysisOrderResponse{
		ID:          ao.ID,
		OrderNumber: ao.OrderNumber,
		TotalAmount: ao.TotalAmount,
		OrderItems:  u.mapAnalysisOrderItems(ao.OrderItems),
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

func (u *PatientUsecase) mapAnalysis(a *entities.Analysis) *models.AnalysisResponse {
	if a == nil {
		return nil
	}
	return &models.AnalysisResponse{
		ID:    a.ID,
		Name:  a.Name,
		Price: a.Price,
	}
}

// Map-функции для вакцин
func (u *PatientUsecase) mapTitle(t *entities.Title) *models.TitleResponse {
	if t == nil {
		return nil
	}
	return &models.TitleResponse{
		ID:    t.ID,
		Value: t.Value,
	}
}

func (u *PatientUsecase) mapMedication(m *entities.Medication) *models.MedicationResponse {
	if m == nil {
		return nil
	}
	return &models.MedicationResponse{
		ID:    m.ID,
		Value: m.Value,
	}
}

func (u *PatientUsecase) mapDose(d *entities.Dose) *models.DoseResponse {
	if d == nil {
		return nil
	}
	return &models.DoseResponse{
		ID:    d.ID,
		Value: d.Value,
	}
}

func (u *PatientUsecase) mapNumber(n *entities.Number) *models.NumberResponse {
	if n == nil {
		return nil
	}
	return &models.NumberResponse{
		ID:    n.ID,
		Value: n.Value,
	}
}

func (u *PatientUsecase) mapCertificateNumber(cn *entities.CertificateNumber) *models.CertificateNumberResponse {
	if cn == nil {
		return nil
	}
	return &models.CertificateNumberResponse{
		ID:    cn.ID,
		Value: cn.Value,
	}
}

func (u *PatientUsecase) mapBodyPart(bp *entities.BodyPart) *models.BodyPartResponse {
	if bp == nil {
		return nil
	}
	return &models.BodyPartResponse{
		ID:    bp.ID,
		Value: bp.Value,
	}
}

func (u *PatientUsecase) mapMethod(m *entities.Method) *models.MethodResponse {
	if m == nil {
		return nil
	}
	return &models.MethodResponse{
		ID:    m.ID,
		Value: m.Value,
	}
}

func (u *PatientUsecase) mapPlace(p *entities.Place) *models.PlaceResponse {
	if p == nil {
		return nil
	}
	return &models.PlaceResponse{
		ID:    p.ID,
		Value: p.Value,
	}
}
