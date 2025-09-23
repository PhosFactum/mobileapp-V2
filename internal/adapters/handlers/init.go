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

	// Доктора
	doctorGroup := protected.Group("/doctors")
	doctorGroup.GET("/:doc_id", h.GetDoctorByID)
	doctorGroup.PUT("/:doc_id", h.UpdateDoctor)

	// // Приёмы больницы
	// hospitalGroup := protected.Group("/hospital")
	// hospitalGroup.GET("/receptions/patients/:pat_id", h.GetAllReceptionsByPatientID) // Все приемы пациента
	// hospitalGroup.GET("/receptions/:doc_id", h.GetReceptionsHospitalByDoctorID)      // Все приемы доктора
	// hospitalGroup.GET("/receptions/:doc_id/:hosp_id", h.GetReceptionHosptalById)
	// hospitalGroup.PUT("/receptions/:recep_id", h.UpdateReceptionHospitalByReceptionID)

	// Поправить пациента на транзакцию

	// Руты рабочие для новго проекта

	// Авторизация
	authGroup := baseRouter.Group("/auth")
	authGroup.POST("/login", h.LoginDoctor)
	authGroup.POST("/logout", jwtMiddleware.JWTAuth(cfg.JWTSecret), h.LogoutDoctor)

	// Организации (страховые)
	organizationGroup := protected.Group("/organization")
	organizationGroup.GET("/", h.GetAllOrganizations)

	//Списки пациентов
	patientGroupsGroup := protected.Group("/groups")
	patientGroupsGroup.GET("/", h.GetPatientGroupsByCodeOrOrgTitle) //arg search
	patientGroupsGroup.GET("/:org_id", h.GetPatientGroupsByOrganization)

	// Пациенты
	patientGroup := baseRouter.Group("/patients")
	patientGroup.GET("/:group_id", h.GetPatientsByGroup)
	patientGroup.POST("/:group_id/create", h.CreatePatient)

	consentGroup := protected.Group("/consent")
	consentGroup.GET("/personal-data", h.GetPersonalDataConsent)
	consentGroup.GET("/medical-exam", h.GetMedicalExamConsent)
	consentGroup.GET("/signature/:recep_id", h.GetSignature)
	consentGroup.POST("/signature/:recep_id", h.SaveSignature)

	return r
}
