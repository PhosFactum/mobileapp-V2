package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type PatientUsecase struct {
	repo          interfaces.PatientRepository
	manualRepo    interfaces.ManualRepository
	txManager     interfaces.TxManager
	receptionRepo interfaces.ReceptionRepository
	analysisRepo  interfaces.AnalysisRepository
}

func NewPatientUsecase(repo interfaces.PatientRepository, manualRepo interfaces.ManualRepository, txManager interfaces.TxManager) interfaces.PatientUsecase {
	return &PatientUsecase{
		repo:       repo,
		manualRepo: manualRepo,
		txManager:  txManager,
	}
}

// CreatePatient создаёт нового пациента со всеми связанными сущностями
func (u *PatientUsecase) CreatePatient(ctx context.Context, req *models.CreatePatientRequest, groupID uint) (*entities.Patient, *errors.AppError) {
	op := "usecase.Patient.CreatePatient"

	// Начинаем транзакцию
	ctx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Отложенный откат/коммит
	defer func() {
		if err != nil {
			u.txManager.Rollback(ctx)
		}
	}()

	// ✅ 1. Создаём ContactInfo
	contactInfo := &entities.ContactInfo{
		Phone:     req.ContactInfo.Phone,
		Email:     req.ContactInfo.Email,
		Address:   req.ContactInfo.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err = u.repo.CreateContactInfo(ctx, contactInfo); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 2. Создаём PersonalInfo
	personalInfo := &entities.PersonalInfo{
		DocNumber:      req.PersonalInfo.DocNumber,
		DocSeries:      req.PersonalInfo.DocSeries,
		SNILS:          req.PersonalInfo.SNILS,
		OMS:            req.PersonalInfo.OMS,
		DocumentTypeID: req.PersonalInfo.DocumentTypeID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err = u.repo.CreatePersonalInfo(ctx, personalInfo); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 3. Создаём AnalysisOrder (с временным номером)
	analysisOrder := &entities.AnalysisOrder{
		OrderNumber: "TEMP",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err = u.analysisRepo.CreateAnalysisOrder(ctx, analysisOrder); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Обновляем номер на основе ID
	analysisOrder.OrderNumber = fmt.Sprintf("ORD-%06d", analysisOrder.ID)
	if err = u.analysisRepo.UpdateAnalysisOrder(ctx, analysisOrder); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 4. Получаем шаблоны заключений по HarmPointID
	templatesFromHarmPoint, err := u.receptionRepo.GetReceptionTemplatesByHarmPointID(ctx, req.HarmPointID)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ Получаем обязательные шаблоны по кодам
	mandatoryTemplateCodes, err := u.manualRepo.GetManualValuesByType(ctx, entities.RefTypeMandatoryReception)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}
	mandatoryTemplates, err := u.receptionRepo.GetReceptionTemplatesByCodes(ctx, mandatoryTemplateCodes)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Объединяем шаблоны
	templates := append(templatesFromHarmPoint, mandatoryTemplates...)

	// ✅ 5. Получаем анализы по HarmPointID
	analysesFromHarmPoint, err := u.analysisRepo.GetAnalysesByHarmPointID(ctx, req.HarmPointID)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ Получаем обязательные анализы по кодам
	mandatoryAnalysisCodes, err := u.manualRepo.GetManualValuesByType(ctx, entities.RefTypeMandatoryAnalysis)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}
	mandatoryAnalyses, err := u.analysisRepo.GetAnalysesByCodes(ctx, mandatoryAnalysisCodes)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ Объединяем анализы, исключая дубликаты по ID
	analysisMap := make(map[uint]entities.Analysis)
	for _, a := range analysesFromHarmPoint {
		analysisMap[a.ID] = a
	}
	for _, a := range mandatoryAnalyses {
		analysisMap[a.ID] = a
	}

	var allAnalyses []entities.Analysis
	for _, a := range analysisMap {
		allAnalyses = append(allAnalyses, a)
	}

	// ✅ 6. Создаём пациента
	patient := &entities.Patient{
		FullName:          req.FullName,
		BirthDate:         req.BirthDate,
		IsMale:            req.IsMale,
		Position:          req.Position,
		Division:          req.Division,
		ExaminationTypeID: req.ExaminationTypeID,
		ExaminationViewID: req.ExaminationViewID,
		HarmPointID:       req.HarmPointID,
		PatientGroupID:    groupID,
		PersonalInfoID:    personalInfo.ID,
		ContactInfoID:     contactInfo.ID,
		AnalysisOrderID:   analysisOrder.ID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err = u.repo.CreatePatient(ctx, patient); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 7. Создаём элементы анализа
	var analysisItems []entities.AnalysisOrderItem
	for _, analysis := range allAnalyses {
		analysisItems = append(analysisItems, entities.AnalysisOrderItem{
			OrderID:     analysisOrder.ID,
			AnalysisID:  analysis.ID,
			IsCompleted: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}
	if len(analysisItems) > 0 {
		if err = u.analysisRepo.CreateAnalysisItems(ctx, analysisItems); err != nil {
			return nil, errors.NewDBError(op, err)
		}
	}

	// ✅ 8. Привязываем заказ к пациенту
	analysisOrder.PatientID = patient.ID
	if err = u.analysisRepo.UpdateAnalysisOrder(ctx, analysisOrder); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 9. Кэшируем специализации
	specializationMap := make(map[uint]entities.Specialization)
	for _, tmpl := range templates {
		specializationMap[tmpl.SpecializationID] = entities.Specialization{ID: tmpl.SpecializationID}
	}
	var specializations []entities.Specialization
	for _, s := range specializationMap {
		specializations = append(specializations, s)
	}
	if len(specializations) > 0 {
		if err = u.repo.CacheSpecializations(ctx, patient, specializations); err != nil {
			return nil, errors.NewDBError(op, err)
		}
	}

	// ✅ 10. Создаём пустые приёмы по шаблонам
	var receptions []entities.Reception
	initialData := []byte(`{"values": {}}`)
	for _, tmpl := range templates {
		receptions = append(receptions, entities.Reception{
			PatientID:        patient.ID,
			SpecializationID: tmpl.SpecializationID,
			TemplateID:       tmpl.ID,
			IsCompleted:      false,
			Data:             initialData,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		})
	}
	if err = u.receptionRepo.CreateReceptions(ctx, receptions); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 11. Создаём статистику
	statistics := &entities.PatientStatistics{
		PatientID:               patient.ID,
		TotalReceptions:         int64(len(templates)),
		CompletedReceptions:     0,
		TotalAnalysisOrderItems: int64(len(allAnalyses)),
		CompletedAnalysisItems:  0,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}
	if err = u.repo.CreatePatientStatistics(ctx, statistics); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// ✅ 12. Предзагружаем пациента со специализациями для возврата
	createdPatient, err := u.repo.PreloadPatientWithSpecializations(ctx, patient.ID)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Коммитим транзакцию
	if err = u.txManager.Commit(ctx); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return createdPatient, nil
}

// // GetPatientsByGroup — без транзакции
// func (u *PatientUsecase) GetPatientsByGroup(ctx context.Context, groupID uint) ([]models.PatientResponse, *errors.AppError) {
// 	patients, err := u.repo.GetPatientsByGroup(ctx, groupID)
// 	if err != nil {
// 		return nil, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to get patients",
// 			err,
// 			true,
// 		)
// 	}

// 	var response []models.PatientResponse
// 	for _, p := range patients {
// 		resp, err := u.buildPatientResponse(p)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response = append(response, resp)
// 	}

// 	return response, nil
// }

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

func (u *PatientUsecase) mapHarmPoint(point *entities.HarmPoint) models.HarmPointResponse {
	return models.HarmPointResponse{
		ID:    point.ID,
		Value: point.Value,
	}
}

func (u *PatientUsecase) mapExaminationType(id uint) (string, *errors.AppError) {
	val, err := u.manualRepo.GetManualValueByTypeAndID(id, entities.RefTypePatientExaminationType)
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
	val, err := u.manualRepo.GetManualValueByTypeAndID(id, entities.RefTypePatientExaminationView)
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
	val, err := u.manualRepo.GetManualValueByTypeAndID(id, entities.RefTypeVaccineTitle)
	if err != nil {
		return "", errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get VaccineTitle",
			err,
			false,
		)
	}

	return val, nil
}

// mapVaccineToResponse маппит вакцинацию
func (u *PatientUsecase) mapVaccineToResponse(v entities.Vaccine) (models.VaccineAllResponse, *errors.AppError) {
	title, err := u.mapVaccineTitle(v.TitleID)
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
	title, err := u.mapVaccineTitle(v.TitleID)
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
	title, err := u.mapVaccineTitle(v.TitleID)
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
	title, err := u.mapVaccineTitle(v.TitleID)
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

// usecase/patient.go
func (u *PatientUsecase) mapReceptions(receptions []entities.Reception) []models.ReceptionResponse {
	if receptions == nil {
		return nil
	}

	result := make([]models.ReceptionResponse, len(receptions))
	for i, r := range receptions {
		result[i] = models.ReceptionResponse{
			ID:               r.ID,
			IsCompleted:      r.IsCompleted,
			SpecializationID: r.SpecializationID,
			Specialization:   u.mapSpecialization(r.Specialization),
			Template: models.ReceptionTemplateResponse{
				ID:     r.Template.ID,
				Code:   r.Template.Code,
				Fields: r.Template.Fields,
			},
			Data: r.Data,
		}
	}
	return result
}

func (u *PatientUsecase) mapDocumentType(id uint) (string, *errors.AppError) {
	val, err := u.manualRepo.GetManualValueByTypeAndID(id, entities.RefTypePersonalDocumentType)
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

func (u *PatientUsecase) mapPersonalInfo(pi *entities.PersonalInfo) (models.PersonalInfoResponse, *errors.AppError) {
	emty := models.PersonalInfoResponse{}
	if pi == nil {
		return emty, errors.NewAppError(
			errors.InternalServerErrorCode,
			"NotFound PersonalInfo",
			errors.ErrEmptyData,
			true,
		)
	}
	docType, err := u.mapDocumentType(pi.DocumentTypeID)
	if err != nil {
		return emty, err
	}
	return models.PersonalInfoResponse{
		ID:           pi.ID,
		DocNumber:    pi.DocNumber,
		DocSeries:    pi.DocSeries,
		SNILS:        pi.SNILS,
		OMS:          pi.OMS,
		DocumentType: docType,
	}, nil
}

func (u *PatientUsecase) mapContactInfo(ci *entities.ContactInfo) models.ContactInfoResponse {
	if ci == nil {
		return models.ContactInfoResponse{}
	}
	return models.ContactInfoResponse{
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
		Organization: flg.Organization,
		Number:       flg.Number,
		Result:       flg.Result,
	}
}

func (u *PatientUsecase) mapAnalysis(a *entities.Analysis) models.AnalysisResponse {
	if a == nil {
		return models.AnalysisResponse{}
	}
	return models.AnalysisResponse{
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

func (u *PatientUsecase) mapAnalysisOrder(ao *entities.AnalysisOrder) models.AnalysisOrderResponse {
	if ao == nil {
		return models.AnalysisOrderResponse{}
	}
	// Вычисляем сумму
	total := uint(0)
	for _, item := range ao.OrderItems {
		total += item.PriceAtAssignment
	}
	return models.AnalysisOrderResponse{
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
		TotalAnalysisOrders:    stat.TotalAnalysisOrderItems,
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

func (u *PatientUsecase) mapSpecialization(spec *entities.Specialization) models.SpecializationResponse {
	if spec == nil {
		return models.SpecializationResponse{}
	}
	return models.SpecializationResponse{
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
