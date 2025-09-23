// repository/migrations.go

package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/auth"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/consent_signatures"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/contactInfo"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/doctor"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/organization"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/patient"
	patientgroup "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/patient_group"
	personalInfo "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/personal_info"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/reception"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/tx"
	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	interfaces.AuthRepository
	interfaces.DoctorRepository
	interfaces.PatientRepository
	interfaces.ContactInfoRepository
	interfaces.PersonalInfoRepository
	interfaces.ReceptionRepository
	interfaces.TxRepository
	interfaces.OrganizationRepository
	interfaces.PatientGroupRepository
	interfaces.ConsentSignatureRepository
}

func NewRepository(cfg *config.Config) (interfaces.Repository, error) {
	//logger := logging.NewModuleLogger("ADAPTER", "POSTGRES", parentLogger)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Вывод в stdout
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Порог для медленных запросов
			LogLevel:                  logger.Info,            // Уровень логирования (Info - все запросы)
			IgnoreRecordNotFoundError: true,                   // Игнорировать ошибки "запись не найдена"
			Colorful:                  true,                   // Цветной вывод
		},
	)

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Выполнение автомиграций
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("ошибка выполнения автомиграций: %w", err)

	}

	return &Repository{
		auth.NewAuthRepository(db),
		doctor.NewDoctorRepository(db),
		patient.NewPatientRepository(db),
		contactInfo.NewContactInfoRepository(db),
		personalInfo.NewPersonalInfoRepository(db),
		reception.NewReceptionRepository(db),
		tx.NewTxRepository(db),
		organization.NewOrganizationRepository(db),
		patientgroup.NewPatientGroupRepository(db),
		consent_signatures.NewConsentSignatureRepository(db),
	}, nil

}

