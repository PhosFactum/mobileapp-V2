package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type PatientUsecase struct {
	repo              interfaces.PatientRepository
	manualRepo        interfaces.ManualRepository
	txManager         interfaces.TxManager
	receptionRepo     interfaces.ReceptionRepository
	analysisRepo      interfaces.AnalysisRepository
	analysisOrderRepo interfaces.AnalysisOrderRepository
	parser            interfaces.ParamsParserService
}

func NewPatientUsecase(repo interfaces.PatientRepository, manualRepo interfaces.ManualRepository, receptionRepo interfaces.ReceptionRepository, analysisRepo interfaces.AnalysisRepository, analysisOrderRepo interfaces.AnalysisOrderRepository, txManager interfaces.TxManager, parser interfaces.ParamsParserService) interfaces.PatientUsecase {
	return &PatientUsecase{
		repo:              repo,
		manualRepo:        manualRepo,
		receptionRepo:     receptionRepo,
		analysisRepo:      analysisRepo,
		analysisOrderRepo: analysisOrderRepo,
		txManager:         txManager,
		parser:            parser,
	}
}

// CreatePatient создаёт нового пациента со всеми связанными сущностями
func (u *PatientUsecase) CreatePatient(ctx context.Context, req models.CreatePatientRequest) (*models.PatientResponse, *errors.AppError) {
	op := "usecase.Patient.CreatePatient"

	// Начинаем транзакцию
	ctx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	shouldRollback := true
	defer func() {
		if shouldRollback {
			// Игнорируем ошибку Rollback — главное попытаться откатить
			_ = u.txManager.Rollback(ctx)
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
	if err = u.analysisOrderRepo.CreateAnalysisOrder(ctx, analysisOrder); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Обновляем номер на основе ID
	analysisOrder.OrderNumber = fmt.Sprintf("ORD-%06d", analysisOrder.ID)
	if err = u.analysisOrderRepo.UpdateAnalysisOrder(ctx, analysisOrder); err != nil {
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
		PatientGroupID:    req.GroupID,
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
		if err = u.analysisOrderRepo.CreateAnalysisItems(ctx, analysisItems); err != nil {
			return nil, errors.NewDBError(op, err)
		}
	}

	// ✅ 8. Привязываем заказ к пациенту
	analysisOrder.PatientID = patient.ID
	if err = u.analysisOrderRepo.UpdateAnalysisOrder(ctx, analysisOrder); err != nil {
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

	// Отключаем откат перед коммитом
	shouldRollback = false
	if err = u.txManager.Commit(ctx); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// 🔁 Преобразуем сущность в ответную модель
	resp, err := u.buildPatientResponse(*createdPatient)
	log.Printf("err is nil: %v", err == nil)
	log.Printf("err type: %T", err)
	if err != nil {
		log.Printf("err string: '%s'", err.Error())
		log.Printf("err detailed: '%+v'", err)
	}
	return &resp, nil
}

// GetPatientsByGroup — основной метод
func (u *PatientUsecase) GetPatientsByGroup(ctx context.Context, groupID uint) ([]models.PatientResponse, *errors.AppError) {
	patients, err := u.repo.GetPatientsByGroup(ctx, groupID)
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
	receptions, err := u.MapReceptions(p.Receptions)
	if err != nil {
		return models.PatientResponse{}, err
	}

	// if p.HarmPoint == nil {
	// 	return models.PatientResponse{}, errors.NewInternalError(
	// 		"buildPatientResponse",
	// 		fmt.Sprintf("harm_point not found for patient.id=%d, harm_point_id=%d", p.ID, p.HarmPointID),
	// 		nil,
	// 	)
	// }
	// if p.PersonalInfo == nil {
	// 	return models.PatientResponse{}, errors.NewInternalError(
	// 		"buildPatientResponse",
	// 		fmt.Sprintf("personal_info not found for patient.id=%d, personal_info_id=%d", p.ID, p.PersonalInfoID),
	// 		nil,
	// 	)
	// }
	// if p.ContactInfo == nil {
	// 	return models.PatientResponse{}, errors.NewInternalError(
	// 		"buildPatientResponse",
	// 		fmt.Sprintf("contact_info not found for patient.id=%d, contact_info_id=%d", p.ID, p.ContactInfoID),
	// 		nil,
	// 	)
	// }
	// if p.Statistics == nil {
	// 	return models.PatientResponse{}, errors.NewInternalError(
	// 		"buildPatientResponse",
	// 		fmt.Sprintf("statistics not found for patient.id=%d", p.ID),
	// 		nil,
	// 	)
	// }

	return models.PatientResponse{
		ID:                p.ID,
		FullName:          p.FullName,
		BirthDate:         p.BirthDate,
		Age:               u.calculateAge(p.BirthDate),
		IsMale:            p.IsMale,
		Position:          p.Position,
		Division:          p.Division,
		PatientGroupID:    p.PatientGroupID,
		ExaminationTypeID: p.ExaminationTypeID,
		ExaminationViewID: p.ExaminationViewID,

		Vaccines:        u.mapVaccines(p.Vaccines, p.VaccineRefusals, p.VaccineWithdrawals, p.Titers),
		Receptions:      receptions,
		AnalysisOrder:   u.mapAnalysisOrder(p.AnalysisOrder),
		HarmPoint:       u.mapHarmPoint(*p.HarmPoint),
		PersonalInfo:    u.mapPersonalInfo(*p.PersonalInfo),
		ContactInfo:     u.mapContactInfo(*p.ContactInfo),
		Flg:             u.mapFlg(p.Flg),
		Statistics:      u.mapStatistics(*p.Statistics),
		Specializations: u.mapSpecializations(p.Specializations),
	}, nil
}

func (u *PatientUsecase) mapHarmPoint(point entities.HarmPoint) models.HarmPointResponse {
	return models.HarmPointResponse{
		ID:    point.ID,
		Value: point.Value,
	}
}

// mapVaccineToResponse маппит вакцинацию
func (u *PatientUsecase) mapVaccineToResponse(v entities.Vaccine) models.VaccineAllResponse {
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "vaccination",
		TitleID:        v.TitleID,
		TiterAmountStr: nil,
	}
}

// mapVaccineRefusalToResponse маппит отказ
func (u *PatientUsecase) mapVaccineRefusalToResponse(v entities.VaccineRefusal) models.VaccineAllResponse {
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "refusal",
		TitleID:        v.TitleID,
		TiterAmountStr: nil,
	}
}

// mapVaccineWithdrawalToResponse маппит отвод
func (u *PatientUsecase) mapVaccineWithdrawalToResponse(v entities.VaccineWithdrawal) models.VaccineAllResponse {
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "withdrawal",
		TitleID:        v.TitleID,
		TiterAmountStr: nil,
	}
}

