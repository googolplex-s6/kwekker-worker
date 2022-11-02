package validation

import (
	"github.com/googolplex-s6/kwekker-protobufs/v2/kwek"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Validation struct {
	Valid  bool
	Errors []string
}

func ValidateCreateKwek(kwek *kwek.CreateKwek) Validation {
	validation := Validation{
		Valid: true,
	}

	validateText(kwek.GetText(), &validation)
	validateUserId(kwek.GetUserId(), &validation)
	validatePostedAt(kwek.GetPostedAt(), &validation)

	return validation
}

func ValidateUpdateKwek(kwek *kwek.UpdateKwek) Validation {
	validation := Validation{
		Valid: true,
	}

	validateKwekId(kwek.GetKwekId(), &validation)
	validateText(kwek.GetText(), &validation)
	validateUpdatedAt(kwek.GetUpdatedAt(), &validation)

	return validation
}

func ValidateDeleteKwek(kwek *kwek.DeleteKwek) Validation {
	validation := Validation{
		Valid: true,
	}

	validateKwekId(kwek.GetKwekId(), &validation)

	return validation
}

func validateText(text string, v *Validation) {
	if !assertNotEmpty(text, "Text", v) {
		return
	}

	if len(text) > 256 {
		v.Valid = false
		v.Errors = append(v.Errors, "Text must be less than 256 characters")
	}
}

func validateUserId(userId string, v *Validation) {
	assertNotEmpty(userId, "UserId", v)
}

func validatePostedAt(postedAt *timestamppb.Timestamp, v *Validation) {
	timestampValidation := Validation{
		Valid: true,
	}

	validateTimestamp(postedAt, &timestampValidation)

	if !timestampValidation.Valid {
		v.Valid = false

		for _, validationErr := range timestampValidation.Errors {
			v.Errors = append(v.Errors, "PostedAt "+validationErr)
		}
	}
}

func validateUpdatedAt(updatedAt *timestamppb.Timestamp, v *Validation) {
	timestampValidation := Validation{
		Valid: true,
	}

	validateTimestamp(updatedAt, &timestampValidation)

	if !timestampValidation.Valid {
		v.Valid = false

		for _, validationErr := range timestampValidation.Errors {
			v.Errors = append(v.Errors, "UpdatedAt "+validationErr)
		}
	}
}

func validateTimestamp(timestamp *timestamppb.Timestamp, v *Validation) {
	if timestamp == nil {
		v.Valid = false
		v.Errors = append(v.Errors, "Timestamp is required")

		return
	}

	if timestamp.AsTime().After(time.Now()) {
		v.Valid = false
		v.Errors = append(v.Errors, "Timestamp cannot be in the future")
	}

	if timestamp.AsTime().Before(time.Now().Add(-24 * 30 * time.Hour)) {
		v.Valid = false
		v.Errors = append(v.Errors, "Timestamp cannot be more than 30 days ago")
	}
}

func validateKwekId(id int64, v *Validation) {
	if id == 0 {
		v.Valid = false
		v.Errors = append(v.Errors, "KwekId is required")
	}

	if id < 0 {
		v.Valid = false
		v.Errors = append(v.Errors, "KwekId cannot be negative")
	}
}

func assertNotEmpty(value string, key string, v *Validation) bool {
	if value == "" {
		v.Valid = false
		v.Errors = append(v.Errors, key+" is required")

		return false
	}

	return true
}