func autoMigrate(db *gorm.DB) error {
	// ✅ Правильный порядок удаления таблиц (зависимые первыми)
	tablesToDelete := []string{
		// Зависимые таблицы
		"receptions",
		"analysis_order_items",
		"analysis_orders",
		"vaccines",
		"fl_gs",
		"contact_infos",
		"personal_infos",
		"patient_statistics",

		//many2many
		"patients_specializations",
		"patients_patient_groups",
		"doctor_specializations",
		"doctor_organizations",
		"harm_points_specializations",

		// Родительские таблицы
		"patients",
		"patient_groups",
		"doctors",
		"specializations",
		"organizations",
		"managers",

		// Справочники
		"document_types",
		"examination_types",
		"examination_views",
		"harm_points",
		"analyses",
		"titles",
		"medications",
		"doses",
		"numbers",
		"certificate_numbers",
		"body_parts",
		"methods",
		"places",
	}

	// Удаление таблиц
	for _, table := range tablesToDelete {
		if db.Migrator().HasTable(table) {
			if err := db.Migrator().DropTable(table); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}
	}

	// ✅ Создание таблиц в правильном порядке (зависимости первыми)
	models := []interface{}{
		// Справочники (без внешних ключей)
		&entities.DocumentType{},
		&entities.ExaminationType{},
		&entities.ExaminationView{},
		&entities.HarmPoint{},
		&entities.Title{},
		&entities.Medication{},
		&entities.Dose{},
		&entities.Number{},
		&entities.CertificateNumber{},
		&entities.BodyPart{},
		&entities.Method{},
		&entities.Place{},
		&entities.Manager{},
		&entities.Analysis{},

		// Основные сущности
		&entities.Specialization{},
		&entities.Doctor{},
		&entities.Organization{},
		&entities.PatientGroup{},

		// Зависимые сущности
		&entities.ContactInfo{},
		&entities.PersonalInfo{},
		&entities.Flg{},

		// Зависимые от основных
		&entities.Patient{},
		&entities.AnalysisOrder{},
		&entities.AnalysisOrderItem{},
		&entities.Vaccine{},
		&entities.Reception{},
		&entities.PatientStatistics{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// ✅ Создаем составной индекс для ускорения фильтрации + сортировки
	if !db.Migrator().HasIndex(&entities.Patient{}, "idx_patients_group_name") {
		if err := db.Exec(`
			CREATE INDEX idx_patients_group_name ON patients (patient_group_id, full_name)
		`).Error; err != nil {
			return fmt.Errorf("failed to create index idx_patients_group_name: %w", err)
		}
		fmt.Println("✅ Created index: idx_patients_group_name on patients(patient_group_id, full_name)")
	}

	if err := seedTestData(db); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	return nil
}

func seedTestData(db *gorm.DB) error {
	// 1. Создаем справочники
	if err := seedDocumentTypes(db); err != nil {
		return fmt.Errorf("failed to seed document types: %w", err)
	}

	if err := seedExaminationTypes(db); err != nil {
		return fmt.Errorf("failed to seed examination types: %w", err)
	}

	if err := seedExaminationViews(db); err != nil {
		return fmt.Errorf("failed to seed examination views: %w", err)
	}

	if err := seedHarmPoints(db); err != nil {
		return fmt.Errorf("failed to seed harm points: %w", err)
	}

	if err := seedSpecializations(db); err != nil {
		return fmt.Errorf("failed to seed specializations: %w", err)
	}

	// 2. Создаем менеджеров
	if err := seedManagers(db); err != nil {
		return fmt.Errorf("failed to seed managers: %w", err)
	}

	// 3. Создаем организации
	if err := seedOrganizations(db); err != nil {
		return fmt.Errorf("failed to seed organizations: %w", err)
	}

	// 4. Создаем группы пациентов
	if err := seedPatientGroups(db); err != nil {
		return fmt.Errorf("failed to seed patient groups: %w", err)
	}

	// 5. Создаем докторов
	if err := seedDoctors(db); err != nil {
		return fmt.Errorf("failed to seed doctors: %w", err)
	}

	// 6. Создаем справочники для вакцин
	if err := seedVaccineDictionaries(db); err != nil {
		return fmt.Errorf("failed to seed vaccine dictionaries: %w", err)
	}

	// 7. Создаем справочники для пациентов
	if err := seedPatientDictionaries(db); err != nil {
		return fmt.Errorf("failed to seed patient dictionaries: %w", err)
	}

	// 8. Создаем анализы
	if err := seedAnalyses(db); err != nil {
		return fmt.Errorf("failed to seed analyses: %w", err)
	}

	// 9. Создаем пациентов
	if err := seedHarmPointsSpecializations(db); err != nil {
		return fmt.Errorf("failed to seed patients: %w", err)
	}

	// 9.1. Создаем пациентов
	if err := seedPatients(db); err != nil {
		return fmt.Errorf("failed to seed patients: %w", err)
	}

	// 9.2. Создаем вакцины для пациентов ← ДОБАВИТЬ ЭТОТ БЛОК
	if err := seedVaccines(db); err != nil {
		return fmt.Errorf("failed to seed vaccines: %w", err)
	}

	// 10. Создаем статистику пациентов
	if err := seedPatientStatistics(db); err != nil {
		return fmt.Errorf("failed to seed patient statistics: %w", err)
	}

	// 11. Создаем направления на анализы
	if err := seedAnalysisOrders(db); err != nil {
		return fmt.Errorf("failed to seed analysis orders: %w", err)
	}

	// 12. Создаем приемы
	if err := seedReceptions(db); err != nil {
		return fmt.Errorf("failed to seed receptions: %w", err)
	}

	return nil
}

// seed/document_types.go
func seedDocumentTypes(db *gorm.DB) error {
	documentTypes := []*entities.DocumentType{
		{Value: "Паспорт РФ"},
		{Value: "Заграничный паспорт"},
		{Value: "Водительское удостоверение"},
		{Value: "Военный билет"},
		{Value: "Свидетельство о рождении"},
	}

	for _, docType := range documentTypes {
		if err := db.Create(docType).Error; err != nil {
			return fmt.Errorf("failed to create document type %s: %w", docType.Value, err)
		}
	}
	return nil
}

// seed/examination_types.go
func seedExaminationTypes(db *gorm.DB) error {
	examinationTypes := []*entities.ExaminationType{
		{Value: "Предварительный"},
		{Value: "Периодический"},
		{Value: "Внеочередной"},
	}

	for _, examType := range examinationTypes {
		if err := db.Create(examType).Error; err != nil {
			return fmt.Errorf("failed to create examination type %s: %w", examType.Value, err)
		}
	}
	return nil
}

// seed/examination_views.go
func seedExaminationViews(db *gorm.DB) error {
	examinationViews := []*entities.ExaminationView{
		{Value: "Осмотр терапевта"},
		{Value: "Осмотр хирурга"},
		{Value: "Осмотр невролога"},
		{Value: "Осмотр офтальмолога"},
		{Value: "Осмотр отоларинголога"},
	}

	for _, examView := range examinationViews {
		if err := db.Create(examView).Error; err != nil {
			return fmt.Errorf("failed to create examination view %s: %w", examView.Value, err)
		}
	}
	return nil
}

// seed/specializations.go
func seedSpecializations(db *gorm.DB) error {
	specializations := []*entities.Specialization{
		{Title: "Терапевт"},
		{Title: "Невролог"},
		{Title: "Травматолог"},
		{Title: "Психиатр"},
		{Title: "Уролог"},
		{Title: "Оториноларинголог"},
		{Title: "Аллерголог"},
		{Title: "Проктолог"},
		{Title: "Кардиолог"},
		{Title: "Хирург"},
	}

	for _, spec := range specializations {
		if err := db.Create(spec).Error; err != nil {
			return fmt.Errorf("failed to create specialization %s: %w", spec.Title, err)
		}
	}
	return nil
}

// seed/managers.go
func seedManagers(db *gorm.DB) error {
	managers := []*entities.Manager{
		{FullName: "Иванов Иван Петрович", Phone: "+79001111111"},
		{FullName: "Петров Петр Иванович", Phone: "+79002222222"},
		{FullName: "Сидоров Сидор Петрович", Phone: "+79003333333"},
	}

	for _, manager := range managers {
		if err := db.Create(manager).Error; err != nil {
			return fmt.Errorf("failed to create manager %s: %w", manager.FullName, err)
		}
	}
	return nil
}

func seedOrganizations(db *gorm.DB) error {
	// Получаем всех менеджеров
	var managers []entities.Manager
	if err := db.Find(&managers).Error; err != nil {
		return fmt.Errorf("failed to get managers: %w", err)
	}

	if len(managers) == 0 {
		return errors.New("no managers found to assign to organizations")
	}

	organizations := []*entities.Organization{
		{Title: "Городская клиника №1", ManagerID: managers[0].ID},
		{Title: "Центральная больница", ManagerID: managers[1].ID},
		{Title: "Медцентр 'Здоровье'", ManagerID: managers[2].ID},
	}

	for _, org := range organizations {
		if err := db.Create(org).Error; err != nil {
			return fmt.Errorf("failed to create organization %s: %w", org.Title, err)
		}
		fmt.Printf("✅ Created organization '%s' with manager ID %d\n", org.Title, org.ManagerID)
	}

	return nil
}

// seed/patient_groups.go
func seedPatientGroups(db *gorm.DB) error {
	var organizations []entities.Organization
	if err := db.Find(&organizations).Error; err != nil {
		return fmt.Errorf("failed to get organizations: %w", err)
	}

	patientGroups := []*entities.PatientGroup{
		{Code: "PG001", OrganizationID: organizations[0].ID},
		{Code: "PG002", OrganizationID: organizations[0].ID},
		{Code: "PG003", OrganizationID: organizations[1].ID},
		{Code: "PG004", OrganizationID: organizations[2].ID},
	}

	for _, group := range patientGroups {
		if err := db.Create(group).Error; err != nil {
			return fmt.Errorf("failed to create patient group %s: %w", group.Code, err)
		}
	}
	return nil
}

// seed/doctors.go
func seedDoctors(db *gorm.DB) error {
	var specializations []entities.Specialization
	var organizations []entities.Organization

	if err := db.Find(&specializations).Error; err != nil {
		return fmt.Errorf("failed to get specializations: %w", err)
	}

	if err := db.Find(&organizations).Error; err != nil {
		return fmt.Errorf("failed to get organizations: %w", err)
	}

	hashPass123 := hashPassword("123")

	doctors := []*entities.Doctor{
		{
			FullName:     "Смирнова Анна Михайловна",
			Phone:        "+79161111111",
			PasswordHash: hashPass123,
		},
		{
			FullName:     "Козлов Владимир Петрович",
			Phone:        "+79162222222",
			PasswordHash: hashPass123,
		},
		{
			FullName:     "Иванов Иван Иванович",
			Phone:        "+79163333333",
			PasswordHash: hashPass123,
		},
		{
			FullName:     "Петрова Мария Сергеевна",
			Phone:        "+79164444444",
			PasswordHash: hashPass123,
		},
		{
			FullName:     "Сидоров Алексей Дмитриевич",
			Phone:        "+79165555555",
			PasswordHash: hashPass123,
		},
	}

	for _, doctor := range doctors {
		if err := db.Create(doctor).Error; err != nil {
			return fmt.Errorf("failed to create doctor %s: %w", doctor.FullName, err)
		}

		// Связываем докторов со специализациями
		if strings.Contains(doctor.FullName, "Смирнова") {
			db.Model(doctor).Association("Specializations").Append(&specializations[0]) // Терапевт
		} else if strings.Contains(doctor.FullName, "Козлов") {
			db.Model(doctor).Association("Specializations").Append(&specializations[0]) // Терапевт
		} else if strings.Contains(doctor.FullName, "Иванов") {
			db.Model(doctor).Association("Specializations").Append(&specializations[1]) // Невролог
		} else if strings.Contains(doctor.FullName, "Петрова") {
			db.Model(doctor).Association("Specializations").Append(&specializations[2]) // Травматолог
		} else if strings.Contains(doctor.FullName, "Сидоров") {
			db.Model(doctor).Association("Specializations").Append(&specializations[3]) // Психиатр
		}

		// Связываем докторов с организациями
		db.Model(doctor).Association("Organizations").Append(&organizations[0])
	}
	return nil
}

// seed/vaccine_dictionaries.go
func seedVaccineDictionaries(db *gorm.DB) error {
	// Titles
	titles := []*entities.Title{
		{Value: "COVID-19 Vaccine"},
		{Value: "Гриппозная вакцина"},
		{Value: "Вакцина против гепатита B"},
		{Value: "Вакцина АКДС"},
	}
	for _, title := range titles {
		if err := db.Create(title).Error; err != nil {
			return fmt.Errorf("failed to create title %s: %w", title.Value, err)
		}
	}

	// Medications
	medications := []*entities.Medication{
		{Value: "Pfizer-BioNTech"},
		{Value: "Moderna"},
		{Value: "ГамКовидВак"},
		{Value: "Совигрипп"},
	}
	for _, med := range medications {
		if err := db.Create(med).Error; err != nil {
			return fmt.Errorf("failed to create medication %s: %w", med.Value, err)
		}
	}

	// Doses
	doses := []*entities.Dose{
		{Value: 0.3},
		{Value: 0.5},
		{Value: 1.0},
	}
	for _, dose := range doses {
		if err := db.Create(dose).Error; err != nil {
			return fmt.Errorf("failed to create dose %f: %w", dose.Value, err)
		}
	}

	// Numbers
	numbers := []*entities.Number{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}
	for _, num := range numbers {
		if err := db.Create(num).Error; err != nil {
			return fmt.Errorf("failed to create number %d: %w", num.Value, err)
		}
	}

	// CertificateNumbers
	certNumbers := []*entities.CertificateNumber{
		{Value: 1001},
		{Value: 1002},
		{Value: 1003},
		{Value: 1004},
	}
	for _, cn := range certNumbers {
		if err := db.Create(cn).Error; err != nil {
			return fmt.Errorf("failed to create certificate number %d: %w", cn.Value, err)
		}
	}

	// BodyParts
	bodyParts := []*entities.BodyPart{
		{Value: "Правая рука"},
		{Value: "Левая рука"},
		{Value: "Правая нога"},
		{Value: "Левая нога"},
		{Value: "Ягодица"},
	}
	for _, bp := range bodyParts {
		if err := db.Create(bp).Error; err != nil {
			return fmt.Errorf("failed to create body part %s: %w", bp.Value, err)
		}
	}

	// Methods
	methods := []*entities.Method{
		{Value: "Внутримышечный"},
		{Value: "Подкожный"},
		{Value: "Внутривенный"},
	}
	for _, method := range methods {
		if err := db.Create(method).Error; err != nil {
			return fmt.Errorf("failed to create method %s: %w", method.Value, err)
		}
	}

	// Places
	places := []*entities.Place{
		{Value: "Поликлиника"},
		{Value: "Стационар"},
		{Value: "Выездной пункт"},
	}
	for _, place := range places {
		if err := db.Create(place).Error; err != nil {
			return fmt.Errorf("failed to create place %s: %w", place.Value, err)
		}
	}

	return nil
}

// seed/patient_dictionaries.go
func seedPatientDictionaries(db *gorm.DB) error {
	documentTypes := []*entities.DocumentType{
		{Value: "Паспорт РФ"},
		{Value: "Заграничный паспорт"},
		{Value: "Водительское удостоверение"},
		{Value: "Военный билет"},
		{Value: "Свидетельство о рождении"},
	}

	for _, docType := range documentTypes {
		if err := db.Create(docType).Error; err != nil {
			return fmt.Errorf("failed to create document type %s: %w", docType.Value, err)
		}
	}

	examinationTypes := []*entities.ExaminationType{
		{Value: "Предварительный"},
		{Value: "Периодический"},
		{Value: "Внеочередной"},
	}

	for _, examType := range examinationTypes {
		if err := db.Create(examType).Error; err != nil {
			return fmt.Errorf("failed to create examination type %s: %w", examType.Value, err)
		}
	}

	examinationViews := []*entities.ExaminationView{
		{Value: "Осмотр терапевта"},
		{Value: "Осмотр хирурга"},
		{Value: "Осмотр невролога"},
		{Value: "Осмотр офтальмолога"},
		{Value: "Осмотр отоларинголога"},
	}

	for _, examView := range examinationViews {
		if err := db.Create(examView).Error; err != nil {
			return fmt.Errorf("failed to create examination view %s: %w", examView.Value, err)
		}
	}

	harmPoints := []*entities.HarmPoint{
		{Value: 1.0},
		{Value: 2.0},
		{Value: 3.0},
		{Value: 3.1},
		{Value: 3.2},
		{Value: 3.3},
		{Value: 4.0},
	}

	for _, harmPoint := range harmPoints {
		if err := db.Create(harmPoint).Error; err != nil {
			return fmt.Errorf("failed to create harm point %f: %w", harmPoint.Value, err)
		}
	}

	return nil
}

// seed/analyses.go
func seedAnalyses(db *gorm.DB) error {
	analyses := []*entities.Analysis{
		{Name: "Общий анализ крови", Price: 500},
		{Name: "Биохимический анализ крови", Price: 1200},
		{Name: "Анализ мочи", Price: 300},
		{Name: "ЭКГ", Price: 400},
		{Name: "Флюрография", Price: 800},
		{Name: "УЗИ брюшной полости", Price: 1500},
		{Name: "Анализ на глюкозу", Price: 200},
		{Name: "Анализ на холестерин", Price: 250},
	}

	for _, analysis := range analyses {
		if err := db.Create(analysis).Error; err != nil {
			return fmt.Errorf("failed to create analysis %s: %w", analysis.Name, err)
		}
	}
	return nil
}

// seed/harm_points.go

func seedHarmPoints(db *gorm.DB) error {
	harmPoints := []*entities.HarmPoint{
		{Value: 1.0},
		{Value: 2.0},
		{Value: 3.0},
		{Value: 3.1},
		{Value: 3.2},
		{Value: 3.3},
		{Value: 4.0},
	}

	for _, harmPoint := range harmPoints {
		if err := db.Create(harmPoint).Error; err != nil {
			return fmt.Errorf("failed to create harm point %f: %w", harmPoint.Value, err)
		}
	}

	return nil
}

// seed/harm_points_specializations.go - связывание HarmPoint со Specialization

func seedHarmPointsSpecializations(db *gorm.DB) error {
	// Получаем все HarmPoint и Specialization
	var harmPoints []entities.HarmPoint
	var specializations []entities.Specialization

	if err := db.Find(&harmPoints).Error; err != nil {
		return fmt.Errorf("failed to get harm points: %w", err)
	}

	if err := db.Find(&specializations).Error; err != nil {
		return fmt.Errorf("failed to get specializations: %w", err)
	}

	// Создаем связи: каждый HarmPoint связан с 2-3 специализациями
	for i, harmPoint := range harmPoints {
		specCount := 2
		if i%3 == 0 {
			specCount = 3
		}

		var specsToLink []entities.Specialization
		for j := 0; j < specCount && j < len(specializations); j++ {
			specIndex := (i + j) % len(specializations)
			specsToLink = append(specsToLink, specializations[specIndex])
		}

		// ✅ GORM автоматически создаст записи в harm_points_specializations
		if err := db.Model(&harmPoint).Association("Specializations").Append(&specsToLink); err != nil {
			return fmt.Errorf("failed to link harm point %f with specializations: %w", harmPoint.Value, err)
		}

		fmt.Printf("✅ Linked harm point %f with %d specializations\n", harmPoint.Value, len(specsToLink))
	}

	return nil
}

// seed/patients.go - ИСПРАВЛЕННЫЙ

func seedPatients(db *gorm.DB) error {
	// Получаем все справочники
	var examinationTypes []entities.ExaminationType
	var examinationViews []entities.ExaminationView
	var harmPoints []entities.HarmPoint
	var docTypes []entities.DocumentType
	var patientGroup []entities.PatientGroup

	if err := db.Find(&examinationTypes).Error; err != nil {
		return fmt.Errorf("failed to get examination types: %w", err)
	}

	if err := db.Find(&examinationViews).Error; err != nil {
		return fmt.Errorf("failed to get examination views: %w", err)
	}

	if err := db.Find(&harmPoints).Error; err != nil {
		return fmt.Errorf("failed to get harm points: %w", err)
	}

	if err := db.Find(&docTypes).Error; err != nil {
		return fmt.Errorf("failed to get docTypes: %w", err)
	}

	if err := db.Find(&patientGroup).Error; err != nil {
		return fmt.Errorf("failed to get patientGroup: %w", err)
	}

	patientsData := []struct {
		FullName       string
		BirthDate      time.Time
		IsMale         bool
		Position       string
		Division       string
		ExamTypeID     uint
		ExamViewID     uint
		HarmPointID    uint
		PatientGroupID uint
	}{
		{
			"Иванов Петр Сергеевич",
			time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
			true,
			"Программист",
			"IT отдел",
			examinationTypes[0].ID,
			examinationViews[0].ID,
			harmPoints[0].ID,
			patientGroup[0].ID,
		},
		{
			"Петрова Мария Ивановна",
			time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC),
			false,
			"Дизайнер",
			"Дизайн отдел",
			examinationTypes[1].ID,
			examinationViews[1].ID,
			harmPoints[1].ID,
			patientGroup[1].ID,
		},
		{
			"Сидоров Алексей Петрович",
			time.Date(1978, 12, 3, 0, 0, 0, 0, time.UTC),
			true,
			"Менеджер",
			"Управление",
			examinationTypes[2].ID,
			examinationViews[2].ID,
			harmPoints[2].ID,
			patientGroup[2].ID,
		},
	}

	for i, pd := range patientsData {
		// ✅ 1. Создаем контактную информацию
		contactInfo := &entities.ContactInfo{
			Phone:     fmt.Sprintf("+7900%06d", 111111+i),
			Email:     fmt.Sprintf("patient%d@example.com", i+1),
			Address:   fmt.Sprintf("г. Москва, ул. Ленина, д. %d", i+1),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(contactInfo).Error; err != nil {
			return fmt.Errorf("failed to create contact info: %w", err)
		}

		// ✅ 2. Создаем персональную информацию
		personalInfo := &entities.PersonalInfo{
			DocNumber:      fmt.Sprintf("%06d", 123456+i),
			DocSeries:      fmt.Sprintf("451%d", i),
			SNILS:          fmt.Sprintf("123-456-789 %02d", i),
			OMS:            fmt.Sprintf("123456789012345%d", i),
			DocumentTypeID: docTypes[0].ID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := db.Create(personalInfo).Error; err != nil {
			return fmt.Errorf("failed to create personal info: %w", err)
		}

		// ✅ 3. Создаем Flg запись (может быть NULL)
		var flgID *uint
		flg := &entities.Flg{
			IsCompleted:  false,
			Organization: "Городская поликлиника",
			Number:       10000 + i,
			Result:       "Норма",
			Date:         time.Now(),
		}
		if err := db.Create(flg).Error; err != nil {
			return fmt.Errorf("failed to create flg: %w", err)
		}
		flgID = &flg.ID

		// ✅ 4. Создаем пустое направление на анализы
		analysisOrder := &entities.AnalysisOrder{
			OrderNumber: fmt.Sprintf("ORD-%06d", 0), // Временный номер
			TotalAmount: 0,                          // Пустое направление
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := db.Create(analysisOrder).Error; err != nil {
			return fmt.Errorf("failed to create analysis order: %w", err)
		}

		// Обновляем номер с правильным ID
		analysisOrder.OrderNumber = fmt.Sprintf("ORD-%06d", analysisOrder.ID)
		if err := db.Save(analysisOrder).Error; err != nil {
			return fmt.Errorf("failed to update analysis order number: %w", err)
		}

		// ✅ 5. Создаем пациента со всеми обязательными связями
		patient := &entities.Patient{
			FullName:          pd.FullName,
			BirthDate:         pd.BirthDate,
			IsMale:            pd.IsMale,
			Position:          pd.Position,
			Division:          pd.Division,
			PatientGroupID:    uint(i%3 + 1), // Группа 1 или 2
			ExaminationTypeID: pd.ExamTypeID,
			ExaminationViewID: pd.ExamViewID,
			HarmPointID:       pd.HarmPointID,
			PersonalInfoID:    personalInfo.ID,
			ContactInfoID:     contactInfo.ID,
			FlgID:             flgID,
			AnalysisOrderID:   analysisOrder.ID,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := db.Create(patient).Error; err != nil {
			return fmt.Errorf("failed to create patient %s: %w", patient.FullName, err)
		}

		// Обновляем AnalysisOrder с PatientID
		analysisOrder.PatientID = patient.ID
		if err := db.Save(analysisOrder).Error; err != nil {
			return fmt.Errorf("failed to update analysis order with patient ID: %w", err)
		}

		// ✅ 6. Автоматически связываем пациента со специализациями через HarmPoint
		var harmPoint entities.HarmPoint
		if err := db.Preload("Specializations").First(&harmPoint, pd.HarmPointID).Error; err != nil {
			return fmt.Errorf("failed to get harm point with specializations: %w", err)
		}

		// GORM автоматически создаст связи в patients_specializations
		if err := db.Model(patient).Association("Specializations").Append(&harmPoint.Specializations); err != nil {
			return fmt.Errorf("failed to link patient with specializations: %w", err)
		}

		fmt.Printf("✅ Created patient %s with %d specializations\n",
			patient.FullName, len(harmPoint.Specializations))
	}

	return nil
}

func seedVaccines(db *gorm.DB) error {
	// 1. Получаем всех пациентов
	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	if len(patients) == 0 {
		fmt.Println("⚠️ No patients found, skipping vaccine seeding")
		return nil
	}

	// 2. Получаем все справочники
	var titles []entities.Title
	var medications []entities.Medication
	var doses []entities.Dose
	var numbers []entities.Number
	var certNumbers []entities.CertificateNumber
	var bodyParts []entities.BodyPart
	var methods []entities.Method
	var places []entities.Place

	if err := db.Find(&titles).Error; err != nil {
		return fmt.Errorf("failed to load titles: %w", err)
	}
	if err := db.Find(&medications).Error; err != nil {
		return fmt.Errorf("failed to load medications: %w", err)
	}
	if err := db.Find(&doses).Error; err != nil {
		return fmt.Errorf("failed to load doses: %w", err)
	}
	if err := db.Find(&numbers).Error; err != nil {
		return fmt.Errorf("failed to load numbers: %w", err)
	}
	if err := db.Find(&certNumbers).Error; err != nil {
		return fmt.Errorf("failed to load certificate numbers: %w", err)
	}
	if err := db.Find(&bodyParts).Error; err != nil {
		return fmt.Errorf("failed to load body parts: %w", err)
	}
	if err := db.Find(&methods).Error; err != nil {
		return fmt.Errorf("failed to load methods: %w", err)
	}
	if err := db.Find(&places).Error; err != nil {
		return fmt.Errorf("failed to load places: %w", err)
	}
	// 3. Для каждого пациента создаем 1-3 вакцины
	for _, patient := range patients {
		vaccineCount := 1 + rand.Intn(3) // 1, 2 или 3 вакцины

		for i := 0; i < vaccineCount; i++ {
			// Случайные индексы (с проверкой на пустые срезы)
			titleIdx := rand.Intn(len(titles))
			medIdx := rand.Intn(len(medications))
			doseIdx := rand.Intn(len(doses))
			numIdx := rand.Intn(len(numbers))
			certIdx := rand.Intn(len(certNumbers))
			bodyPartIdx := rand.Intn(len(bodyParts))
			methodIdx := rand.Intn(len(methods))
			placeIdx := rand.Intn(len(places))

			// Генерируем случайную дату за последние 2 года
			startDate := time.Now().AddDate(-2, 0, 0)
			days := rand.Intn(730) // 2 года = 730 дней
			vaccineDate := startDate.AddDate(0, 0, days)

			// Случайный результат
			results := []string{"Успешно", "Реакция", "Отменено", "Перенесено"}
			result := results[rand.Intn(len(results))]

			// Случайные флаги
			isCompleted := rand.Intn(2) == 1 // 50% шанс
			isRefusal := !isCompleted && rand.Intn(2) == 1
			isExemption := !isCompleted && !isRefusal && rand.Intn(2) == 1
			IsTiter := !isCompleted && !isRefusal && rand.Intn(2) == 1

			// Случайные числовые поля
			var titerAmount *int
			if rand.Intn(2) == 1 {
				val := 100 + rand.Intn(900) // 100-999
				titerAmount = &val
			}

			var medWithdrawlNum *int
			if isExemption && rand.Intn(2) == 1 {
				val := 2024000 + rand.Intn(1000)
				medWithdrawlNum = &val
			}

			// Создаем вакцину
			vaccine := &entities.Vaccine{
				Date:                vaccineDate,
				IsCompleted:         isCompleted,
				PatientID:           patient.ID,
				IsRefusal:           isRefusal,
				IsExemption:         isExemption,
				IsTiter:             IsTiter,
				TiterAmount:         titerAmount,
				MedWithdrawlNum:     medWithdrawlNum,
				Result:              &result,
				TitleID:             &titles[titleIdx].ID,
				MedicationID:        &medications[medIdx].ID,
				DoseID:              &doses[doseIdx].ID,
				NumberID:            &numbers[numIdx].ID,
				CertificateNumberID: &certNumbers[certIdx].ID,
				BodyPartID:          &bodyParts[bodyPartIdx].ID,
				MethodID:            &methods[methodIdx].ID,
				PlaceID:             &places[placeIdx].ID,
				CreatedAt:           time.Now(),
			}

			if err := db.Create(vaccine).Error; err != nil {
				return fmt.Errorf("failed to create vaccine for patient %d: %w", patient.ID, err)
			}

			fmt.Printf("✅ Created vaccine '%s' for patient %s (ID: %d)\n",
				titles[titleIdx].Value, patient.FullName, vaccine.ID)
		}
	}

	fmt.Printf("✅ Seeded vaccines for %d patients\n", len(patients))
	return nil
}

// seed/patient_statistics.go - создание нулевой статистики
func seedPatientStatistics(db *gorm.DB) error {
	// Получаем всех пациентов
	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	// Создаем нулевую статистику для каждого пациента
	for _, patient := range patients {
		stats := &entities.PatientStatistics{
			PatientID:              patient.ID,
			TotalReceptions:        0, // ✅ Нулевая статистика
			CompletedReceptions:    0, // ✅ Нулевая статистика
			TotalAnalysisOrders:    0, // ✅ Нулевая статистика
			CompletedAnalysisItems: 0, // ✅ Нулевая статистика
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}

		if err := db.Create(stats).Error; err != nil {
			if !strings.Contains(err.Error(), "duplicate") {
				return fmt.Errorf("failed to create stats for patient %d: %w", patient.ID, err)
			}
		}

		// Обновляем пациента с StatisticsID
		patient.StatisticsID = stats.ID
		if err := db.Save(&patient).Error; err != nil {
			return fmt.Errorf("failed to update patient %d with statistics ID: %w", patient.ID, err)
		}

		fmt.Printf("✅ Created zero statistics for patient %d\n", patient.ID)
	}

	return nil
}

// seed/analysis_orders.go - создание направлений на анализы
func seedAnalysisOrders(db *gorm.DB) error {
	// Получаем пациентов и анализы
	var patients []entities.Patient
	var analyses []entities.Analysis

	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	if err := db.Find(&analyses).Error; err != nil {
		return fmt.Errorf("failed to get analyses: %w", err)
	}

	// Обновляем существующие направления пациентов
	for i, patient := range patients {
		var analysisOrder entities.AnalysisOrder
		if err := db.First(&analysisOrder, patient.AnalysisOrderID).Error; err != nil {
			return fmt.Errorf("failed to get analysis order for patient %d: %w", patient.ID, err)
		}

		// Добавляем 2-3 анализа в направление
		orderCount := 2
		if i%2 == 0 {
			orderCount = 3
		}

		var orderItems []entities.AnalysisOrderItem
		totalAmount := uint(0)

		for j := 0; j < orderCount && j < len(analyses); j++ {
			analysisIndex := (i + j) % len(analyses)
			analysis := analyses[analysisIndex]

			totalAmount += analysis.Price

			item := entities.AnalysisOrderItem{
				OrderID:           analysisOrder.ID,
				AnalysisID:        analysis.ID,
				PriceAtAssignment: analysis.Price,
				IsCompleted:       i%3 != 0, // 2 из 3 анализов завершены
				CreatedAt:         time.Now().Add(-time.Duration((i+j)*24) * time.Hour),
				UpdatedAt:         time.Now().Add(-time.Duration((i+j)*12) * time.Hour),
			}

			if item.IsCompleted {
				completedAt := time.Now().Add(-time.Duration((i+j)*6) * time.Hour)
				item.CompletedAt = &completedAt
			}

			orderItems = append(orderItems, item)
		}

		// Обновляем направление
		analysisOrder.TotalAmount = totalAmount
		analysisOrder.UpdatedAt = time.Now()
		if err := db.Save(&analysisOrder).Error; err != nil {
			return fmt.Errorf("failed to update analysis order for patient %d: %w", patient.ID, err)
		}

		// Создаем элементы направления
		if len(orderItems) > 0 {
			if err := db.Create(&orderItems).Error; err != nil {
				return fmt.Errorf("failed to create order items for patient %d: %w", patient.ID, err)
			}
		}

		fmt.Printf("✅ Updated analysis order %d for patient %d with %d items\n",
			analysisOrder.ID, patient.ID, len(orderItems))
	}

	return nil
}

// seed/receptions.go - ОБНОВЛЁННАЯ ВЕРСИЯ

func seedReceptions(db *gorm.DB) error {
	var patients []entities.Patient

	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	// Получаем все специализации по названиям — будем использовать по Title
	var specializations []entities.Specialization
	if err := db.Find(&specializations).Error; err != nil {
		return fmt.Errorf("failed to get specializations: %w", err)
	}

	// Вспомогательная функция: найти специализацию по названию
	findSpecByTitle := func(title string) *entities.Specialization {
		for i := range specializations {
			if specializations[i].Title == title {
				return &specializations[i]
			}
		}
		return nil
	}

	// Явно заданные приемы для пациентов
	type ReceptionTemplate struct {
		SpecializationTitle string
		Values              map[string]interface{}
		Schema              []map[string]interface{}
	}

	patientReceptions := map[string][]ReceptionTemplate{
		"Иванов Петр Сергеевич": {
			{
				SpecializationTitle: "Травматолог",
				Values: map[string]interface{}{
					"injury_type":      "Перелом",
					"injury_mechanism": "Падение с высоты",
					"localization":     "Правая нога, большеберцовая кость",
					"xray_results":     "Перелом без смещения",
					"fracture":         true,
					"treatment_plan":   "Иммобилизация, повторный осмотр через 2 недели",
					"surgeon_name":     "Др. Сидоров А.В.",
					"operation_date":   "2024-02-15",
				},
				Schema: []map[string]interface{}{
					{"name": "injury_type", "type": "string", "required": true, "description": "Тип травмы", "example": "Перелом", "min_length": 1, "max_length": 100},
					{"name": "injury_mechanism", "type": "string", "required": true, "description": "Механизм получения травмы", "min_length": 1},
					{"name": "fracture", "type": "boolean", "required": true, "description": "Наличие перелома"},
					{"name": "treatment_plan", "type": "string", "required": true, "description": "План лечения"},
					{"name": "surgeon_name", "type": "string", "required": false, "description": "Имя хирурга"},
				},
			},
			{
				SpecializationTitle: "Невролог",
				Values: map[string]interface{}{
					"mental_status":  "ясное сознание",
					"motor_function": "слабость в правой руке",
					"diagnosis":      "ДЦП",
					"treatment_plan": "физиотерапия, ЛФК",
					"eeg_results":    "норма",
					"mri_scan":       "есть отклонения",
				},
				Schema: []map[string]interface{}{
					{"name": "mental_status", "type": "string", "required": true, "description": "Сознание пациента"},
					{"name": "diagnosis", "type": "string", "required": true, "description": "Диагноз"},
					{"name": "eeg_results", "type": "string", "required": false, "description": "Результаты ЭЭГ"},
				},
			},
		},
		"Петрова Мария Ивановна": {
			{
				SpecializationTitle: "Психиатр",
				Values: map[string]interface{}{
					"mental_status": "ясное сознание",
					"mood":          "подавленное",
					"risk_assessment": map[string]interface{}{
						"suicide":   false,
						"self_harm": false,
						"violence":  false,
					},
					"diagnosis_icd": "F32.0",
					"therapy_plan":  "Психотерапия",
				},
				Schema: []map[string]interface{}{
					{"name": "mental_status", "type": "string", "required": true, "description": "Психическое состояние"},
					{"name": "mood", "type": "string", "required": true, "description": "Настроение"},
					{"name": "risk_assessment", "type": "object", "required": true, "description": "Оценка рисков"},
				},
			},
			{
				SpecializationTitle: "Невролог",
				Values: map[string]interface{}{
					"mental_status":  "спутанное",
					"diagnosis":      "Мигрень",
					"treatment_plan": "Покой, препараты",
					"eeg_results":    "отклонения",
				},
				Schema: []map[string]interface{}{
					{"name": "mental_status", "type": "string", "required": true, "description": "Сознание пациента"},
					{"name": "diagnosis", "type": "string", "required": true, "description": "Диагноз"},
					{"name": "eeg_results", "type": "string", "required": false, "description": "Результаты ЭЭГ"},
				},
			},
		},
		"Сидоров Алексей Петрович": {
			{
				SpecializationTitle: "Травматолог",
				Values: map[string]interface{}{
					"injury_type":      "Растяжение",
					"injury_mechanism": "Спорт",
					"localization":     "Левое плечо",
					"xray_results":     "Без перелома",
					"fracture":         false,
					"treatment_plan":   "Покой, компрессы",
					"surgeon_name":     "",
					"operation_date":   "",
				},
				Schema: []map[string]interface{}{
					{"name": "injury_type", "type": "string", "required": true, "description": "Тип травмы", "min_length": 1, "max_length": 100},
					{"name": "injury_mechanism", "type": "string", "required": true, "description": "Механизм получения травмы", "min_length": 1},
					{"name": "fracture", "type": "boolean", "required": true, "description": "Наличие перелома"},
					{"name": "treatment_plan", "type": "string", "required": true, "description": "План лечения"},
					{"name": "surgeon_name", "type": "string", "required": false, "description": "Имя хирурга"},
				},
			},
		},
	}

	// Создаем приемы
	for i, patient := range patients {
		receptions, exists := patientReceptions[patient.FullName]
		if !exists {
			fmt.Printf("⚠️  No reception templates defined for patient: %s\n", patient.FullName)
			continue
		}

		for j, rt := range receptions {
			spec := findSpecByTitle(rt.SpecializationTitle)
			if spec == nil {
				return fmt.Errorf("specialization '%s' not found for patient %s", rt.SpecializationTitle, patient.FullName)
			}

			// Сериализуем Values и Schema в JSON
			valuesJSON, err := json.Marshal(rt.Values)
			if err != nil {
				return fmt.Errorf("failed to marshal values for patient %s: %w", patient.FullName, err)
			}

			schemaJSON, err := json.Marshal(rt.Schema)
			if err != nil {
				return fmt.Errorf("failed to marshal schema for patient %s: %w", patient.FullName, err)
			}

			// Объединяем в один JSON объект
			combinedData := map[string]json.RawMessage{
				"values": json.RawMessage(valuesJSON),
				"schema": json.RawMessage(schemaJSON),
			}

			combinedJSON, err := json.Marshal(combinedData)
			if err != nil {
				return fmt.Errorf("failed to marshal combined data: %w", err)
			}

			reception := &entities.Reception{
				PatientID:        patient.ID,
				SpecializationID: spec.ID,
				IsCompleted:      i%3 != 0, // сохраним старую логику для разнообразия
				Data:             json.RawMessage(combinedJSON),
				CreatedAt:        time.Now().Add(-time.Duration((i+j)*24) * time.Hour), // немного разнесем по времени
				UpdatedAt:        time.Now().Add(-time.Duration((i+j)*12) * time.Hour),
			}

			if err := db.Create(reception).Error; err != nil {
				if !strings.Contains(err.Error(), "duplicate") &&
					!strings.Contains(err.Error(), "idx_patient_specialization") {
					return fmt.Errorf("failed to create reception for patient %d, spec %s: %w",
						patient.ID, spec.Title, err)
				}
			}

			fmt.Printf("✅ Created reception for %s with specialization '%s'\n", patient.FullName, spec.Title)
		}
	}

	return nil
}

// Вспомогательная функция для хэширования пароля
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
