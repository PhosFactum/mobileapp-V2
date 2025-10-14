package interfaces

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"gorm.io/gorm"
)

type Service interface {
	ParamsParserService
	FilterBuilderService
	TxManager
}

// ParamsParserService Сервис преобразования типов
// Парсинг строковых параметров и приведение к единому типу
type ParamsParserService interface {
	ParseDateString(dateStr string) (time.Time, error)
	ParseTimeString(timeStr string) (time.Time, error)
	ParseUintString(uintStr string) (uint, error)
	ParseIntString(intStr string) (int, error)
	ParseUint(value interface{}) (uint, error)

	FormatDateToString(t time.Time) string
	FormatTimeToString(t time.Time) string
	ConvertJSONSchemaToFields(schemaJSON json.RawMessage) ([]models.FieldDescriptor, error)
}

type FilterBuilderService interface {
	ParseFilterString(filterStr string, modelFields map[string]string) (string, []interface{}, error)
	ParseOrderString(orderStr string, modelFields map[string]string) (string, error)
}

type TxManager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetTransaction(ctx context.Context) *gorm.DB
}
