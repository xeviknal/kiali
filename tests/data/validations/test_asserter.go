package validations

import (
	"testing"

	"github.com/kiali/kiali/models"
	"github.com/stretchr/testify/assert"
)

type IstioCheckTestAsserter struct {
	T           *testing.T
	Validations []*models.IstioCheck
	Valid       bool
}

func (tb IstioCheckTestAsserter) AssertNoValidations() {
	assert := assert.New(tb.T)

	assert.Empty(tb.Validations)
	assert.True(tb.Valid)
}

func (tb IstioCheckTestAsserter) AssertValidationsPresent(count int) {
	assert := assert.New(tb.T)

	assert.False(tb.Valid)
	assert.NotEmpty(tb.Validations)
	assert.Len(tb.Validations, count)
}

func (tb IstioCheckTestAsserter) AssertValidationAt(i int, severity models.SeverityLevel, path, message string) {
	assert := assert.New(tb.T)

	if len(tb.Validations) < i {
		tb.T.Error("Wrong memory access to validations array")
	}

	validation := tb.Validations[i]
	assert.NotNil(validation)
	assert.Equal(severity, validation.Severity)
	assert.Equal(path, validation.Path)
	assert.Equal(models.CheckMessage(message), validation.Message)
}

type ValidationsTestAsserter struct {
	T           *testing.T
	Validations models.IstioValidations
}

func (vta ValidationsTestAsserter) AssertNoValidations() {
	assert := assert.New(vta.T)

	assert.Empty(vta.Validations)
}

func (vta ValidationsTestAsserter) AssertValidationsPresent(count int) {
	assert := assert.New(vta.T)
	assert.NotEmpty(vta.Validations)
	assert.Len(vta.Validations, count)
}

func (vta ValidationsTestAsserter) AssertValidationAt(key models.IstioValidationKey, severity models.SeverityLevel, path, message string) {
	assert := assert.New(vta.T)

	// Assert specific's object validation
	validation, ok := vta.Validations[key]
	assert.True(ok)
	if validation == nil {
		return
	}

	assert.False(validation.Valid)

	// Assert object's checks
	assert.NotEmpty(validation.Checks)
	assert.Equal(severity, validation.Checks[0].Severity)
	assert.Equal(path, validation.Checks[0].Path)
	assert.Equal(models.CheckMessage(message), validation.Checks[0].Message)
}
