package team

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sirupsen/logrus"
)

func TestHandler_GetTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(*mockService)
		expectedStatus int
		expectedError  string
		validateBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "successful get with multiple users",
			queryParams: "team_name=backend",
			setupMock: func(m *mockService) {
				m.On("GetTeam", mock.Anything, "backend").Return(teamsrv.Team{
					TeamName: "backend",
					Members: []teamsrv.TeamMember{
						{UserID: "user-001", Username: "alice", IsActive: true},
						{UserID: "user-002", Username: "bob", IsActive: true},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response GetTeamResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "backend", response.Team.TeamName)
				assert.Len(t, response.Team.Members, 2)
				assert.Equal(t, "user-001", response.Team.Members[0].UserID)
				assert.Equal(t, "alice", response.Team.Members[0].Username)
			},
		},
		{
			name:        "successful get with single user",
			queryParams: "team_name=frontend",
			setupMock: func(m *mockService) {
				m.On("GetTeam", mock.Anything, "frontend").Return(teamsrv.Team{
					TeamName: "frontend",
					Members: []teamsrv.TeamMember{
						{UserID: "user-003", Username: "charlie", IsActive: true},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response GetTeamResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "frontend", response.Team.TeamName)
				assert.Len(t, response.Team.Members, 1)
			},
		},
		{
			name:        "successful get with no users",
			queryParams: "team_name=empty-team",
			setupMock: func(m *mockService) {
				m.On("GetTeam", mock.Anything, "empty-team").Return(teamsrv.Team{
					TeamName: "empty-team",
					Members:  []teamsrv.TeamMember{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response GetTeamResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "empty-team", response.Team.TeamName)
				assert.Empty(t, response.Team.Members)
			},
		},
		{
			name:        "missing team_name parameter",
			queryParams: "",
			setupMock:   func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
				assert.Contains(t, w.Body.String(), "team_name parameter is required")
			},
		},
		{
			name:        "empty team_name parameter",
			queryParams: "team_name=",
			setupMock:   func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
				assert.Contains(t, w.Body.String(), "team_name parameter is required")
			},
		},
		{
			name:        "team not found",
			queryParams: "team_name=nonexistent",
			setupMock: func(m *mockService) {
				m.On("GetTeam", mock.Anything, "nonexistent").Return(teamsrv.Team{}, teamsrv.ErrTeamNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  string(models.NOTFOUND),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOTFOUND))
				assert.Contains(t, w.Body.String(), "team not found")
			},
		},
		{
			name:        "internal server error",
			queryParams: "team_name=backend",
			setupMock: func(m *mockService) {
				m.On("GetTeam", mock.Anything, "backend").Return(teamsrv.Team{}, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "INTERNAL_ERROR",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
				assert.Contains(t, w.Body.String(), "database connection failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			tt.setupMock(mockSvc)

			logger := logrus.New()
			logger.SetLevel(logrus.ErrorLevel)
			handler := &Handler{
				service: mockSvc,
				logger:  logger,
			}

			router := gin.New()
			router.GET("/team/get", handler.GetTeam)

			url := "/team/get"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.validateBody != nil {
				tt.validateBody(t, w)
			}

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

