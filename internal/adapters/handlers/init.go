package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	jwtMiddleware "github.com/AlexanderMorozov1919/mobileapp/internal/middleware/jwt"
	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/logging"
	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/swagger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/swaggo/files"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	logger  *logging.Logger
	usecase interfaces.Usecases
	service interfaces.Service
}

// NewHandler создает новый экземпляр Handler со всеми зависимостями
func NewHandler(usecase interfaces.Usecases, parentLogger *logging.Logger, service interfaces.Service) *Handler {
	handlerLogger := parentLogger.WithPrefix("HANDLER")
	handlerLogger.Info("Handler initialized",
		"component", "GENERAL",
	)
	return &Handler{
		logger:  handlerLogger,
		usecase: usecase,
		service: service,
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler, cfg *config.Config, swagCfg *swagger.Config) http.Handler {
	r := gin.Default()

	// CORS
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     cfg.Server.AllowedOrigins,
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	// Swagger-роутер
	swagger.Setup(r, swagCfg)

	// Logger
	r.Use(LoggingMiddleware(h.logger))

	// Общая группа для API
	baseRouter := r.Group("/api/v1")

	//Версия
	baseRouter.GET("/version", h.GetVersionProject)

	protected := baseRouter.Group("/")
	protected.Use(jwtMiddleware.JWTAuth(cfg.JWTSecret))

	// Руты рабочие для новго проекта

	// Авторизация
	authGroup := baseRouter.Group("/auth")
	authGroup.POST("/login", h.LoginDoctor)
	authGroup.POST("/logout", jwtMiddleware.JWTAuth(cfg.JWTSecret), h.LogoutDoctor)

	// Организации
	organizationGroup := protected.Group("/organizations")
	organizationGroup.GET("", h.GetAllDoctorOrganizations)

	// Группы пациентов
	patientGroupsGroup := protected.Group("/patient-groups")
	patientGroupsGroup.GET("/by-organization/:organization_id", h.GetPatientGroupsByOrganizationID)
	patientGroupsGroup.GET("/:group_id/patients", h.GetPatientsByGroup)

	// Пациенты
	patientGroup := baseRouter.Group("/patients")
	patientGroup.POST("", h.CreatePatient)

	// Флюрографии
	flgGroup := baseRouter.Group("/flg")
	flgGroup.POST("", h.CreateFlgWithPhoto)

	// Справочники
	manualGroup := baseRouter.Group("/manuals")
	manualGroup.GET("", h.GetAllManuals)

	// Приемы
	receptionGroup := baseRouter.Group("/receptions")
	receptionGroup.POST("/update", h.UpdateReceptionData)

	// Анализы ← исправлена опечатка!
	analysisGroup := baseRouter.Group("/analysis")
	analysisGroup.POST("/update", h.UpdateAnalysisOrder)

	// Вакцины
	vaccineGroup := baseRouter.Group("/vaccines")
	vaccineGroup.POST("", h.CreateVaccine)
	vaccineGroup.POST("/refusals", h.CreateVaccineRefusal)
	vaccineGroup.POST("/withdrawals", h.CreateVaccineWithdrawal)
	vaccineGroup.POST("/titrs", h.CreateTitr)

	// Доктора
	doctorGroup := protected.Group("/doctors")
	doctorGroup.GET("/current", h.GetDoctorByID)

	// Согласия
	consentGroup := protected.Group("/consent")
	consentGroup.GET("/personal-data", h.GetPersonalDataConsent)
	consentGroup.GET("/medical-exam", h.GetMedicalExamConsent)
	consentGroup.GET("/signature/:recep_id", h.GetSignature)
	consentGroup.POST("/signature/:recep_id", h.SaveSignature)

	return r
}