// mapTitrToResponse маппит титрование
func (u *PatientUsecase) mapTitrToResponse(v entities.Titr) models.VaccineAllResponse {
	amount := v.Amount
	return models.VaccineAllResponse{
		ID:             v.ID,
		Date:           v.Date,
		Type:           "titer",
		TitleID:        v.TitleID,
		TiterAmountStr: &amount,
	}
}

func (u *PatientUsecase) mapVaccines(
	vaccines []entities.Vaccine,
	refusals []entities.VaccineRefusal,
	withdrawals []entities.VaccineWithdrawal,
	titers []entities.Titr,
) []models.VaccineAllResponse {

	var result []models.VaccineAllResponse

	// Маппинг вакцинаций
	for _, v := range vaccines {
		resp := u.mapVaccineToResponse(v)

		result = append(result, resp)
	}

	// Маппинг отказов
	for _, r := range refusals {
		resp := u.mapVaccineRefusalToResponse(r)

		result = append(result, resp)
	}

	// Маппинг отводов
	for _, w := range withdrawals {
		resp := u.mapVaccineWithdrawalToResponse(w)

		result = append(result, resp)
	}

	// Маппинг титров
	for _, t := range titers {
		resp := u.mapTitrToResponse(t)

		result = append(result, resp)
	}

	return result
}

func (u *PatientUsecase) mapPersonalInfo(pi entities.PersonalInfo) models.PersonalInfoResponse {
	return models.PersonalInfoResponse{
		ID:             pi.ID,
		DocNumber:      pi.DocNumber,
		DocSeries:      pi.DocSeries,
		SNILS:          pi.SNILS,
		OMS:            pi.OMS,
		DocumentTypeID: pi.DocumentTypeID,
	}
}

func (u *PatientUsecase) mapContactInfo(ci entities.ContactInfo) models.ContactInfoResponse {
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
		Date:         u.parser.FormatDateToString(flg.Date),
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

func (u *PatientUsecase) mapStatistics(stat entities.PatientStatistics) *models.PatientStatisticsResponse {
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

// В пакете mapper
func MapSpecializations(specs []entities.Specialization) []models.SpecializationResponse {
	if specs == nil {
		return nil
	}
	result := make([]models.SpecializationResponse, len(specs))
	for i, s := range specs {
		result[i] = MapSpecialization(&s)
	}
	return result
}

func MapSpecialization(spec *entities.Specialization) models.SpecializationResponse {
	if spec == nil {
		return models.SpecializationResponse{}
	}
	return models.SpecializationResponse{
		ID:    spec.ID,
		Title: spec.Title,
	}
}

// MapReceptions преобразует список сущностей Reception в модели ответа.
func (u *PatientUsecase) MapReceptions(receptions []entities.Reception) ([]models.ReceptionResponse, *errors.AppError) {
	if receptions == nil {
		return nil, nil
	}

	result := make([]models.ReceptionResponse, len(receptions))
	for i, r := range receptions {
		if r.Template.ID == 0 {
			return nil, errors.NewInternalError(
				"MapReception",
				fmt.Sprintf("template not loaded for reception.id=%d", r.ID),
				nil,
			)
		}

		receptionResp, err := u.MapReception(r)
		if err != nil {
			return nil, err
		}
		result[i] = receptionResp
	}
	return result, nil
}

// MapReception преобразует одну сущность Reception в модель ответа.
func (u *PatientUsecase) MapReception(r entities.Reception) (models.ReceptionResponse, *errors.AppError) {
	fields, err := u.convertTemplateSchema(r.Template.Schema, r.Template.ID)
	if err != nil {
		return models.ReceptionResponse{}, err
	}

	return models.ReceptionResponse{
		ID:               r.ID,
		IsCompleted:      r.IsCompleted,
		SpecializationID: r.SpecializationID,
		Specialization:   MapSpecialization(r.Specialization), // ← нужно реализовать
		Template: models.ReceptionTemplateResponse{
			ID:     r.Template.ID,
			Code:   r.Template.Code,
			Fields: fields,
		},
		Data: r.Data,
	}, nil
}

// convertTemplateSchema — приватная вспомогательная функция
func (u *PatientUsecase) convertTemplateSchema(schema []byte, templateID uint) ([]models.FieldDescriptor, *errors.AppError) {
	fields, err := u.parser.ConvertJSONSchemaToFields(schema)
	if err != nil {
		return nil, errors.NewInternalError(
			"mapper.convertTemplateSchema",
			fmt.Sprintf("failed to convert schema for template ID %d", templateID),
			err,
		)
	}
	return fields, nil
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
