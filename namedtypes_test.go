package namedtypes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	namedtypes "github.com/gomatic/yze-namedtypes"
)

func TestBarePrimitiveParameterIsReported(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), namedtypes.Analyzer, "a")
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, namedtypes.Registration.Validate())
	assert.Equal(t, "yze/namedtypes", namedtypes.Registration.RuleID())
	assert.Same(t, namedtypes.Analyzer, namedtypes.Registration.Analyzer)
}
