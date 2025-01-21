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

func (suite *SchedulerTestSuite) TestDeleteLSwipes() {
	for i := 0; i < 9; i++ {
		lSwipe := models.LSwipe{
			UserID:   uint(i) + 1,
			LSwipeID: uint(i) + 2,
			DeleteOn: time.Now().AddDate(0, 0, -1),
		}

		DBTest.Create(&lSwipe)
	}

	scheduler.DeleteLSwipes()
	var lSwipes []models.LSwipe
	suite.DB.Find(&lSwipes)

	assert.Equal(suite.T(), true, len(lSwipes) == 0)
}

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
