package validation

import (
	"github.com/googolplex-s6/kwekker-protobufs/v3/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/mail"
	"strings"
)

func ValidateCreateUser(user *user.CreateUser) Validation {
	validation := Validation{
		Valid: true,
	}

	validateUserId(user.GetUserId(), &validation)
	validateUsername(user.GetUsername(), &validation)
	validateEmail(user.GetEmail(), &validation)
	validateDisplayName(user.GetDisplayName(), &validation)
	validateAvatarUrl(user.GetAvatarUrl(), &validation)
	validateCreatedAt(user.GetCreatedAt(), &validation)

	return validation
}

func ValidateUpdateUser(user *user.UpdateUser) Validation {
	validation := Validation{
		Valid: true,
	}

	validateUserId(user.GetUserId(), &validation)
	validateUpdatedAt(user.GetUpdatedAt(), &validation)

	if user.GetUsername() != "" {
		validateUsername(user.GetUsername(), &validation)
	}

	if user.GetEmail() != "" {
		validateEmail(user.GetEmail(), &validation)
	}

	if user.GetDisplayName() != "" {
		validateDisplayName(user.GetDisplayName(), &validation)
	}

	if user.GetAvatarUrl() != "" {
		validateAvatarUrl(user.GetAvatarUrl(), &validation)
	}

	return validation
}

func ValidateDeleteUser(user *user.DeleteUser) Validation {
	validation := Validation{
		Valid: true,
	}

	validateUserId(user.GetUserId(), &validation)

	return validation
}

func validateUserId(userId string, v *Validation) {
	if !assertNotEmpty(userId, "UserId", v) {
		return
	}
}

func validateUsername(username string, v *Validation) {
	if !assertNotEmpty(username, "Username", v) {
		return
	}

	if len(username) < 3 {
		v.Valid = false
		v.Errors = append(v.Errors, "Username must be at least 3 characters")
	} else if len(username) > 15 {
		v.Valid = false
		v.Errors = append(v.Errors, "Username must be less than 15 characters")
	}
}

func validateEmail(email string, v *Validation) {
	if !assertNotEmpty(email, "Email", v) {
		return
	}

	_, err := mail.ParseAddress(email)

	if err != nil {
		v.Valid = false
		v.Errors = append(v.Errors, "Email is not valid")
	}
}

func validateDisplayName(name string, v *Validation) {
	if !assertNotEmpty(name, "DisplayName", v) {
		return
	}

	if len(name) > 30 {
		v.Valid = false
		v.Errors = append(v.Errors, "DisplayName must be less than 30 characters")
	}
}

func validateAvatarUrl(url string, v *Validation) {
	if !assertNotEmpty(url, "AvatarUrl", v) {
		return
	}

	if len(url) > 256 {
		v.Valid = false
		v.Errors = append(v.Errors, "AvatarUrl must be less than 256 characters")
	}

	if !strings.HasPrefix(url, "https://") {
		v.Valid = false
		v.Errors = append(v.Errors, "AvatarUrl must start with https://")
	}
}

func validateCreatedAt(createdAt *timestamppb.Timestamp, v *Validation) {
	validateTimestamp(createdAt, "CreatedAt", v)
}
