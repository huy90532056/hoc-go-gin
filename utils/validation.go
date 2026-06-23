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
			switch e.Tag() {
			case "gt":
				errors[e.Field()] = e.Field() + " phai lon hon gia tri toi thieu"
			case "lt":
				errors[e.Field()] = e.Field() + " phai lon hon gia tri toi thieu"
			case "gte":
				errors[e.Field()] = e.Field() + " phai lon hon hoac bang gia tri toi thieu"
			case "lte":
				errors[e.Field()] = e.Field() + " phai nho hon hoac bang gia tri toi thieu"
			case "uuid":
				errors[e.Field()] = e.Field() + " phai la UUID hop le"
			case "slug":
				errors[e.Field()] = e.Field() + " phai la Slug hop le"
			case "min":
				errors[e.Field()] = fmt.Sprintf("%s phai nhieu hon %s ky tu", e.Field(), e.Param())
			case "max":
				errors[e.Field()] = fmt.Sprintf("%s phai it hon %s ky tu", e.Field(), e.Param())
			case "min_int":
				errors[e.Field()] = fmt.Sprintf("%s phai co gia tri lon hon %s", e.Field(), e.Param())
			case "max_int":
				errors[e.Field()] = fmt.Sprintf("%s phai co gia tri be hon %s", e.Field(), e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[e.Field()] = fmt.Sprintf("%s phai la 1 trong cac gia tri: %s", e.Field(), allowedValues)
			case "required":
				errors[e.Field()] = e.Field() + " la bat buoc"
			case "search":
				errors[e.Field()] = e.Field() + " chi duoc phep chu thuong, in hoa, so va khoang trang"
			case "email":
				errors[e.Field()] = e.Field() + " phai dung dinh dang la email"
			case "datetime":
				errors[e.Field()] = e.Field() + " phai dung dinh dang la nam, thang, ngay"
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[e.Field()] = fmt.Sprintf("%s phai la 1 trong cac gia tri: %s", e.Field(), allowedValues)
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
