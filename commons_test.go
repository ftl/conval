package conval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNoMember(t *testing.T) {
	validator := PropertyValidators[NoMemberProperty]
	assert.NoError(t, validator.ValidateProperty("nm", nil), "nm")
	assert.NoError(t, validator.ValidateProperty("NM", nil), "NM")
	assert.Error(t, validator.ValidateProperty("", nil), "empty")
	assert.Error(t, validator.ValidateProperty("   ", nil), "whitespace")
	assert.Error(t, validator.ValidateProperty("12345", nil), "number")
}
