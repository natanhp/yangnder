package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/controllers"
	"github.com/natanhp/yangnder/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SwipesControllerTestSuite struct {
	suite.Suite
	DB        *gorm.DB
	Routes    *gin.Engine
	LeftPath  string
	RightPath string
}

func (suite *SwipesControllerTestSuite) SetupSuite() {
	ConnectTest()
	MigrateTest()
	PopulateUsersTest()

	config.DB = DBTest
	suite.DB = DBTest

	suite.LeftPath = "/swipes/left"
	suite.RightPath = "/swipes/right"

	gin.SetMode(gin.TestMode)
	router := gin.New()
	suite.Routes = router

	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH(suite.RightPath, controllers.Right)
	suite.Routes.PATCH(suite.LeftPath, controllers.Left)
}

func (suite *SwipesControllerTestSuite) TestRight() {
	suite.DB.Model(&models.User{}).Where("id = ?", 1).Update("swipe_num", 10)
	suite.deleteSwipes()

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.RightPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	user := responseBody["user"].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), user["id"])
	assert.Equal(suite.T(), float64(1), data["user_id"].(float64))
	assert.Equal(suite.T(), float64(2), data["r_swipe_id"].(float64))
}

func (suite *SwipesControllerTestSuite) deleteSwipes() {
	suite.DB.Where("user_id = ?", 1).Delete(&models.LSwipe{})
	suite.DB.Where("user_id = ?", 1).Delete(&models.RSwipe{})
}

func (suite *SwipesControllerTestSuite) TestRightUserNotFound() {
	payload := `{ "r_swipe_id": 1001 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.RightPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "User not found", data)
}

func (suite *SwipesControllerTestSuite) TestRightOutOfSwipe() {
	suite.DB.Model(&models.User{}).Where("id = ?", 1).Update("swipe_num", 0)
	suite.deleteSwipes()

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.RightPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Out of swipes", data)
}

func (suite *SwipesControllerTestSuite) TestRightAlreadyRSwiped() {
	suite.DB.Create(&models.RSwipe{
		UserID:   1,
		RSwipeID: 2,
	})

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.RightPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Already swiped", data)
}

func (suite *SwipesControllerTestSuite) TestRightAlreadyLSwiped() {
	suite.DB.Create(&models.LSwipe{
		UserID:   1,
		LSwipeID: 2,
		DeleteOn: time.Now(),
	})

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.RightPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Already swiped", data)
}

func (suite *SwipesControllerTestSuite) TestLeft() {
	payload := `{ "l_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.LeftPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	user := responseBody["user"].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), user["id"])
	assert.Equal(suite.T(), float64(1), data["user_id"].(float64))
	assert.Equal(suite.T(), float64(2), data["l_swipe_id"].(float64))
}

func (suite *SwipesControllerTestSuite) TestLeftUserNotFound() {
	payload := `{ "l_swipe_id": 1001 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.LeftPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "User not found", data)
}

func (suite *SwipesControllerTestSuite) TestLeftOutOfSwipe() {
	suite.DB.Model(&models.User{}).Where("id = ?", 1).Update("swipe_num", 0)
	suite.deleteSwipes()

	payload := `{ "l_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.LeftPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Out of swipes", data)
}

func (suite *SwipesControllerTestSuite) TestLeftAlreadyRSwiped() {
	suite.DB.Create(&models.RSwipe{
		UserID:   1,
		RSwipeID: 2,
	})

	payload := `{ "l_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.LeftPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Already swiped", data)
}

func (suite *SwipesControllerTestSuite) TestLeftAlreadyLSwiped() {
	suite.DB.Create(&models.LSwipe{
		UserID:   1,
		LSwipeID: 2,
		DeleteOn: time.Now(),
	})

	payload := `{ "l_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, suite.LeftPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Already swiped", data)
}

func TestSwipeControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SwipesControllerTestSuite))
}
