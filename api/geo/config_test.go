package geo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (ts *ConfigTestSuite) SetupSuite() {}

func (ts *ConfigTestSuite) TearDownSuite() {}

func (ts *ConfigTestSuite) SetupTest() {}

func (ts *ConfigTestSuite) TearDownTest() {}

//With no file
func (ts *ConfigTestSuite) TestNoConfigFile() {
	config, err := InitConfig()
	fmt.Println(config.RefreshDuration.String())
	assert.NoError(ts.T(), err)
	assert.NotEmpty(ts.T(), config)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
