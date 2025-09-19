// repository/migrations.go

package repository

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/auth"
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
	}, nil

}

func autoMigrate(db *gorm.DB) error {
	log.Println("🚀 Starting auto-migration...")

	// ✅ Правильный порядок удаления таблиц (зависимые первыми)
	tablesToDelete := []string{
		// Зависимые таблицы
		"receptions",
		"analysis_order_items",
		"analysis_orders",
		"vaccines",
		"fl_gs",
		"patients_specializations",
		"patients_patient_groups",
		"doctor_specializations",
		"doctor_organizations",
		"contact_infos",
		"personal_infos",
		"patient_statistics",

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

	// Удаляем таблицы в правильном порядке
	log.Println("🗑️  Dropping existing tables...")
	for _, table := range tablesToDelete {
		if db.Migrator().HasTable(table) {
			if err := db.Migrator().DropTable(table); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
			log.Printf("✅ Dropped table: %s", table)
		}
	}

	// ✅ Создаем таблицы в правильном порядке (зависимости первыми)
	log.Println("🏗️  Creating new tables...")
	models := []interface{}{
		// Справочники (без внешних ключей)
		&entities.Manager{},
		&entities.Specialization{},
		&entities.Organization{},
		&entities.PatientGroup{},
		&entities.Doctor{},
		&entities.Patient{},
		&entities.ContactInfo{},
		&entities.PersonalInfo{},
		&entities.Flg{},
		&entities.AnalysisOrder{},
		&entities.AnalysisOrderItem{},
		&entities.Vaccine{},
		&entities.Reception{},
		&entities.PatientStatistics{},

		// Справочники для вакцин
		&entities.Title{},
		&entities.Medication{},
		&entities.Dose{},
		&entities.Number{},
		&entities.CertificateNumber{},
		&entities.BodyPart{},
		&entities.Method{},
		&entities.Place{},

		// Справочники для пациентов
		&entities.DocumentType{},
		&entities.ExaminationType{},
		&entities.ExaminationView{},
		&entities.HarmPoint{},
		&entities.Analysis{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}
	log.Println("✅ Tables created successfully")

	// ✅ Заполняем тестовыми данными
	log.Println("🌱 Seeding test data...")
	if err := seedTestData(db); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}
	log.Println("✅ Test data seeded successfully")

	return nil
}

// seed/test_data.go

func seedTestData(db *gorm.DB) error {
	// 1. Создаем справочники
	if err := seedManagers(db); err != nil {
		return err
	}

	if err := seedDocumentTypes(db); err != nil {
		return err
	}

	if err := seedExaminationTypes(db); err != nil {
		return err
	}

	if err := seedExaminationViews(db); err != nil {
		return err
	}

	if err := seedHarmPoints(db); err != nil {
		return err
	}

	if err := seedSpecializations(db); err != nil {
		return err
	}

	if err := seedTitles(db); err != nil {
		return err
	}

	if err := seedMedications(db); err != nil {
		return err
	}

	if err := seedDoses(db); err != nil {
		return err
	}

	if err := seedNumbers(db); err != nil {
		return err
	}

	if err := seedCertificateNumbers(db); err != nil {
		return err
	}

	if err := seedBodyParts(db); err != nil {
		return err
	}

	if err := seedMethods(db); err != nil {
		return err
	}

	if err := seedPlaces(db); err != nil {
		return err
	}

	// 2. Создаем основные сущности
	if err := seedOrganizations(db); err != nil {
		return err
	}

	if err := seedPatientGroups(db); err != nil {
		return err
	}

	if err := seedDoctors(db); err != nil {
		return err
	}

	// 3. Создаем пациентов
	if err := seedPatients(db); err != nil {
		return err
	}

	// 4. Создаем зависимые сущности
	if err := seedContactInfos(db); err != nil {
		return err
	}

	if err := seedPersonalInfos(db); err != nil {
		return err
	}

	if err := seedFlgs(db); err != nil {
		return err
	}

	if err := seedVaccines(db); err != nil {
		return err
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

// seed/titles.go
func seedTitles(db *gorm.DB) error {
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
	return nil
}

// seed/medications.go
func seedMedications(db *gorm.DB) error {
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
	return nil
}

// seed/doses.go
func seedDoses(db *gorm.DB) error {
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
	return nil
}

// seed/numbers.go
func seedNumbers(db *gorm.DB) error {
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
	return nil
}

// seed/certificate_numbers.go
func seedCertificateNumbers(db *gorm.DB) error {
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
	return nil
}

// seed/body_parts.go
func seedBodyParts(db *gorm.DB) error {
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
	return nil
}

// seed/methods.go
func seedMethods(db *gorm.DB) error {
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
	return nil
}

// seed/places.go
func seedPlaces(db *gorm.DB) error {
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

// seed/organizations.go
func seedOrganizations(db *gorm.DB) error {
	var managers []entities.Manager
	if err := db.Find(&managers).Error; err != nil {
		return fmt.Errorf("failed to get managers: %w", err)
	}

	organizations := []*entities.Organization{
		{Title: "Городская поликлиника №1", ManagerID: managers[0].ID},
		{Title: "Областная больница", ManagerID: managers[1].ID},
		{Title: "Частная клиника МедСервис", ManagerID: managers[2].ID},
	}

	for _, org := range organizations {
		if err := db.Create(org).Error; err != nil {
			return fmt.Errorf("failed to create organization %s: %w", org.Title, err)
		}
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

// seed/patients.go - ИСПРАВЛЕННЫЙ
func seedPatients(db *gorm.DB) error {
	var organizations []entities.Organization
	var patientGroups []entities.PatientGroup

	if err := db.Find(&organizations).Error; err != nil {
		return fmt.Errorf("failed to get organizations: %w", err)
	}

	if err := db.Find(&patientGroups).Error; err != nil {
		return fmt.Errorf("failed to get patient groups: %w", err)
	}

	patientsData := []struct {
		FullName       string
		BirthDate      time.Time
		IsMale         bool
		Position       string
		Division       string
		OrganizationID *uint
	}{
		{
			"Иванов Петр Сергеевич",
			time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
			true,
			"Программист",
			"IT отдел",
			&organizations[0].ID,
		},
		{
			"Петрова Мария Ивановна",
			time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC),
			false,
			"Дизайнер",
			"Дизайн отдел",
			&organizations[1].ID,
		},
		{
			"Сидоров Алексей Петрович",
			time.Date(1978, 12, 3, 0, 0, 0, 0, time.UTC),
			true,
			"Менеджер",
			"Управление",
			&organizations[2].ID,
		},
	}

	for i, pd := range patientsData {
		patient := &entities.Patient{
			FullName:       pd.FullName,
			BirthDate:      pd.BirthDate,
			IsMale:         pd.IsMale,
			Position:       pd.Position,
			Division:       pd.Division,
			OrganizationID: pd.OrganizationID,
			CreatedAt:      time.Now().Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt:      time.Now().Add(-time.Duration(i*12) * time.Hour),
		}

		if err := db.Create(patient).Error; err != nil {
			return fmt.Errorf("failed to create patient %s: %w", patient.FullName, err)
		}

		// Связываем пациента с группами
		if len(patientGroups) > 0 {
			db.Model(patient).Association("PatientGroups").Append(&patientGroups[0])
		}

		// Связываем пациента со специализациями
		var specializations []entities.Specialization
		db.Limit(3).Find(&specializations)
		db.Model(patient).Association("Specializations").Append(&specializations)
	}

	return nil
}

// seed/contact_infos.go
func seedContactInfos(db *gorm.DB) error {
	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	for i, patient := range patients {
		contactInfo := &entities.ContactInfo{
			Phone:     fmt.Sprintf("+7900%06d", 111111+i),
			Email:     fmt.Sprintf("patient%d@example.com", i+1),
			Address:   fmt.Sprintf("г. Москва, ул. Ленина, д. %d", i+1),
			CreatedAt: time.Now().Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(i*12) * time.Hour),
		}

		if err := db.Create(contactInfo).Error; err != nil {
			return fmt.Errorf("failed to create contact info: %w", err)
		}

		// Обновляем пациента с ContactInfoID
		patient.ContactInfoID = &contactInfo.ID
		if err := db.Save(&patient).Error; err != nil {
			return fmt.Errorf("failed to update patient with contact info: %w", err)
		}
	}

	return nil
}

// seed/personal_infos.go
func seedPersonalInfos(db *gorm.DB) error {
	var patients []entities.Patient
	var documentTypes []entities.DocumentType

	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	if err := db.Find(&documentTypes).Error; err != nil {
		return fmt.Errorf("failed to get document types: %w", err)
	}

	for i, patient := range patients {
		var docTypeID *uint
		if len(documentTypes) > 0 {
			docTypeID = &documentTypes[0].ID
		}

		personalInfo := &entities.PersonalInfo{
			DocNumber:      fmt.Sprintf("%06d", 123456+i),
			DocSeries:      fmt.Sprintf("451%d", i),
			SNILS:          fmt.Sprintf("123-456-789 %02d", i),
			OMS:            fmt.Sprintf("123456789012345%d", i),
			DocumentTypeID: docTypeID,
			CreatedAt:      time.Now().Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt:      time.Now().Add(-time.Duration(i*12) * time.Hour),
		}

		if err := db.Create(personalInfo).Error; err != nil {
			return fmt.Errorf("failed to create personal info: %w", err)
		}

		// Обновляем пациента с PersonalInfoID
		patient.PersonalInfoID = &personalInfo.ID
		if err := db.Save(&patient).Error; err != nil {
			return fmt.Errorf("failed to update patient with personal info: %w", err)
		}
	}

	return nil
}

// seed/flgs.go
func seedFlgs(db *gorm.DB) error {
	var patients []entities.Patient
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	for i, patient := range patients {
		flg := &entities.Flg{
			IsCompleted:  false,
			Organization: "Городская поликлиника",
			Number:       10000 + i,
			Result:       "Норма",
			Date:         time.Now().Add(-time.Duration(i*24) * time.Hour),
		}

		if err := db.Create(flg).Error; err != nil {
			return fmt.Errorf("failed to create flg: %w", err)
		}

		// Обновляем пациента с FlgID
		patient.FlgID = &flg.ID
		if err := db.Save(&patient).Error; err != nil {
			return fmt.Errorf("failed to update patient with flg: %w", err)
		}
	}

	return nil
}

// seed/vaccines.go
func seedVaccines(db *gorm.DB) error {
	var patients []entities.Patient
	var titles []entities.Title
	var medications []entities.Medication
	var doses []entities.Dose
	var numbers []entities.Number

	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to get patients: %w", err)
	}

	if err := db.Find(&titles).Error; err != nil {
		return fmt.Errorf("failed to get titles: %w", err)
	}

	if err := db.Find(&medications).Error; err != nil {
		return fmt.Errorf("failed to get medications: %w", err)
	}

	if err := db.Find(&doses).Error; err != nil {
		return fmt.Errorf("failed to get doses: %w", err)
	}

	if err := db.Find(&numbers).Error; err != nil {
		return fmt.Errorf("failed to get numbers: %w", err)
	}

	for i, patient := range patients {
		vaccineCount := 2
		if i%2 == 0 {
			vaccineCount = 3
		}

		for j := 0; j < vaccineCount; j++ {
			var titleID, medicationID, doseID, numberID *uint

			if len(titles) > 0 {
				titleID = &titles[j%len(titles)].ID
			}
			if len(medications) > 0 {
				medicationID = &medications[j%len(medications)].ID
			}
			if len(doses) > 0 {
				doseID = &doses[j%len(doses)].ID
			}
			if len(numbers) > 0 {
				numberID = &numbers[j%len(numbers)].ID
			}

			vaccine := &entities.Vaccine{
				Date:            time.Now().Add(-time.Duration((i+j)*24) * time.Hour),
				IsCompleted:     j%3 != 0,
				IsRefusal:       j%4 == 0,
				IsExemption:     j%5 == 0,
				TiterAmount:     intPtr(100 + j*50),
				MedWithdrawlNum: intPtr(1000 + j*100),
				Result:          stringPtr(fmt.Sprintf("Результат %d", j+1)),
				TitleID:         titleID,
				MedicationID:    medicationID,
				DoseID:          doseID,
				NumberID:        numberID,
				PatientID:       patient.ID,
				CreatedAt:       time.Now().Add(-time.Duration((i+j)*24) * time.Hour),
			}

			if err := db.Create(vaccine).Error; err != nil {
				return fmt.Errorf("failed to create vaccine: %w", err)
			}
		}
	}

	return nil
}

// Вспомогательные функции
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
