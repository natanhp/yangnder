package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/controllers"
	"github.com/natanhp/yangnder/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserControllerTestSuite struct {
	suite.Suite
	DB              *gorm.DB
	Routes          *gin.Engine
	FindAllPath     string
	FineOnePath     string
	CreatePath      string
	UploadPhotoPath string
	LoginPath       string
}

func (suite *UserControllerTestSuite) SetupSuite() {
	ConnectTest()
	MigrateTest()
	PopulateUsersTest()
	os.Setenv("JWT_SECRET", "secret")

	config.DB = DBTest
	suite.DB = DBTest

	gin.SetMode(gin.TestMode)
	router := gin.New()
	suite.Routes = router

	suite.FindAllPath = "/users"
	suite.FineOnePath = "/users/detail/1"
	suite.CreatePath = "/users/register"
	suite.UploadPhotoPath = "/users/upload-photo"
	suite.LoginPath = "/users/login"

	claims := jwt.MapClaims{
		"sub": float64(1),
	}

	suite.Routes.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})

	suite.Routes.GET(suite.FindAllPath, controllers.FindAll)
	suite.Routes.GET(suite.FineOnePath, controllers.FindOne)
	suite.Routes.POST(suite.CreatePath, controllers.Create)
	suite.Routes.POST(suite.UploadPhotoPath, controllers.UploadPhoto)
	suite.Routes.POST(suite.LoginPath, controllers.Login)
}

func (suite *UserControllerTestSuite) TestFindAll() {
	req, _ := http.NewRequest(http.MethodGet, suite.FindAllPath, nil)
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), true, len(data) > 10)
}

func (suite *UserControllerTestSuite) TestFindOne() {
	req, _ := http.NewRequest(http.MethodGet, suite.FineOnePath, nil)
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), data["id"])
}

func (suite *UserControllerTestSuite) TestCreate() {
	payload := `
		{
			"email": "asdas3@asdsad.com",
			"password": "asd123",
			"name": "Adadasd",
			"desc": "Lorem ipsum",
			"dob": "2020-01-01"
		}
	`
	req, _ := http.NewRequest(http.MethodPost, suite.CreatePath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	assert.Equal(suite.T(), "asdas3@asdsad.com", data["email"])
}

func (suite *UserControllerTestSuite) TestCreateEmailTaken() {
	payload := createNewUser(suite)

	req, _ := http.NewRequest(http.MethodPost, suite.CreatePath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Email already taken", data)
}

func createNewUser(suite *UserControllerTestSuite) string {
	payload := `
		{
			"email": "asdas3@asdsad.com",
			"password": "asd123",
			"name": "Adadasd",
			"desc": "Lorem ipsum",
			"dob": "2020-01-01"
		}
	`

	createUser := models.User{}
	json.Unmarshal([]byte(payload), &createUser)
	createUser.Password, _ = argon2id.CreateHash(createUser.Password, argon2id.DefaultParams)
	suite.DB.Create(&createUser)
	return payload
}

func (suite *UserControllerTestSuite) TestUploadPhoto() {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile("photo", "test.jpg")
	assert.NoError(suite.T(), err)

	_, err = io.Copy(fileWriter, bytes.NewReader([]byte("dummy file content")))
	assert.NoError(suite.T(), err)

	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, suite.UploadPhotoPath, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	assert.NotEmpty(suite.T(), data["photo"])
}

func (suite *UserControllerTestSuite) TestUploadPhotoEmpty() {
	req, _ := http.NewRequest(http.MethodPost, suite.UploadPhotoPath, nil)
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Photo is required", data)
}

func (suite *UserControllerTestSuite) TestLogin() {
	payload := createNewUser(suite)
	req, _ := http.NewRequest(http.MethodPost, suite.LoginPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["data"].(map[string]interface{})
	token := responseBody["token"].(string)
	assert.Equal(suite.T(), "asdas3@asdsad.com", data["email"])
	assert.NotEmpty(suite.T(), token)
}

func (suite *UserControllerTestSuite) TestLoginInvalidEmail() {
	payload := `{ "email": "gjls@gibberish.com" }`
	req, _ := http.NewRequest(http.MethodPost, suite.LoginPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Invalid email or password", data)
}

func (suite *UserControllerTestSuite) TestLoginInvalidHash() {
	createNewUser(suite)
	payloadWrongPass := `{ "email": "asdas3@asdsad.com", "password": "wrongpass" }`

	req, _ := http.NewRequest(http.MethodPost, suite.LoginPath, strings.NewReader(payloadWrongPass))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.Routes.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	data := responseBody["error"].(string)
	assert.Equal(suite.T(), "Invalid email or password", data)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
