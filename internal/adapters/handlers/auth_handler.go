// internal/adapters/handlers/auth.go
package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// LoginDoctor аутентифицирует врача
// @Summary Вход в систему
// @Description Аутентифицирует врача по номеру телефона и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.DoctorLoginRequest true "Данные для входа"
// @Success 200 {object} models.DoctorAuthResponse "Успешное создание"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Неверные учётные данные"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (h *Handler) LoginDoctor(c *gin.Context) {
	var req models.DoctorLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Error decoding request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	h.logger.Info("Auth attempt", "phone", req.Phone)

	// Передаём всю структуру req, а не отдельные поля
	id, token, err := h.usecase.LoginDoctor(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Auth failed", "error", err)
		// Возвращаем точное сообщение из ошибки, если это безопасно
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid request payload"})
		return
	}

	c.JSON(http.StatusOK, models.DoctorAuthResponse{
		ID:    id,
		Token: token,
	})
}

// LogoutDoctor осуществляет выход из системы
// @Summary Выход из системы
// @Description Удаляет токен на клиенте
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ResultResponse "Успешный выход"
// @Router /auth/logout [post]
func (h *Handler) LogoutDoctor(c *gin.Context) {
	// Просто логируем факт выхода и возвращаем успех
	// Клиент сам должен удалить токен из localStorage/cookies

	h.logger.Info("User successfully logged out")

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
		"action":  "Please remove the token from client storage",
	})
}

// С логаутом такая тема, что это ответственность фронтендера - удалить токен пользователя
// При тесте в сваггере или постмане логаута не будет - ибо это Stateless JWT

// Ну и в целом токен в JWT хранится на клиенте, потому только клиентская часть может его удалить
