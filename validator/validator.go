package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/i18n"
)

// ValidateHuma 在 Resolver 中使用的验证函数
// 完全兼容 Huma 的验证标签，支持 message 标签自定义错误消息
func Validate(input interface{}) []error {
	// 使用默认 context
	return ValidateWithContext(context.Background(), input)
}

// ValidateHumaWithContext 带上下文的验证函数，支持多语言
func ValidateWithContext(ctx context.Context, input interface{}) []error {
	var errors []error

	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// 处理指针
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 遍历所有字段
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 获取字段名（优先使用 json 标签）
		fieldName := getFieldName(field)

		// 跳过未导出字段
		if !field.IsExported() {
			continue
		}

		// 验证该字段
		if fieldErrors := validateField(ctx, field, fieldValue, fieldName); len(fieldErrors) > 0 {
			errors = append(errors, fieldErrors...)
		}
	}

	return errors
}

// validateField 验证单个字段
func validateField(ctx context.Context, field reflect.StructField, value reflect.Value, fieldName string) []error {
	var errors []error

	// 获取自定义消息
	customMessage := field.Tag.Get("message")

	// 获取字段的显示名称（优先使用 doc 标签）
	displayName := field.Tag.Get("doc")
	if displayName == "" {
		displayName = fieldName
	}

	// 获取字段值的字符串表示
	strValue := getStringValue(value)

	// required 验证
	if tag := field.Tag.Get("required"); tag == "true" {
		if isZeroValue(value) {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.required", map[string]any{"Field": displayName})
			}
			errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
		}
	}

	// 如果字段为空且不是必需的，跳过其他验证
	if isZeroValue(value) && field.Tag.Get("required") != "true" {
		return errors
	}

	// minLength 验证
	if tag := field.Tag.Get("minLength"); tag != "" {
		minLen, _ := strconv.Atoi(tag)
		if len(strValue) < minLen {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.minLength", map[string]any{"Field": displayName, "Min": minLen})
			}
			errors = append(errors, createFieldError(fieldName, msg, strValue))
		}
	}

	// maxLength 验证
	if tag := field.Tag.Get("maxLength"); tag != "" {
		maxLen, _ := strconv.Atoi(tag)
		if len(strValue) > maxLen {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.maxLength", map[string]any{"Field": displayName, "Max": maxLen})
			}
			errors = append(errors, createFieldError(fieldName, msg, strValue))
		}
	}

	// minimum 验证（数字）
	if tag := field.Tag.Get("minimum"); tag != "" {
		if numValue, ok := getNumericValue(value); ok {
			min, _ := strconv.ParseFloat(tag, 64)
			if numValue < min {
				msg := customMessage
				if msg == "" {
					msg = i18n.T(ctx, "validator.minimum", map[string]any{"Field": displayName, "Min": tag})
				}
				errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
			}
		}
	}

	// maximum 验证（数字）
	if tag := field.Tag.Get("maximum"); tag != "" {
		if numValue, ok := getNumericValue(value); ok {
			max, _ := strconv.ParseFloat(tag, 64)
			if numValue > max {
				msg := customMessage
				if msg == "" {
					msg = i18n.T(ctx, "validator.maximum", map[string]any{"Field": displayName, "Max": tag})
				}
				errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
			}
		}
	}

	// exclusiveMinimum 验证
	if tag := field.Tag.Get("exclusiveMinimum"); tag != "" {
		if numValue, ok := getNumericValue(value); ok {
			min, _ := strconv.ParseFloat(tag, 64)
			if numValue <= min {
				msg := customMessage
				if msg == "" {
					msg = i18n.T(ctx, "validator.exclusiveMinimum", map[string]any{"Field": displayName, "Min": tag})
				}
				errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
			}
		}
	}

	// exclusiveMaximum 验证
	if tag := field.Tag.Get("exclusiveMaximum"); tag != "" {
		if numValue, ok := getNumericValue(value); ok {
			max, _ := strconv.ParseFloat(tag, 64)
			if numValue >= max {
				msg := customMessage
				if msg == "" {
					msg = i18n.T(ctx, "validator.exclusiveMaximum", map[string]any{"Field": displayName, "Max": tag})
				}
				errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
			}
		}
	}

	// multipleOf 验证
	if tag := field.Tag.Get("multipleOf"); tag != "" {
		if numValue, ok := getNumericValue(value); ok {
			multiple, _ := strconv.ParseFloat(tag, 64)
			if multiple != 0 {
				remainder := float64(int64(numValue/multiple) * int64(multiple))
				if numValue != remainder {
					msg := customMessage
					if msg == "" {
						msg = i18n.T(ctx, "validator.multipleOf", map[string]any{"Field": displayName, "Multiple": tag})
					}
					errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
				}
			}
		}
	}

	// pattern 验证
	if pattern := field.Tag.Get("pattern"); pattern != "" {
		matched, err := regexp.MatchString(pattern, strValue)
		if err != nil || !matched {
			msg := customMessage
			if msg == "" {
				// 检查是否有 patternDescription
				if desc := field.Tag.Get("patternDescription"); desc != "" {
					msg = desc
				} else {
					msg = i18n.T(ctx, "validator.pattern", map[string]any{"Field": displayName})
				}
			}
			errors = append(errors, createFieldError(fieldName, msg, strValue))
		}
	}

	// format 验证
	if format := field.Tag.Get("format"); format != "" {
		if !validateFormat(format, strValue) {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.format."+format, map[string]any{"Field": displayName})
			}
			errors = append(errors, createFieldError(fieldName, msg, strValue))
		}
	}

	// enum 验证
	if enum := field.Tag.Get("enum"); enum != "" {
		values := strings.Split(enum, ",")
		found := false
		for _, v := range values {
			if strings.TrimSpace(v) == strValue {
				found = true
				break
			}
		}
		if !found {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.enum", map[string]any{"Field": displayName, "Values": enum})
			}
			errors = append(errors, createFieldError(fieldName, msg, strValue))
		}
	}

	// minItems 验证（数组/切片）
	if tag := field.Tag.Get("minItems"); tag != "" && value.Kind() == reflect.Slice {
		minItems, _ := strconv.Atoi(tag)
		if value.Len() < minItems {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.minItems", map[string]any{"Field": displayName, "Min": minItems})
			}
			errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
		}
	}

	// maxItems 验证（数组/切片）
	if tag := field.Tag.Get("maxItems"); tag != "" && value.Kind() == reflect.Slice {
		maxItems, _ := strconv.Atoi(tag)
		if value.Len() > maxItems {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.maxItems", map[string]any{"Field": displayName, "Max": maxItems})
			}
			errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
		}
	}

	// uniqueItems 验证
	if tag := field.Tag.Get("uniqueItems"); tag == "true" && value.Kind() == reflect.Slice {
		if !hasUniqueItems(value) {
			msg := customMessage
			if msg == "" {
				msg = i18n.T(ctx, "validator.uniqueItems", map[string]any{"Field": displayName})
			}
			errors = append(errors, createFieldError(fieldName, msg, value.Interface()))
		}
	}

	return errors
}

