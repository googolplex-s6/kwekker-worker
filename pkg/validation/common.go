package validation

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Validation struct {
	Valid  bool
	Errors []string
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
