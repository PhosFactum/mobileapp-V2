package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

const (
	DATE_LAYOUT = "2006-01-02"
	TIME_LAYOUT = "15:04:05"
)

var (
	datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)    // YYYY-MM-DD
	timePattern = regexp.MustCompile(`^\d{2}:\d{2}(:\d{2})?$`) // HH:MM or HH:MM:SS
)

type ParamsParser struct {
}

func NewParamsParser() interfaces.ParamsParserService {
	return &ParamsParser{}
}

func (s *ParamsParser) ParseDateString(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(DATE_LAYOUT, strings.TrimSpace(dateStr))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected '%s': %v", DATE_LAYOUT, err)
	}
	return parsedDate, nil
}

func (s *ParamsParser) ParseTimeString(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(TIME_LAYOUT, strings.TrimSpace(timeStr))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format, expected '%s': %v", TIME_LAYOUT, err)
	}
	return parsedTime, nil
}

func (s *ParamsParser) ParseUintString(uintStr string) (uint, error) {
	uintStr = strings.TrimSpace(uintStr)
	if uintStr == "" {
		return 0, fmt.Errorf("empty string provided, expected unsigned integer")
	}

	value, err := strconv.ParseUint(uintStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint format, expected unsigned integer: %v", err)
	}

	return uint(value), nil
}

func (s *ParamsParser) ParseIntString(intStr string) (int, error) {
	intStr = strings.TrimSpace(intStr)
	if intStr == "" {
		return 0, fmt.Errorf("empty string provided, expected integer")
	}

	value, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int format, expected signed integer: %v", err)
	}

	return int(value), nil
}

func (s *ParamsParser) FormatDateToString(t time.Time) string {
	return t.Format(DATE_LAYOUT)
}

func (s *ParamsParser) FormatTimeToString(t time.Time) string {
	return t.Format(TIME_LAYOUT)
}

func (s *ParamsParser) ParseUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case string:
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint(id), nil
		}
	case uint:
		return v, nil
	case uint8, uint16, uint32, uint64:
		return uint(reflect.ValueOf(v).Uint()), nil
	case int, int8, int16, int32, int64:
		iv := reflect.ValueOf(v).Int()
		if iv >= 0 {
			return uint(iv), nil
		}
	case float32, float64:
		fv := reflect.ValueOf(v).Float()
		if fv >= 0 && fv == float64(uint(fv)) {
			return uint(fv), nil
		}
	case json.Number:
		if iv, err := v.Int64(); err == nil && iv >= 0 {
			return uint(iv), nil
		}
	}

	return 0, fmt.Errorf("cannot convert %v (type %T) to uint", value, value)
}

// ConvertJSONSchemaToFields преобразует JSON Schema в список FieldDescriptor
func (s *ParamsParser) ConvertJSONSchemaToFields(schemaJSON json.RawMessage) ([]models.FieldDescriptor, error) {
	if len(schemaJSON) == 0 {
		return []models.FieldDescriptor{}, nil
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return nil, fmt.Errorf("invalid schema JSON: %w", err)
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'properties' in schema")
	}

	requiredMap := make(map[string]bool)
	if reqs, ok := schema["required"].([]interface{}); ok {
		for _, r := range reqs {
			if name, ok := r.(string); ok {
				requiredMap[name] = true
			}
		}
	}

	var fields []models.FieldDescriptor
	for name, prop := range properties {
		// prop имеет тип interface{} — передаём как есть
		field, err := convertProperty(name, prop, requiredMap[name])
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %q: %w", name, err)
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func convertProperty(name string, p interface{}, required bool) (models.FieldDescriptor, error) {
	// Приводим тип ВНУТРИ функции
	prop, ok := p.(map[string]interface{})
	if !ok {
		return models.FieldDescriptor{}, fmt.Errorf("property %q is not a JSON object", name)
	}

	field := models.FieldDescriptor{
		Name:     name,
		Title:    name,
		Required: required,
	}

	// Type
	if t, ok := prop["type"].(string); ok {
		field.Type = t
	}

	// Description → Title
	if desc, ok := prop["description"].(string); ok {
		field.Title = desc
		field.Description = &desc
	}

	// Format
	if f, ok := prop["format"].(string); ok {
		field.Format = &f
	}

	// String ограничения
	if minLen, ok := prop["minLength"].(float64); ok {
		val := int(minLen)
		field.MinLength = &val
	}
	if maxLen, ok := prop["maxLength"].(float64); ok {
		val := int(maxLen)
		field.MaxLength = &val
	}
	if pattern, ok := prop["pattern"].(string); ok {
		field.Pattern = &pattern
	}

	// Number ограничения
	if min, ok := prop["minimum"].(float64); ok {
		field.Minimum = &min
	}
	if max, ok := prop["maximum"].(float64); ok {
		field.Maximum = &max
	}
	if exMin, ok := prop["exclusiveMinimum"].(float64); ok {
		field.ExclusiveMinimum = &exMin
	}
	if exMax, ok := prop["exclusiveMaximum"].(float64); ok {
		field.ExclusiveMaximum = &exMax
	}
	if mult, ok := prop["multipleOf"].(float64); ok {
		field.MultipleOf = &mult
	}

	// Array
	if minItems, ok := prop["minItems"].(float64); ok {
		val := int(minItems)
		field.MinItems = &val
	}
	if maxItems, ok := prop["maxItems"].(float64); ok {
		val := int(maxItems)
		field.MaxItems = &val
	}

	// Enum
	if enums, ok := prop["enum"].([]interface{}); ok {
		for _, e := range enums {
			if s, ok := e.(string); ok {
				field.Enum = append(field.Enum, s)
			}
		}
	}

	// Tag
	field.Tag = inferTag(field.Type, prop, field.Enum != nil, field.Pattern != nil)

	return field, nil
}

func inferTag(fieldType string, prop map[string]interface{}, hasEnum bool, hasPattern bool) string {
	switch fieldType {
	case "string":
		if hasEnum {
			return "select"
		}
		if hasPattern {
			return "input" // или "masked-input", если нужна маска
		}
		// Проверяем format для специальных input-ов
		if format, ok := prop["format"].(string); ok {
			switch format {
			case "email":
				return "email"
			case "date":
				return "date"
			case "textarea": // кастомный format для многострочного ввода
				return "textarea"
			default:
				return "input"
			}
		}
		return "input"

	case "number", "integer":
		return "number"

	case "boolean":
		return "checkbox"

	case "array":
		// Если массив строк с enum — мультиселект
		if items, ok := prop["items"].(map[string]interface{}); ok {
			if itemType, ok := items["type"].(string); ok && itemType == "string" {
				if _, hasEnum := items["enum"]; hasEnum {
					return "multiselect"
				}
			}
		}
		return "list" // общий тег для массивов

	default:
		return "input"
	}
}
