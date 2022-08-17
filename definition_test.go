package conval

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYAML_ValidNameIdentifierOfficialRules(t *testing.T) {
	yaml := bytes.NewBufferString(`name: Test Contest
identifier: TEST-CONTEST-VALID
official_rules: https://github.com/ftl/conval/testdata
`)
	expected := Definition{
		Name:          "Test Contest",
		Identifier:    "TEST-CONTEST-VALID",
		OfficialRules: "https://github.com/ftl/conval/testdata",
	}
	actual, err := LoadYAML(yaml)
	assert.NoError(t, err)

	assert.Equal(t, expected, *actual)
}
