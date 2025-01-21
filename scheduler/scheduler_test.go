package scheduler_test

import (
	"testing"
	"time"

	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
	"github.com/natanhp/yangnder/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SchedulerTestSuite struct {
	suite.Suite
	DB *gorm.DB
}

func (suite *SchedulerTestSuite) SetupSuite() {
	ConnectTest()
	MigrateTest()
	PopulateUsersTest()

	config.DB = DBTest
	suite.DB = DBTest
}

func (suite *SchedulerTestSuite) TestResetSwipeNumber() {
	scheduler.ResetSwipeNumber()
	var user models.User
	suite.DB.First(&user, 1)

	assert.Equal(suite.T(), 10, user.SwipeNum)
	assert.Equal(suite.T(), true, user.NextSwipeReset.After(time.Now()))
}

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
