package validation

import (
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"testing"
	"time"
)

func TestValidateCreateKwekWithValidKwek(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithEmptyGuid(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithInvalidGuid(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "invalid",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithEmptyText(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithTextWithMaxLength(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     strings.Repeat("a", 256),
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithTooLongText(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     strings.Repeat("a", 257),
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithEmptyUserId(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "Hello world!",
		UserId:   "",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithEmptyPostedAt(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: nil,
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithPostedAtInFuture(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix() + 100},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithPostedAtInTooDistantPast(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix() - 60*60*24*30},
	}

	validation := ValidateCreateKwek(&createKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateCreateKwekWithPostedAtInDistantButAcceptablePast(t *testing.T) {
	createKwek := kwek.CreateKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:     "Hello world!",
		UserId:   "123",
		PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix() - 60*60*24*29},
	}

	validation := ValidateCreateKwek(&createKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithValidKwek(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      "Hello world!",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithEmptyKwekGuid(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "",
		Text:      "Hello world!",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithInvalidKwekGuid(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "invalid",
		Text:      "Hello world!",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithEmptyText(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      "",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithTextWithMaxLength(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      strings.Repeat("a", 256),
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithTooLongText(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      strings.Repeat("a", 257),
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithEmptyUpdatedAt(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      "Hello world!",
		UpdatedAt: nil,
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithUpdatedAtInFuture(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      "Hello world!",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix() + 100},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithUpdatedAtInTooDistantPast(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      "Hello world!",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix() - 60*60*24*30},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateUpdateKwekWithUpdatedAtInDistantButAcceptablePast(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekGuid:  "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
		Text:      "Hello world!",
		UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix() - 60*60*24*29},
	}

	validation := ValidateUpdateKwek(&updateKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateDeleteKwekWithValidKwek(t *testing.T) {
	deleteKwek := kwek.DeleteKwek{
		KwekGuid: "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc",
	}

	validation := ValidateDeleteKwek(&deleteKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateDeleteKwekWithEmptyKwekGuid(t *testing.T) {
	deleteKwek := kwek.DeleteKwek{
		KwekGuid: "",
	}

	validation := ValidateDeleteKwek(&deleteKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateDeleteKwekWithInvalidKwekGuid(t *testing.T) {
	deleteKwek := kwek.DeleteKwek{
		KwekGuid: "invalid",
	}

	validation := ValidateDeleteKwek(&deleteKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}
