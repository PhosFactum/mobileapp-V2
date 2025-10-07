// repository/migrations.go

package repository

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	analysis "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/analysys"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/auth"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/consent_signatures"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/doctor"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/manual"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/organization"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/patient"
	patientgroup "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/patient_group"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/reception"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/vaccine"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	interfaces.AuthRepository
	interfaces.DoctorRepository
	interfaces.PatientRepository
	interfaces.ReceptionRepository
	interfaces.OrganizationRepository
	interfaces.PatientGroupRepository
	interfaces.ConsentSignatureRepository
	interfaces.ManualRepository
	interfaces.AnalysisRepository
	interfaces.VaccineRepository
}

func NewRepository(db *gorm.DB) (interfaces.Repository, error) {
	// Выполнение автомиграций
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("ошибка выполнения автомиграций: %w", err)

	}

	return &Repository{
		auth.NewAuthRepository(db),
		doctor.NewDoctorRepository(db),
		patient.NewPatientRepository(db),
		reception.NewReceptionRepository(db),
		organization.NewOrganizationRepository(db),
		patientgroup.NewPatientGroupRepository(db),
		consent_signatures.NewConsentSignatureRepository(db),
		manual.NewManualRepository(db),
		analysis.NewAnalysisRepository(db),
		vaccine.NewVaccineRepository(db),
	}, nil

}

