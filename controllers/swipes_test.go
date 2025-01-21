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
	DB     *gorm.DB
	Routes *gin.Engine
}

func (suite *SwipesControllerTestSuite) SetupSuite() {
	ConnectTest()
	MigrateTest()
	PopulateUsersTest()

	config.DB = DBTest
	suite.DB = DBTest

	gin.SetMode(gin.TestMode)
	router := gin.New()
	suite.Routes = router
}

func (suite *SwipesControllerTestSuite) TestRight() {
	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH("/swipes/right", controllers.Right)

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, "/swipes/right", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	user := responseBody["user"].(map[string]interface{})
	assert.Equal(suite.T(), float64(2), user["id"])
	assert.Equal(suite.T(), float64(1), data["user_id"].(float64))
	assert.Equal(suite.T(), float64(2), data["r_swipe_id"].(float64))
}

func (suite *SwipesControllerTestSuite) TestRightUserNotFound() {
	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH("/swipes/right", controllers.Right)

	payload := `{ "r_swipe_id": 1001 }`
	req, _ := http.NewRequest(http.MethodPatch, "/swipes/right", strings.NewReader(payload))
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
	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH("/swipes/right", controllers.Right)

	suite.DB.Model(&models.User{}).Where("id = ?", 2).Update("swipe_num", 0)

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, "/swipes/right", strings.NewReader(payload))
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
	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH("/swipes/right", controllers.Right)

	suite.DB.Create(&models.RSwipe{
		UserID:   1,
		RSwipeID: 2,
	})

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, "/swipes/right", strings.NewReader(payload))
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
	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH("/swipes/right", controllers.Right)

	suite.DB.Create(&models.LSwipe{
		UserID:   1,
		LSwipeID: 2,
		DeleteOn: time.Now(),
	})

	payload := `{ "r_swipe_id": 2 }`
	req, _ := http.NewRequest(http.MethodPatch, "/swipes/right", strings.NewReader(payload))
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
