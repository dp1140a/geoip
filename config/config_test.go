package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (ts *ConfigTestSuite) SetupSuite() {}

func (ts *ConfigTestSuite) TearDownSuite() {}

func (ts *ConfigTestSuite) SetupTest() {
}

func (ts *ConfigTestSuite) TearDownTest() {}

//With File
func (ts *ConfigTestSuite) TestGoodConfigFile() {
	err := InitConfig("./testdata/config_good.toml")
	assert.NoError(ts.T(), err)
	assert.NotEmpty(ts.T(), viper.AllKeys())
}

//With no file
func (ts *ConfigTestSuite) TestNoConfigFile() {
	err := InitConfig("")
	assert.NoError(ts.T(), err)
	assert.NotEmpty(ts.T(), viper.AllKeys())
}

//With Bad File
func (ts *ConfigTestSuite) TestBadConfigFile() {
	err := InitConfig("./testdata/config_bad.toml")
	assert.Error(ts.T(), err)
	assert.Contains(ts.T(), err.Error(), "While parsing config")
}

//With ENV Vars
func (ts *ConfigTestSuite) TestEnvVarsConfig() {
	expected := "BAR"
	os.Setenv("GEOIP_FOO", expected)
	err := InitConfig("")
	got := viper.GetString("foo")
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), expected, got)
}

//Print Config
func (ts *ConfigTestSuite) TestPrintConfig() {
	err := InitConfig("")
	assert.NoError(ts.T(), err)
	PrintConfig()
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