func autoMigrate(db *gorm.DB) error {
	tablesToDelete := []string{
		// Зависимые от Patient
		"harm_point_analyses",
		"harm_point_reception_templates",
		"reception_templates",
		"receptions",
		"analysis_order_items",
		"vaccines",
		"vaccine_refusals",
		"vaccine_withdrawals",
		"titrs",
		"patient_statistics",

		// Many-to-many
		"patients_specializations",
		"doctor_specializations",
		"doctor_organizations",
		"harm_points_specializations",

		// Слабые сущности
		"fl_gs",
		"contact_infos",
		"personal_infos",
		"analysis_orders",

		// Основные
		"patients",
		"patient_groups",
		"doctors",
		"organizations",
		"specializations",
		"harm_points",
		"managers",

		// Независимые
		"analyses",
		"manuals",
	}

	for _, table := range tablesToDelete {
		if db.Migrator().HasTable(table) {
			if err := db.Migrator().DropTable(table); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}
	}

	models := []interface{}{
		// Единственный справочник
		&entities.Manual{},

		// Независимые
		&entities.Analysis{},
		&entities.HarmPoint{},
		&entities.Manager{},
		&entities.Specialization{},

		// Основные
		&entities.Doctor{},
		&entities.Organization{},
		&entities.PatientGroup{},

		// Слабые
		&entities.ContactInfo{},
		&entities.PersonalInfo{},
		&entities.Flg{},

		// Центральная сущность
		&entities.Patient{},

		// Зависимые от Patient
		&entities.AnalysisOrder{},
		&entities.AnalysisOrderItem{},
		&entities.Reception{},
		&entities.ReceptionTemplate{},
		&entities.Vaccine{},
		&entities.VaccineRefusal{},
		&entities.VaccineWithdrawal{},
		&entities.Titr{},
		&entities.PatientStatistics{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Индекс для пациентов
	if !db.Migrator().HasIndex(&entities.Patient{}, "idx_patients_group_name") {
		if err := db.Exec(`
			CREATE INDEX idx_patients_group_name ON patients (patient_group_id, full_name)
		`).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		fmt.Println("✅ Created index: idx_patients_group_name")
	}

	if err := seedTestData(db); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	return nil
}
func seedTestData(db *gorm.DB) error {
	// 1. Справочники
	if err := seedReferenceEntries(db); err != nil {
		return fmt.Errorf("failed to seed reference entries: %w", err)
	}

	// 2. Независимые сущности
	if err := seedAnalyses(db); err != nil {
		return fmt.Errorf("failed to seed analyses: %w", err)
	}
	if err := seedManagers(db); err != nil {
		return fmt.Errorf("failed to seed managers: %w", err)
	}

	// 3. Основные сущности
	if err := seedOrganizations(db); err != nil {
		return fmt.Errorf("failed to seed organizations: %w", err)
	}
	if err := seedPatientGroups(db); err != nil {
		return fmt.Errorf("failed to seed patient groups: %w", err)
	}
	if err := seedSpecializations(db); err != nil {
		return fmt.Errorf("failed to seed specializations: %w", err)
	}
	if err := seedHarmPoints(db); err != nil {
		return fmt.Errorf("failed to seed harm points: %w", err)
	}
	if err := seedReceptionTemplatesAndLinks(db); err != nil {
		return fmt.Errorf("failed to seed reception templates: %w", err)
	}
	if err := seedDoctors(db); err != nil {
		return fmt.Errorf("failed to seed doctors: %w", err)
	}

	// 4. Пациенты (включая ContactInfo, PersonalInfo, Flg, AnalysisOrder, PatientStatistics)
	if err := seedPatients(db); err != nil {
		return fmt.Errorf("failed to seed patients: %w", err)
	}

	// 5. Зависимые от Patient
	if err := seedVaccines(db); err != nil {
		return fmt.Errorf("failed to seed vaccines: %w", err)
	}
	if err := seedVaccineRefusals(db); err != nil {
		return fmt.Errorf("failed to seed vaccine refusals: %w", err)
	}
	if err := seedVaccineWithdrawals(db); err != nil {
		return fmt.Errorf("failed to seed vaccine withdrawals: %w", err)
	}
	if err := seedTitrs(db); err != nil {
		return fmt.Errorf("failed to seed titrs: %w", err)
	}
	if err := seedReceptions(db); err != nil {
		return fmt.Errorf("failed to seed receptions: %w", err)
	}

	fmt.Println("✅ All test data seeded successfully")
	return nil
}

func seedReferenceEntries(db *gorm.DB) error {
	entries := []entities.Manual{
		// Document types
		{Type: entities.RefTypePersonalDocumentType, Value: "Паспорт РФ"},
		{Type: entities.RefTypePersonalDocumentType, Value: "Свидетельство о рождении"},

		// Examination
		{Type: entities.RefTypePatientExaminationType, Value: "Предварительный"},
		{Type: entities.RefTypePatientExaminationType, Value: "Периодический"},
		{Type: entities.RefTypePatientExaminationView, Value: "Осмотр терапевта"},
		{Type: entities.RefTypePatientExaminationView, Value: "Осмотр невролога"},

		// Vaccine
		{Type: entities.RefTypeVaccineTitle, Value: "COVID-19"},
		{Type: entities.RefTypeVaccineTitle, Value: "Грипп"},
		{Type: entities.RefTypeVaccineMedication, Value: "Спутник V"},
		{Type: entities.RefTypeVaccineMedication, Value: "Совигрипп"},
		{Type: entities.RefTypeVaccineDose, Value: "0.5"},
		{Type: entities.RefTypeVaccineNumber, Value: "1"},
		{Type: entities.RefTypeVaccineNumber, Value: "2"},
		{Type: entities.RefTypeVaccineCertificateNumber, Value: "CERT-2024-001"},
		{Type: entities.RefTypeVaccineBodyPart, Value: "Левая рука"},
		{Type: entities.RefTypeVaccineMethod, Value: "Внутримышечно"},
		{Type: entities.RefTypeVaccinePlace, Value: "Поликлиника №1"},

		// Обязательные заключений
		{Type: entities.RefTypeMandatoryReception, Value: "THERAPY_ANAMNESIS_V1"},
		{Type: entities.RefTypeMandatoryReception, Value: "NEURO_EXAM_V1"},

		// Обязательные анализы
		{Type: entities.RefTypeMandatoryAnalysis, Value: "14-1231"},
		{Type: entities.RefTypeMandatoryAnalysis, Value: "15-4214"},
	}

	now := time.Now()
	for i := range entries {
		entries[i].CreatedAt = now
		entries[i].UpdatedAt = now
		if err := db.Create(&entries[i]).Error; err != nil {
			return fmt.Errorf("failed to create ref entry: %w", err)
		}
	}
	fmt.Println("✅ Seeded reference entries")
	return nil
}

func seedAnalyses(db *gorm.DB) error {
	analyses := []entities.Analysis{
		// Обязательные (уже есть в RefTypeMandatoryAnalysis)
		{Code: "14-1231", Title: "Общий анализ крови", Price: 500},
		{Code: "15-4214", Title: "ЭКГ", Price: 400},

		// Дополнительные (могут быть привязаны к HarmPoint)
		{Code: "10-5423", Title: "Флюорография", Price: 800},
		{Code: "11-9721", Title: "Общий анализ мочи", Price: 300},
		{Code: "20-8810", Title: "Биохимия крови (АЛТ, АСТ)", Price: 600},
		{Code: "21-3344", Title: "Глюкоза крови", Price: 250},
		{Code: "22-7766", Title: "Холестерин общий", Price: 300},
		{Code: "23-1122", Title: "Мазок на флору", Price: 450},
		{Code: "24-9900", Title: "ПЦР на ВИЧ", Price: 1200},
		{Code: "25-5566", Title: "УЗИ брюшной полости", Price: 1500},
	}

	for _, a := range analyses {
		// Проверим, не существует ли уже анализ с таким кодом
		var count int64
		db.Model(&entities.Analysis{}).Where("code = ?", a.Code).Count(&count)
		if count == 0 {
			if err := db.Create(&a).Error; err != nil {
				return fmt.Errorf("failed to create analysis %s: %w", a.Code, err)
			}
		}
	}
	fmt.Println("✅ Seeded analyses")
	return nil
}

func seedManagers(db *gorm.DB) error {
	managers := []entities.Manager{
		{FullName: "Иванов И.П.", Phone: "+79001111111"},
		{FullName: "Петров П.И.", Phone: "+79002222222"},
	}
	for i := range managers {
		db.Create(&managers[i])
	}
	fmt.Println("✅ Seeded managers")
	return nil
}

func seedOrganizations(db *gorm.DB) error {
	var managers []entities.Manager
	db.Find(&managers)
	orgs := []entities.Organization{
		{Title: "Клиника А", ManagerID: managers[0].ID},
		{Title: "Клиника Б", ManagerID: managers[1].ID},
	}
	for i := range orgs {
		db.Create(&orgs[i])
	}
	fmt.Println("✅ Seeded organizations")
	return nil
}

func seedPatientGroups(db *gorm.DB) error {
	var orgs []entities.Organization
	db.Find(&orgs)
	groups := []entities.PatientGroup{
		{Code: "GRP-001", OrganizationID: orgs[0].ID},
		{Code: "GRP-002", OrganizationID: orgs[1].ID},
	}
	for i := range groups {
		db.Create(&groups[i])
	}
	fmt.Println("✅ Seeded patient groups")
	return nil
}

func seedSpecializations(db *gorm.DB) error {
	specs := []entities.Specialization{
		{Title: "Терапевт"},
		{Title: "Невролог"},
		{Title: "Травматолог"},
		{Title: "Психиатр"},
	}
	for i := range specs {
		db.Create(&specs[i])
	}
	fmt.Println("✅ Seeded specializations")
	return nil
}

func seedHarmPoints(db *gorm.DB) error {
	var templates []entities.ReceptionTemplate
	var analyses []entities.Analysis

	db.Find(&templates)
	db.Find(&analyses)

	// Получаем обязательные коды
	var mandatoryReceptionCodes, mandatoryAnalysisCodes []string
	db.Model(&entities.Manual{}).
		Where("type = ?", entities.RefTypeMandatoryReception).
		Pluck("value", &mandatoryReceptionCodes)
	db.Model(&entities.Manual{}).
		Where("type = ?", entities.RefTypeMandatoryAnalysis).
		Pluck("value", &mandatoryAnalysisCodes)

	mandatoryReceptionSet := make(map[string]struct{})
	mandatoryAnalysisSet := make(map[string]struct{})
	for _, c := range mandatoryReceptionCodes {
		mandatoryReceptionSet[c] = struct{}{}
	}
	for _, c := range mandatoryAnalysisCodes {
		mandatoryAnalysisSet[c] = struct{}{}
	}

	// Фильтруем: только НЕ обязательные
	var harmPointTemplates []entities.ReceptionTemplate
	var harmPointAnalyses []entities.Analysis

	for _, t := range templates {
		if _, isMandatory := mandatoryReceptionSet[t.Code]; !isMandatory {
			harmPointTemplates = append(harmPointTemplates, t)
		}
	}
	for _, a := range analyses {
		if _, isMandatory := mandatoryAnalysisSet[a.Code]; !isMandatory {
			harmPointAnalyses = append(harmPointAnalyses, a)
		}
	}

	hp := []entities.HarmPoint{
		{Value: "3.1"},
		{Value: "3.2"},
		{Value: "3.3"},
	}

	for i := range hp {
		if err := db.Create(&hp[i]).Error; err != nil {
			return fmt.Errorf("create harm point: %w", err)
		}

		// 🔗 Привязываем ТОЛЬКО шаблоны и анализы
		if len(harmPointTemplates) > 0 {
			start := i % len(harmPointTemplates)
			end := start + 2
			if end > len(harmPointTemplates) {
				end = len(harmPointTemplates)
			}
			if err := db.Model(&hp[i]).Association("ReceptionTemplates").Append(harmPointTemplates[start:end]); err != nil {
				return fmt.Errorf("link reception templates: %w", err)
			}
		}

		if len(harmPointAnalyses) > 0 {
			start := i % len(harmPointAnalyses)
			end := start + 3
			if end > len(harmPointAnalyses) {
				end = len(harmPointAnalyses)
			}
			if err := db.Model(&hp[i]).Association("Analyses").Append(harmPointAnalyses[start:end]); err != nil {
				return fmt.Errorf("link analyses: %w", err)
			}
		}
	}

	fmt.Println("✅ Seeded harm points with templates and analyses (no specializations)")
	return nil
}

func seedReceptionTemplatesAndLinks(db *gorm.DB) error {
	var specializations []entities.Specialization
	if err := db.Find(&specializations).Error; err != nil {
		return fmt.Errorf("failed to load specializations: %w", err)
	}

	if len(specializations) == 0 {
		return fmt.Errorf("no specializations found")
	}

	// Соберём текущие обязательные коды из справочника
	var mandatoryReceptionCodes []string
	db.Model(&entities.Manual{}).
		Where("type = ?", entities.RefTypeMandatoryReception).
		Pluck("value", &mandatoryReceptionCodes)

	mandatorySet := make(map[string]struct{})
	for _, code := range mandatoryReceptionCodes {
		mandatorySet[code] = struct{}{}
	}

	// Генератор шаблонов
	createTemplate := func(specID uint, code string, fields []map[string]interface{}) error {
		fieldsJSON, _ := json.Marshal(fields)
		tmpl := entities.ReceptionTemplate{
			Code:             code,
			SpecializationID: specID,
			Fields:           json.RawMessage(fieldsJSON),
		}
		return db.FirstOrCreate(&tmpl, entities.ReceptionTemplate{Code: code}).Error
	}

	// === 1. Обязательные шаблоны (уже есть в RefTypeMandatoryReception) ===
	// Они НЕ должны быть привязаны к HarmPoint напрямую
	for _, spec := range specializations {
		switch spec.Title {
		case "Терапевт":
			if _, ok := mandatorySet["THERAPY_ANAMNESIS_V1"]; ok {
				createTemplate(spec.ID, "THERAPY_ANAMNESIS_V1", []map[string]interface{}{
					{"name": "complaints", "type": "string", "required": true},
					{"name": "bp_systolic", "type": "integer", "required": true, "min": 80, "max": 200},
					{"name": "bp_diastolic", "type": "integer", "required": true, "min": 50, "max": 120},
					{"name": "heart_rate", "type": "integer", "required": true, "min": 40, "max": 200},
					{"name": "temperature", "type": "number", "required": true, "min": 35.0, "max": 42.0},
					{"name": "diagnosis", "type": "string", "required": true},
				})
			}
			// Дополнительный обязательный (если нужен)
			if _, ok := mandatorySet["THERAPY_FOLLOWUP_V1"]; ok {
				createTemplate(spec.ID, "THERAPY_FOLLOWUP_V1", []map[string]interface{}{
					{"name": "previous_diagnosis", "type": "string", "required": true},
					{"name": "current_symptoms", "type": "string", "required": true},
					{"name": "medication_effect", "type": "string", "required": false},
					{"name": "new_diagnosis", "type": "string", "required": true},
				})
			}

		case "Невролог":
			if _, ok := mandatorySet["NEURO_EXAM_V1"]; ok {
				createTemplate(spec.ID, "NEURO_EXAM_V1", []map[string]interface{}{
					{"name": "mental_status", "type": "string", "required": true},
					{"name": "motor_function", "type": "string", "required": true},
					{"name": "sensory_function", "type": "string", "required": true},
					{"name": "reflexes", "type": "string", "required": true},
					{"name": "diagnosis", "type": "string", "required": true},
				})
			}

		case "Травматолог":
			if _, ok := mandatorySet["TRAUMA_INITIAL_V1"]; ok {
				createTemplate(spec.ID, "TRAUMA_INITIAL_V1", []map[string]interface{}{
					{"name": "injury_type", "type": "string", "required": true},
					{"name": "localization", "type": "string", "required": true},
					{"name": "xray_results", "type": "string", "required": false},
					{"name": "swelling", "type": "boolean", "required": true},
					{"name": "treatment_plan", "type": "string", "required": true},
				})
			}

		case "Психиатр":
			if _, ok := mandatorySet["PSYCH_INITIAL_V1"]; ok {
				createTemplate(spec.ID, "PSYCH_INITIAL_V1", []map[string]interface{}{
					{"name": "mood", "type": "string", "required": true},
					{"name": "sleep_quality", "type": "string", "required": true},
					{"name": "appetite", "type": "string", "required": true},
					{"name": "suicidal_ideation", "type": "boolean", "required": true},
					{"name": "diagnosis_icd", "type": "string", "required": true},
				})
			}
		}
	}

	// === 2. Дополнительные шаблоны (НЕ обязательные, могут быть привязаны к HarmPoint) ===
	// Их коды НЕ должны быть в mandatorySet
	extraTemplates := []struct {
		SpecTitle string
		Code      string
		Fields    []map[string]interface{}
	}{
		{"Терапевт", "THERAPY_CARDIO_V1", []map[string]interface{}{
			{"name": "ecg_result", "type": "string", "required": true},
			{"name": "cholesterol", "type": "number", "required": true},
			{"name": "cardio_diagnosis", "type": "string", "required": true},
		}},
		{"Терапевт", "THERAPY_PULMO_V1", []map[string]interface{}{
			{"name": "respiratory_rate", "type": "integer", "required": true},
			{"name": "o2_saturation", "type": "number", "required": true},
			{"name": "lung_exam", "type": "string", "required": true},
		}},
		{"Невролог", "NEURO_EEG_V1", []map[string]interface{}{
			{"name": "eeg_result", "type": "string", "required": true},
			{"name": "seizure_history", "type": "string", "required": true},
			{"name": "neuro_diagnosis", "type": "string", "required": true},
		}},
		{"Травматолог", "TRAUMA_FOLLOWUP_V1", []map[string]interface{}{
			{"name": "healing_progress", "type": "string", "required": true},
			{"name": "pain_level", "type": "integer", "required": true, "min": 0, "max": 10},
			{"name": "rehab_plan", "type": "string", "required": true},
		}},
	}

	for _, et := range extraTemplates {
		specID, exists := func() (uint, bool) {
			for _, s := range specializations {
				if s.Title == et.SpecTitle {
					return s.ID, true
				}
			}
			return 0, false
		}()
		if !exists {
			continue
		}
		// Убедимся, что код не в mandatory
		if _, isMandatory := mandatorySet[et.Code]; isMandatory {
			continue // пропускаем, если вдруг оказался обязательным
		}
		createTemplate(specID, et.Code, et.Fields)
	}

	fmt.Println("✅ Seeded reception templates (mandatory + extra)")
	return nil
}

func seedDoctors(db *gorm.DB) error {
	var specs []entities.Specialization
	var orgs []entities.Organization
	db.Find(&specs)
	db.Find(&orgs)

	doctors := []entities.Doctor{
		{FullName: "Смирнова А.М.", Phone: "+79161111111", PasswordHash: hashPassword("123")},
		{FullName: "Козлов В.П.", Phone: "+79162222222", PasswordHash: hashPassword("123")},
	}
	for i := range doctors {
		db.Create(&doctors[i])
		// Связываем со специализацией и организацией
		db.Model(&doctors[i]).Association("Specializations").Append(&specs[i%len(specs)])
		db.Model(&doctors[i]).Association("Organizations").Append(&orgs[i%len(orgs)])
	}
	fmt.Println("✅ Seeded doctors")
	return nil
}

func seedContactInfos(db *gorm.DB) ([]entities.ContactInfo, error) {
	contacts := []entities.ContactInfo{
		{Phone: "+79001111111", Email: "ivanov@example.com", Address: "г. Москва, ул. Ленина, д.1"},
		{Phone: "+79002222222", Email: "petrova@example.com", Address: "г. Москва, ул. Пушкина, д.2"},
	}
	for i := range contacts {
		contacts[i].CreatedAt = time.Now()
		contacts[i].UpdatedAt = time.Now()
		db.Create(&contacts[i])
	}
	fmt.Println("✅ Seeded contact infos")
	return contacts, nil
}

func seedPersonalInfos(db *gorm.DB) ([]entities.PersonalInfo, error) {
	var docTypes []entities.Manual
	db.Where("type = ?", entities.RefTypePersonalDocumentType).Find(&docTypes)

	personal := []entities.PersonalInfo{
		{DocNumber: "123456", DocSeries: "4510", SNILS: "123-456-789 00", OMS: "1234567890123456", DocumentTypeID: docTypes[0].ID},
		{DocNumber: "789012", DocSeries: "4511", SNILS: "987-654-321 00", OMS: "6543210987654321", DocumentTypeID: docTypes[1].ID},
	}
	for i := range personal {
		personal[i].CreatedAt = time.Now()
		personal[i].UpdatedAt = time.Now()
		db.Create(&personal[i])
	}
	fmt.Println("✅ Seeded personal infos")
	return personal, nil
}

func seedFlgs(db *gorm.DB) ([]*uint, error) {
	flgs := []entities.Flg{
		{Organization: "Stavropol", Number: "984212", Result: "COVID", Date: time.Now()},
		{Organization: "Moscow", Number: "984213", Result: "Negative", Date: time.Now()},
	}
	var ids []*uint
	for i := range flgs {
		db.Create(&flgs[i])
		id := flgs[i].ID
		ids = append(ids, &id)
	}
	fmt.Println("✅ Seeded FLGs")
	return ids, nil
}

func seedPatients(db *gorm.DB) error {
	var groups []entities.PatientGroup
	var harmPoints []entities.HarmPoint
	var examTypes []entities.Manual
	var examViews []entities.Manual
	db.Find(&groups)
	db.Find(&harmPoints)
	db.Where("type = ?", entities.RefTypePatientExaminationType).Find(&examTypes)
	db.Where("type = ?", entities.RefTypePatientExaminationView).Find(&examViews)

	contacts, _ := seedContactInfos(db)
	personal, _ := seedPersonalInfos(db)
	flgIDs, _ := seedFlgs(db)

	patients := []entities.Patient{
		{
			FullName:          "Иванов Иван Иванович",
			BirthDate:         time.Date(1980, 5, 15, 0, 0, 0, 0, time.UTC),
			IsMale:            true,
			Position:          "Программист",
			Division:          "IT",
			ExaminationTypeID: examTypes[0].ID,
			ExaminationViewID: examViews[0].ID,
			PatientGroupID:    groups[0].ID,
			HarmPointID:       harmPoints[0].ID,
			PersonalInfoID:    personal[0].ID,
			ContactInfoID:     contacts[0].ID,
			FlgID:             flgIDs[0],
		},
		{
			FullName:          "Петрова Мария Сергеевна",
			BirthDate:         time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC),
			IsMale:            false,
			Position:          "Дизайнер",
			Division:          "Дизайн",
			ExaminationTypeID: examTypes[1].ID,
			ExaminationViewID: examViews[1].ID,
			PatientGroupID:    groups[1].ID,
			HarmPointID:       harmPoints[1].ID,
			PersonalInfoID:    personal[1].ID,
			ContactInfoID:     contacts[1].ID,
			FlgID:             flgIDs[1],
		},
	}

	for i := range patients {

		// === 2. Создаём AnalysisOrder ===
		order := entities.AnalysisOrder{
			OrderNumber: fmt.Sprintf("ORD-%06d", i+1),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		db.Create(&order)

		// === 3. Получаем анализы для пациента ===
		var mandatoryAnalysisCodes []string
		db.Model(&entities.Manual{}).
			Where("type = ?", entities.RefTypeMandatoryAnalysis).
			Pluck("value", &mandatoryAnalysisCodes)

		var mandatoryAnalyses []entities.Analysis
		if len(mandatoryAnalysisCodes) > 0 {
			db.Where("code IN ?", mandatoryAnalysisCodes).Find(&mandatoryAnalyses)
		}

		var harmPoint entities.HarmPoint
		db.Preload("Analyses").First(&harmPoint, patients[i].HarmPointID)
		harmPointAnalyses := harmPoint.Analyses

		analysisMap := make(map[uint]entities.Analysis)
		for _, a := range mandatoryAnalyses {
			analysisMap[a.ID] = a
		}
		for _, a := range harmPointAnalyses {
			analysisMap[a.ID] = a
		}

		var allAnalyses []entities.Analysis
		for _, a := range analysisMap {
			allAnalyses = append(allAnalyses, a)
		}

		// === 4. Создаём AnalysisOrderItem ===
		var orderItems []entities.AnalysisOrderItem
		for _, analysis := range allAnalyses {
			orderItems = append(orderItems, entities.AnalysisOrderItem{
				OrderID:           order.ID,
				AnalysisID:        analysis.ID,
				PriceAtAssignment: analysis.Price,
				IsCompleted:       false,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			})
		}
		if len(orderItems) > 0 {
			db.Create(&orderItems)
		}

		// === 5. Создаём пациента (уже с AnalysisOrderID) ===
		patients[i].AnalysisOrderID = order.ID
		db.Create(&patients[i])

		// Обновляем заказ с PatientID
		db.Model(&order).Update("patient_id", patients[i].ID)

		// === 6. Получаем шаблоны для пациента (для специализаций) ===
		var harmPointForTemplates entities.HarmPoint
		db.Preload("ReceptionTemplates").First(&harmPointForTemplates, patients[i].HarmPointID)
		harmPointTemplates := harmPointForTemplates.ReceptionTemplates

		var mandatoryTemplateCodes []string
		db.Model(&entities.Manual{}).
			Where("type = ?", entities.RefTypeMandatoryReception).
			Pluck("value", &mandatoryTemplateCodes)

		var mandatoryTemplates []entities.ReceptionTemplate
		if len(mandatoryTemplateCodes) > 0 {
			db.Where("code IN ?", mandatoryTemplateCodes).Find(&mandatoryTemplates)
		}

		templateMap := make(map[uint]entities.ReceptionTemplate)
		for _, t := range harmPointTemplates {
			templateMap[t.ID] = t
		}
		for _, t := range mandatoryTemplates {
			templateMap[t.ID] = t
		}

		// Собираем уникальные специализации
		specializationMap := make(map[uint]entities.Specialization)
		for _, tmpl := range templateMap {
			specializationMap[tmpl.SpecializationID] = entities.Specialization{ID: tmpl.SpecializationID}
		}

		var specializations []entities.Specialization
		for _, s := range specializationMap {
			specializations = append(specializations, s)
		}

		if len(specializations) > 0 {
			db.Model(&patients[i]).Association("Specializations").Append(specializations)
		}

		// === 7. Создаём статистику с реальными числами ===
		stats := entities.PatientStatistics{
			PatientID:               patients[i].ID,
			TotalReceptions:         int64(len(templateMap)), // по количеству шаблонов
			CompletedReceptions:     0,
			TotalAnalysisOrderItems: int64(len(allAnalyses)), // по количеству анализов
			CompletedAnalysisItems:  0,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}
		db.Create(&stats)
	}

	fmt.Println("✅ Seeded patients with orders, items, specializations, and stats")
	return nil
}

func seedVaccines(db *gorm.DB) error {
	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return err
	}
	if len(patients) == 0 {
		return nil
	}

	// Получаем справочники
	var titles []entities.Manual
	var meds []entities.Manual
	var doses []entities.Manual
	var nums []entities.Manual
	var certs []entities.Manual
	var bodyParts []entities.Manual
	var methods []entities.Manual
	var places []entities.Manual
	var results []entities.Manual

	db.Where("type = ?", entities.RefTypeVaccineTitle).Find(&titles)
	db.Where("type = ?", entities.RefTypeVaccineMedication).Find(&meds)
	db.Where("type = ?", entities.RefTypeVaccineDose).Find(&doses)
	db.Where("type = ?", entities.RefTypeVaccineNumber).Find(&nums)
	db.Where("type = ?", entities.RefTypeVaccineCertificateNumber).Find(&certs)
	db.Where("type = ?", entities.RefTypeVaccineBodyPart).Find(&bodyParts)
	db.Where("type = ?", entities.RefTypeVaccineMethod).Find(&methods)
	db.Where("type = ?", entities.RefTypeVaccinePlace).Find(&places)
	db.Where("type = ?", entities.RefTypeVaccinePlace).Find(&results)

	for i, p := range patients {
		vaccine := entities.Vaccine{
			PatientID:           p.ID,
			Date:                time.Now().AddDate(0, -1, 0),
			ResultID:            results[i%len(results)].ID,
			TitleID:             titles[i%len(titles)].ID,
			MedicationID:        meds[i%len(meds)].ID,
			DoseID:              doses[i%len(doses)].ID,
			NumberID:            nums[i%len(nums)].ID,
			CertificateNumberID: certs[i%len(certs)].ID,
			BodyPartID:          bodyParts[i%len(bodyParts)].ID,
			MethodID:            methods[i%len(methods)].ID,
			PlaceID:             places[i%len(places)].ID,
		}
		db.Create(&vaccine)
	}
	fmt.Println("✅ Seeded vaccines")
	return nil
}

func seedVaccineRefusals(db *gorm.DB) error {
	var titles []entities.Manual
	db.Where("type = ?", entities.RefTypeVaccineTitle).Find(&titles)

	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients for vaccine refusals: %w", err)
	}
	if len(patients) == 0 {
		fmt.Println("⚠️ No patients found, skipping vaccine refusals seeding")
		return nil
	}

	for _, patient := range patients {
		// Создаём отказ только у ~30% пациентов
		refusalDate := time.Now().AddDate(0, 0, -rand.Intn(365)) // за последний год
		refusal := &entities.VaccineRefusal{
			PatientID: patient.ID,
			TitleID:   titles[0].ID,
			Date:      refusalDate,
			CreatedAt: time.Now(),
		}
		if err := db.Create(refusal).Error; err != nil {
			return fmt.Errorf("failed to create vaccine refusal for patient %d: %w", patient.ID, err)
		}
		fmt.Printf("✅ Created vaccine refusal for %s (date: %s)\n", patient.FullName, refusalDate.Format("2006-01-02"))
	}

	fmt.Printf("✅ Seeded vaccine refusals for %d patients\n", len(patients))
	return nil
}

