package team

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	teamsrv "github.com/aabbuukkaarr8/PRService/internal/service/team"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CreateTeam(ctx context.Context, team teamsrv.Team) (teamsrv.Team, error) {
	args := m.Called(ctx, team)
	if args.Get(0) == nil {
		return teamsrv.Team{}, args.Error(1)
	}
	return args.Get(0).(teamsrv.Team), args.Error(1)
}

func (m *mockService) GetTeam(ctx context.Context, teamName string) (teamsrv.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return teamsrv.Team{}, args.Error(1)
	}
	return args.Get(0).(teamsrv.Team), args.Error(1)
}

func TestHandler_CreateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mockService)
		expectedStatus int
		expectedError  string
		validateBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful creation",
			requestBody: CreateTeamRequest{
				TeamName: "backend",
				Members: []MemberTeam{
					{UserID: "user-001", Username: "alice", IsActive: true},
					{UserID: "user-002", Username: "bob", IsActive: true},
				},
			},
			setupMock: func(m *mockService) {
				m.On("CreateTeam", mock.Anything, teamsrv.Team{
					TeamName: "backend",
					Members: []teamsrv.TeamMember{
						{UserID: "user-001", Username: "alice", IsActive: true},
						{UserID: "user-002", Username: "bob", IsActive: true},
					},
				}).Return(teamsrv.Team{
					TeamName: "backend",
					Members: []teamsrv.TeamMember{
						{UserID: "user-001", Username: "alice", IsActive: true},
						{UserID: "user-002", Username: "bob", IsActive: true},
					},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response CreateTeamResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "backend", response.Team.TeamName)
				assert.Len(t, response.Team.Members, 2)
				assert.Equal(t, "user-001", response.Team.Members[0].UserID)
				assert.Equal(t, "alice", response.Team.Members[0].Username)
			},
		},
		{
			name: "invalid JSON - missing team_name",
			requestBody: map[string]interface{}{
				"members": []map[string]interface{}{
					{"user_id": "user-001", "username": "alice", "is_active": true},
				},
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "invalid JSON - missing members",
			requestBody: map[string]interface{}{
				"team_name": "backend",
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "invalid JSON - missing member fields",
			requestBody: map[string]interface{}{
				"team_name": "backend",
				"members": []map[string]interface{}{
					{"user_id": "user-001"},
				},
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "team already exists",
			requestBody: CreateTeamRequest{
				TeamName: "backend",
				Members: []MemberTeam{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockService) {
				m.On("CreateTeam", mock.Anything, teamsrv.Team{
					TeamName: "backend",
					Members: []teamsrv.TeamMember{
						{UserID: "user-001", Username: "alice", IsActive: true},
					},
				}).Return(teamsrv.Team{}, teamsrv.ErrTeamExists)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  string(models.TEAMEXISTS),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.TEAMEXISTS))
				assert.Contains(t, w.Body.String(), "team_name already exists")
			},
		},
		{
			name: "internal server error",
			requestBody: CreateTeamRequest{
				TeamName: "backend",
				Members: []MemberTeam{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockService) {
				m.On("CreateTeam", mock.Anything, teamsrv.Team{
					TeamName: "backend",
					Members: []teamsrv.TeamMember{
						{UserID: "user-001", Username: "alice", IsActive: true},
					},
				}).Return(teamsrv.Team{}, errors.New("database connection failed"))
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
		{
			name:           "invalid JSON format",
			requestBody:    "invalid json string",
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "empty team_name",
			requestBody: CreateTeamRequest{
				TeamName: "",
				Members: []MemberTeam{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "empty members array",
			requestBody: CreateTeamRequest{
				TeamName: "backend",
				Members:  []MemberTeam{},
			},
			setupMock: func(m *mockService) {
				m.On("CreateTeam", mock.Anything, teamsrv.Team{
					TeamName: "backend",
					Members:  []teamsrv.TeamMember{},
				}).Return(teamsrv.Team{
					TeamName: "backend",
					Members:  []teamsrv.TeamMember{},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response CreateTeamResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "backend", response.Team.TeamName)
				assert.Empty(t, response.Team.Members)
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
			router.POST("/team/add", handler.CreateTeam)

			var bodyBytes []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

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