// validateFormat 验证格式
func validateFormat(format, value string) bool {
	switch format {
	case "email":
		// 简单的邮箱验证
		pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "uri", "url":
		// URL验证
		pattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "uuid":
		// UUID验证
		pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "date":
		// 日期格式 YYYY-MM-DD
		pattern := `^\d{4}-\d{2}-\d{2}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "date-time":
		// ISO 8601 日期时间格式
		pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})?$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "time":
		// 时间格式 HH:MM:SS
		pattern := `^\d{2}:\d{2}:\d{2}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "ipv4":
		pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "ipv6":
		pattern := `^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|::1|::)$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	case "hostname":
		pattern := `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	default:
		// 未知格式，默认通过
		return true
	}
}

// hasUniqueItems 检查切片是否有唯一项
func hasUniqueItems(value reflect.Value) bool {
	seen := make(map[interface{}]bool)
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		if seen[item] {
			return false
		}
		seen[item] = true
	}
	return true
}

// getFieldName 获取字段名
func getFieldName(field reflect.StructField) string {
	// 优先使用 json 标签
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" && parts[0] != "-" {
			return parts[0]
		}
	}
	// 其次使用 path 标签
	if pathTag := field.Tag.Get("path"); pathTag != "" {
		return pathTag
	}
	// 其次使用 query 标签
	if queryTag := field.Tag.Get("query"); queryTag != "" {
		return queryTag
	}
	// 其次使用 header 标签
	if headerTag := field.Tag.Get("header"); headerTag != "" {
		return headerTag
	}
	// 最后使用字段名
	return field.Name
}

// getStringValue 获取字段的字符串值
func getStringValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	default:
		return fmt.Sprintf("%v", value.Interface())
	}
}

// getNumericValue 获取数字值
func getNumericValue(value reflect.Value) (float64, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()), true
	case reflect.Float32, reflect.Float64:
		return value.Float(), true
	default:
		return 0, false
	}
}

// isZeroValue 检查是否为零值
func isZeroValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
		return value.IsNil()
	default:
		return value.IsZero()
	}
}

// createFieldError 创建字段错误
func createFieldError(fieldName, message string, value interface{}) error {
	return &huma.ErrorDetail{
		Location: "body." + fieldName,
		Message:  message,
		Value:    value,
	}
}

// ValidateField 验证单个字段值（用于自定义验证）
func ValidateField(input interface{}, fieldName string, value interface{}, message string) error {
	// 尝试从结构体中获取 message 标签
	if message == "" {
		typ := reflect.TypeOf(input)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			name := getFieldName(field)
			if name == fieldName {
				if msg := field.Tag.Get("message"); msg != "" {
					message = msg
				}
				break
			}
		}
	}

	return createFieldError(fieldName, message, value)
}