func seedVaccineWithdrawals(db *gorm.DB) error {
	var titles []entities.Manual
	db.Where("type = ?", entities.RefTypeVaccineTitle).Find(&titles)

	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients for vaccine withdrawals: %w", err)
	}
	if len(patients) == 0 {
		fmt.Println("⚠️ No patients found, skipping vaccine withdrawals seeding")
		return nil
	}

	for _, patient := range patients {
		withdrawalDate := time.Now().AddDate(0, 0, -rand.Intn(180)) // за последние 6 месяцев
		num := 20240000 + rand.Intn(10000)                          // условный номер приказа/документа

		withdrawal := &entities.VaccineWithdrawal{
			PatientID: patient.ID,
			TitleID:   titles[0].ID,
			Date:      withdrawalDate,
			Num:       num,
			CreatedAt: time.Now(),
		}
		if err := db.Create(withdrawal).Error; err != nil {
			return fmt.Errorf("failed to create vaccine withdrawal for patient %d: %w", patient.ID, err)
		}
		fmt.Printf("✅ Created vaccine withdrawal #%d for %s (date: %s)\n",
			num, patient.FullName, withdrawalDate.Format("2006-01-02"))
	}

	fmt.Printf("✅ Seeded vaccine withdrawals for %d patients\n", len(patients))
	return nil
}

