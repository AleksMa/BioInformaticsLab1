package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EngineTestSuite struct {
	suite.Suite
}

func (s *EngineTestSuite) SetupTest() {

}

func (s *EngineTestSuite) TestInc() {
	appID := int64(2388722)

	s.counters.Inc(appID)
	s.counters.Inc(appID)

	s.Equal(int64(2), s.counters.diffCounter[appID])
}

func TestVkAudienceCountersSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}
