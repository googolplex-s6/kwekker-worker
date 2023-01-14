package validation

import (
	"fmt"
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	"github.com/googolplex-s6/kwekker-protobufs/v3/user"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Validation struct {
	Valid  bool
	Errors []string
}

func Validate(message proto.Message) Validation {
	switch message.(type) {
	case *kwek.CreateKwek:
		return ValidateCreateKwek(message.(*kwek.CreateKwek))
	case *kwek.UpdateKwek:
		return ValidateUpdateKwek(message.(*kwek.UpdateKwek))
	case *kwek.DeleteKwek:
		return ValidateDeleteKwek(message.(*kwek.DeleteKwek))
	case *user.CreateUser:
		return ValidateCreateUser(message.(*user.CreateUser))
	case *user.UpdateUser:
		return ValidateUpdateUser(message.(*user.UpdateUser))
	case *user.DeleteUser:
		return ValidateDeleteUser(message.(*user.DeleteUser))
	default:
		return Validation{
			Valid:  false,
			Errors: []string{"Unknown message type"},
		}
	}
}

func validateGuid(guid string, key string, v *Validation) {
	if !assertNotEmpty(guid, key, v) {
		return
	}

	if len(guid) != 36 {
		v.Valid = false
		v.Errors = append(v.Errors, fmt.Sprintf("%s must be a valid GUID", key))
	}
}

func validateUpdatedAt(updatedAt *timestamppb.Timestamp, v *Validation) {
	validateTimestamp(updatedAt, "UpdatedAt", v)
}

func validateTimestamp(timestamp *timestamppb.Timestamp, key string, v *Validation) {
	if timestamp == nil {
		v.Valid = false
		v.Errors = append(v.Errors, fmt.Sprintf("%s is required", key))

		return
	}

	if timestamp.AsTime().After(time.Now()) {
		v.Valid = false
		v.Errors = append(v.Errors, fmt.Sprintf("%s cannot be in the future", key))
	} else if timestamp.AsTime().Before(time.Now().Add(-24 * 30 * time.Hour)) {
		v.Valid = false
		v.Errors = append(v.Errors, fmt.Sprintf("%s cannot be more than 30 days ago", key))
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