func seedTitrs(db *gorm.DB) error {
	var titles []entities.Manual
	db.Where("type = ?", entities.RefTypeVaccineTitle).Find(&titles)

	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients for titrs: %w", err)
	}
	if len(patients) == 0 {
		fmt.Println("⚠️ No patients found, skipping titrs seeding")
		return nil
	}

	for _, patient := range patients {
		// Создаём титрование у ~40% пациентов
		if rand.Intn(100) < 40 {
			titrDate := time.Now().AddDate(0, 0, -rand.Intn(90)) // за последние 3 месяца
			amount := 100 + rand.Intn(900)                       // например, титр от 100 до 999

			titr := &entities.Titr{
				PatientID: patient.ID,
				TitleID:   titles[0].ID,
				Date:      titrDate,
				CreatedAt: time.Now(),
			}
			if err := db.Create(titr).Error; err != nil {
				return fmt.Errorf("failed to create titer for patient %d: %w", patient.ID, err)
			}
			fmt.Printf("✅ Created titer (amount: %d) for %s (date: %s)\n",
				amount, patient.FullName, titrDate.Format("2006-01-02"))
		}
	}

	fmt.Printf("✅ Seeded titrs for %d patients\n", len(patients))
	return nil
}

func seedReceptions(db *gorm.DB) error {
	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("load patients: %w", err)
	}

	// Для каждого пациента: по 1 приёму на специализацию
	for _, patient := range patients {
		// Получим все уникальные специализации пациента
		var specializations []entities.Specialization
		if err := db.Model(&patient).Association("Specializations").Find(&specializations); err != nil {
			return fmt.Errorf("load patient specializations: %w", err)
		}

		for _, spec := range specializations {
			// Возьмём первый шаблон для этой специализации
			var template entities.ReceptionTemplate
			if err := db.Where("specialization_id = ?", spec.ID).First(&template).Error; err != nil {
				fmt.Printf("⚠️ No template for specialization %s, skip\n", spec.Title)
				continue
			}

			values := generateReceptionValues(spec.Title)
			dataJSON, _ := json.Marshal(values)

			reception := entities.Reception{
				PatientID:        patient.ID,
				SpecializationID: spec.ID,
				TemplateID:       template.ID,
				Data:             json.RawMessage(dataJSON),
				IsCompleted:      true,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			db.Create(&reception)
		}
	}

	fmt.Printf("✅ Seeded receptions (1 per specialization per patient)\n")
	return nil
}

