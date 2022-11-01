package conval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNoMember(t *testing.T) {
	validator := PropertyValidators[NoMemberProperty]
	assert.NoError(t, validator.ValidateProperty("nm"), "nm")
	assert.NoError(t, validator.ValidateProperty("NM"), "NM")
	assert.Error(t, validator.ValidateProperty(""), "empty")
	assert.Error(t, validator.ValidateProperty("   "), "whitespace")
	assert.Error(t, validator.ValidateProperty("12345"), "number")
}
