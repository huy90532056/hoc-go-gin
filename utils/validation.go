package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			root := strings.Split(e.Namespace(), ".")[0]

			rawPath := strings.TrimPrefix(e.Namespace(), root+".")

			parts := strings.Split(rawPath, ".")

			for i, part := range parts {
				if strings.Contains(part, "[") {
					idx := strings.Index(part, "[")
					base := camelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = camelToSnake(part)
				}
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s phai lon hon %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s phai nho hon %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s phai lon hon hoac bang %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s phai nho hon hoac bang $s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s phai la UUID hop le", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s phai la Slug hop le", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s phai nhieu hon %s ky tu", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s phai it hon %s ky tu", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s phai co gia tri lon hon %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s phai co gia tri be hon %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[fieldPath] = fmt.Sprintf("%s phai la 1 trong cac gia tri: %s", fieldPath, allowedValues)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s la bat buoc", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s chi duoc phep chu thuong, in hoa, so va khoang trang", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s phai dung dinh dang la email", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s phai dung dinh dang la nam, thang, ngay", fieldPath)
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[fieldPath] = fmt.Sprintf("%s phai la 1 trong cac gia tri: %s", fieldPath, allowedValues)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{"error": "Yeu cau khong hop le" + err.Error()}
}

func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minVal, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false
		}
		return fl.Field().Int() >= minVal
	})

	v.RegisterValidation("max_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		maxVal, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false
		}
		return fl.Field().Int() <= maxVal
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		filename := fl.Field().String()

		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}

		allowedExt := strings.Fields(allowedStr)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

		for _, allowed := range allowedExt {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}
		return false
	})

	return nil
}