func generateReceptionValues(specTitle string) map[string]interface{} {
	switch specTitle {
	case "Терапевт":
		return map[string]interface{}{
			"complaints":      "Головная боль, слабость",
			"bp_systolic":     120 + rand.Intn(30),
			"bp_diastolic":    80 + rand.Intn(20),
			"heart_rate":      60 + rand.Intn(30),
			"temperature":     36.6 + rand.Float32()*0.8,
			"diagnosis":       "ОРВИ",
			"recommendations": "Постельный режим, обильное питьё",
		}
	case "Невролог":
		return map[string]interface{}{
			"mental_status":    "ясное сознание",
			"motor_function":   "норма",
			"sensory_function": "снижена в правой руке",
			"reflexes":         "живые, симметричные",
			"diagnosis":        "ДЦП",
			"mri_results":      "очаговые изменения в белом веществе",
		}
	case "Травматолог":
		return map[string]interface{}{
			"injury_type":    "Ушиб",
			"localization":   "Левое колено",
			"xray_results":   "Переломов нет",
			"swelling":       true,
			"treatment_plan": "Покой, холод, НПВС",
		}
	case "Психиатр":
		return map[string]interface{}{
			"mood":              "депрессивное",
			"sleep_quality":     "нарушено",
			"appetite":          "снижено",
			"suicidal_ideation": false,
			"diagnosis_icd":     "F32.1",
			"therapy_plan":      "Антидепрессанты + когнитивно-поведенческая терапия",
		}
	default:
		return map[string]interface{}{
			"general_condition": "удовлетворительное",
			"diagnosis":         "Наблюдение",
			"notes":             "Без патологий",
		}
	}
}

// Вспомогательная функция для хэширования пароля
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
