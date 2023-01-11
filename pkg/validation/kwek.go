package validation

import (
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ValidateCreateKwek(kwek *kwek.CreateKwek) Validation {
	validation := Validation{
		Valid: true,
	}

	validateKwekGuid(kwek.GetKwekGuid(), &validation)
	validateText(kwek.GetText(), &validation)
	validateUserId(kwek.GetUserId(), &validation)
	validatePostedAt(kwek.GetPostedAt(), &validation)

	return validation
}

func ValidateUpdateKwek(kwek *kwek.UpdateKwek) Validation {
	validation := Validation{
		Valid: true,
	}

	validateKwekGuid(kwek.GetKwekGuid(), &validation)
	validateText(kwek.GetText(), &validation)
	validateUpdatedAt(kwek.GetUpdatedAt(), &validation)

	return validation
}

func ValidateDeleteKwek(kwek *kwek.DeleteKwek) Validation {
	validation := Validation{
		Valid: true,
	}

	validateKwekGuid(kwek.GetKwekGuid(), &validation)

	return validation
}

func validateKwekGuid(kwekId string, v *Validation) {
	validateGuid(kwekId, "KwekGuid", v)
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

func validatePostedAt(postedAt *timestamppb.Timestamp, v *Validation) {
	validateTimestamp(postedAt, "PostedAt", v)
}
