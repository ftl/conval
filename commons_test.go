package conval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNoMember(t *testing.T) {
	validator := PropertyValidators[NoMemberProperty]
	assert.NoError(t, validator.ValidateProperty("nm"))
	assert.NoError(t, validator.ValidateProperty("NM"))
	assert.NoError(t, validator.ValidateProperty(""))
	assert.NoError(t, validator.ValidateProperty("   "))
	assert.Error(t, validator.ValidateProperty("12345"))
}
