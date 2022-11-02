package util

import (
	"github.com/googolplex-s6/kwekker-protobufs/v2/kwek"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"testing"
	"time"
)

func TestValidateCreateKwekWithValidKwek(t *testing.T) {
	createKwek := kwek.CreateKwek{
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

func TestValidateCreateKwekWithEmptyText(t *testing.T) {
	createKwek := kwek.CreateKwek{
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
		KwekId:    1,
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

func TestValidateUpdateKwekWithEmptyKwekId(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekId:    0,
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

func TestValidateUpdateKwekWithNegativeKwekId(t *testing.T) {
	updateKwek := kwek.UpdateKwek{
		KwekId:    -1,
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
		KwekId:    1,
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
		KwekId:    1,
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
		KwekId:    1,
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
		KwekId:    1,
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
		KwekId:    1,
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
		KwekId:    1,
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
		KwekId:    1,
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
		KwekId: 123,
	}

	validation := ValidateDeleteKwek(&deleteKwek)

	if !validation.Valid {
		t.Errorf("Validation should be valid, but is not")
	}

	if len(validation.Errors) > 0 {
		t.Errorf("Validation should have no errors, but has %d", len(validation.Errors))
	}
}

func TestValidateDeleteKwekWithEmptyKwekId(t *testing.T) {
	deleteKwek := kwek.DeleteKwek{
		KwekId: 0,
	}

	validation := ValidateDeleteKwek(&deleteKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}

func TestValidateDeleteKwekWithNegativeKwekId(t *testing.T) {
	deleteKwek := kwek.DeleteKwek{
		KwekId: -1,
	}

	validation := ValidateDeleteKwek(&deleteKwek)

	if validation.Valid {
		t.Errorf("Validation should be invalid, but is not")
	}

	if len(validation.Errors) != 1 {
		t.Errorf("Validation should have one error, but has %d", len(validation.Errors))
	}
}
