package patient

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

// GetPatientsByGroup - получение пациентов по группе
func (r *PatientRepositoryImpl) GetPatientsByGroup(groupID uint) ([]entities.Patient, *errors.AppError) {
	op := "repo.Patient.GetPatientsByGroup"

	query := r.db.
		Preload("HarmPoint").
		Preload("PersonalInfo").
		Preload("ContactInfo").
		Preload("Flg").
		Preload("AnalysisOrder.OrderItems.Analysis").
		Preload("Vaccines").
		Preload("VaccineRefusals").
		Preload("VaccineWithdrawals").
		Preload("Titers").
		Preload("Receptions.Specialization").
		Preload("Receptions.Template").
		Preload("Specializations").
		Preload("Statistics").
		Where("patient_group_id = ?", groupID).
		Order("full_name")

	var patients []entities.Patient
	if err := query.Find(&patients).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return patients, nil
}

// getDB извлекает транзакцию из контекста или возвращает основное подключение
func (r *PatientRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return r.db
}

// CreateContactInfo создаёт контактную информацию
func (r *PatientRepositoryImpl) CreateContactInfo(ctx context.Context, contactInfo *entities.ContactInfo) error {
	op := "repo.Patient.CreateContactInfo"
	if err := r.getDB(ctx).WithContext(ctx).Create(contactInfo).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreatePersonalInfo создаёт персональную информацию
func (r *PatientRepositoryImpl) CreatePersonalInfo(ctx context.Context, personalInfo *entities.PersonalInfo) error {
	op := "repo.Patient.CreatePersonalInfo"
	if err := r.getDB(ctx).WithContext(ctx).Create(personalInfo).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateAnalysisOrder создаёт направление на анализы
func (r *PatientRepositoryImpl) CreateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error {
	op := "repo.Patient.CreateAnalysisOrder"
	if err := r.getDB(ctx).WithContext(ctx).Create(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// UpdateAnalysisOrder обновляет направление
func (r *PatientRepositoryImpl) UpdateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error {
	op := "repo.Patient.UpdateAnalysisOrder"
	if err := r.getDB(ctx).WithContext(ctx).Save(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreatePatient создаёт пациента
func (r *PatientRepositoryImpl) CreatePatient(ctx context.Context, patient *entities.Patient) error {
	op := "repo.Patient.CreatePatient"
	if err := r.getDB(ctx).WithContext(ctx).Create(patient).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CacheSpecializations связывает пациента со специализациями
func (r *PatientRepositoryImpl) CacheSpecializations(ctx context.Context, patient *entities.Patient, specializations []entities.Specialization) error {
	op := "repo.Patient.CacheSpecializations"
	if err := r.getDB(ctx).WithContext(ctx).Model(patient).Association("Specializations").Append(&specializations); err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateReceptions создаёт приёмы
func (r *PatientRepositoryImpl) CreateReceptions(ctx context.Context, receptions []entities.Reception) error {
	if len(receptions) == 0 {
		return nil
	}
	op := "repo.Patient.CreateReceptions"
	if err := r.getDB(ctx).WithContext(ctx).Create(&receptions).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreatePatientStatistics создаёт статистику пациента
func (r *PatientRepositoryImpl) CreatePatientStatistics(ctx context.Context, stats *entities.PatientStatistics) error {
	op := "repo.Patient.CreatePatientStatistics"
	if err := r.getDB(ctx).WithContext(ctx).Create(stats).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// GetReceptionTemplatesByHarmPointID возвращает шаблоны заключений по HarmPointID
func (r *PatientRepositoryImpl) GetReceptionTemplatesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.ReceptionTemplate, error) {
	op := "repo.Patient.GetReceptionTemplatesByHarmPointID"

	var harmPoint entities.HarmPoint
	if err := r.getDB(ctx).WithContext(ctx).Preload("ReceptionTemplates").First(&harmPoint, harmPointID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return harmPoint.ReceptionTemplates, nil
}

// PreloadPatientWithSpecializations загружает пациента со специализациями
func (r *PatientRepositoryImpl) PreloadPatientWithSpecializations(ctx context.Context, patientID uint) (*entities.Patient, error) {
	op := "repo.Patient.PreloadPatientWithSpecializations"
	var patient entities.Patient
	if err := r.getDB(ctx).WithContext(ctx).Preload("Specializations").First(&patient, patientID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return &patient, nil
}

func (r *PatientRepositoryImpl) CreateAnalysisItems(ctx context.Context, items []entities.AnalysisOrderItem) error {
	op := "repo.Patient.CreateAnalysisItems"
	if err := r.getDB(ctx).WithContext(ctx).Create(&items).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *PatientRepositoryImpl) GetMandatoryReceptionTemplateCodes(ctx context.Context) ([]string, error) {
	var codes []string
	if err := r.db.
		Model(&entities.Manual{}).
		Where("type = ?", entities.RefTypeMandatoryReception).
		Pluck("value", &codes).Error; err != nil {
		return nil, errors.NewDBError("repo.GetMandatoryReceptionTemplateCodes", err)
	}
	return codes, nil
}

func (r *PatientRepositoryImpl) GetMandatoryAnalysisCodes(ctx context.Context) ([]string, error) {
	var codes []string
	if err := r.db.
		Model(&entities.Manual{}).
		Where("type = ?", entities.RefTypeMandatoryAnalysis).
		Pluck("value", &codes).Error; err != nil {
		return nil, errors.NewDBError("repo.GetMandatoryAnalysisCodes", err)
	}
	return codes, nil
}

// PatientRepositoryImpl
func (r *PatientRepositoryImpl) GetReceptionTemplatesByCodes(ctx context.Context, codes []string) ([]entities.ReceptionTemplate, error) {
	if len(codes) == 0 {
		return []entities.ReceptionTemplate{}, nil
	}
	var templates []entities.ReceptionTemplate
	if err := r.getDB(ctx).WithContext(ctx).Where("code IN ?", codes).Find(&templates).Error; err != nil {
		return nil, errors.NewDBError("repo.GetReceptionTemplatesByCodes", err)
	}
	return templates, nil
}

func (r *PatientRepositoryImpl) GetAnalysesByCodes(ctx context.Context, codes []string) ([]entities.Analysis, error) {
	if len(codes) == 0 {
		return []entities.Analysis{}, nil
	}
	var analyses []entities.Analysis
	if err := r.getDB(ctx).WithContext(ctx).Where("code IN ?", codes).Find(&analyses).Error; err != nil {
		return nil, errors.NewDBError("repo.GetAnalysesByCodes", err)
	}
	return analyses, nil
}
