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

type PremiumControllerTestSuite struct {
	suite.Suite
	DB      *gorm.DB
	Routes  *gin.Engine
	BuyPath string
}

func (suite *PremiumControllerTestSuite) SetupSuite() {
	ConnectTest()
	MigrateTest()
	PopulateUsersTest()

	config.DB = DBTest
	suite.DB = DBTest

	suite.BuyPath = "/premiums/buy"

	gin.SetMode(gin.TestMode)
	router := gin.New()
	suite.Routes = router
}

func (suite *PremiumControllerTestSuite) TestGetPremium() {
	claims := jwt.MapClaims{
		"sub": float64(1),
	}
	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.PATCH(suite.BuyPath, controllers.BuyPremium)

	req, _ := http.NewRequest(http.MethodPatch, suite.BuyPath, nil)
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), data["id"])
	assert.Equal(suite.T(), true, data["is_verified"].(bool))
}

func TestPremiumControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PremiumControllerTestSuite))
}
