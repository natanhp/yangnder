package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserControllerTestSuite struct {
	suite.Suite
	DB          *gorm.DB
	Routes      *gin.Engine
	FindAllPath string
}

func (suite *UserControllerTestSuite) SetupSuite() {
	ConnectTest()
	MigrateTest()
	PopulateUsersTest()

	config.DB = DBTest
	suite.DB = DBTest

	gin.SetMode(gin.TestMode)
	router := gin.New()
	suite.Routes = router

	suite.FindAllPath = "/users/findAll"

	claims := jwt.MapClaims{
		"sub": float64(1),
	}

	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})
}

func (suite *UserControllerTestSuite) TestFindAll() {
	suite.Routes.PATCH(suite.FindAllPath, controllers.FindAll)

	req, _ := http.NewRequest(http.MethodPatch, suite.FindAllPath, nil)
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), 10, len(data))
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
