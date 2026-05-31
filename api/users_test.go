package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "moneytor/database/mocks"
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		requestBody   gin.H
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
				}
				mockStore.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			requestBody: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pgconn.PgError{Code: db.UniqueViolation})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			requestBody: gin.H{
				"username": user.Username,
				"email":    "not-an-email",
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordTooShort",
			requestBody: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": "short",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordTooLong",
			requestBody: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": utils.RandomString(73),
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UsernameTooShort",
			requestBody: gin.H{
				"username": "ab",
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UsernameTooLong",
			requestBody: gin.H{
				"username": utils.RandomString(51),
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "EmailTooLong",
			requestBody: gin.H{
				"username": user.Username,
				"email":    utils.RandomString(200) + "@test.com",
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			requestBody: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			testCase.buildStub(mockStore)

			server := newTestServer(t, mockStore)

			data, err := json.Marshal(testCase.requestBody)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestLogin(t *testing.T) {
	user, password := createRandomUser(t)

	var sessionID pgtype.UUID
	require.NoError(t, sessionID.Scan("12345678-1234-5678-1234-567812345678"))

	testCases := []struct {
		name          string
		requestBody   gin.H
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
				mockStore.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{ID: sessionID}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			requestBody: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
				mockStore.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "WrongPassword",
			requestBody: gin.H{
				"email":    user.Email,
				"password": "wrongpassword",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
				mockStore.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			requestBody: gin.H{
				"email":    "not-an-email",
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			requestBody: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			testCase.buildStub(mockStore)

			server := newTestServer(t, mockStore)

			data, err := json.Marshal(testCase.requestBody)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestRefresh(t *testing.T) {
	user, _ := createRandomUser(t)

	var sessionID pgtype.UUID
	require.NoError(t, sessionID.Scan("12345678-1234-5678-1234-567812345678"))

	testCases := []struct {
		name          string
		buildStub     func(mockStore *mockdb.MockStore, refreshToken string)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().
					GetSession(gomock.Any(), sessionID).
					Times(1).
					Return(db.Session{
						ID:           sessionID,
						UserID:       user.ID,
						RefreshToken: refreshToken,
						ExpiresAt:    time.Now().Add(time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidRefreshToken",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().GetSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionNotFound",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().
					GetSession(gomock.Any(), sessionID).
					Times(1).
					Return(db.Session{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "RefreshTokenMismatch",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().
					GetSession(gomock.Any(), sessionID).
					Times(1).
					Return(db.Session{
						ID:           sessionID,
						UserID:       user.ID,
						RefreshToken: "different-token",
						ExpiresAt:    time.Now().Add(time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UserIDMismatch",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().
					GetSession(gomock.Any(), sessionID).
					Times(1).
					Return(db.Session{
						ID:           sessionID,
						UserID:       user.ID + 1,
						RefreshToken: refreshToken,
						ExpiresAt:    time.Now().Add(time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionExpired",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().
					GetSession(gomock.Any(), sessionID).
					Times(1).
					Return(db.Session{
						ID:           sessionID,
						UserID:       user.ID,
						RefreshToken: refreshToken,
						ExpiresAt:    time.Now().Add(-time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			buildStub: func(mockStore *mockdb.MockStore, refreshToken string) {
				mockStore.EXPECT().
					GetSession(gomock.Any(), sessionID).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, mockStore)

			var refreshToken string
			if testCase.name == "InvalidRefreshToken" {
				refreshToken = "invalid.token.here"
			} else {
				var err error
				refreshToken, _, err = server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
				require.NoError(t, err)
			}

			testCase.buildStub(mockStore, refreshToken)

			body := gin.H{
				"sessionId":    sessionID,
				"refreshToken": refreshToken,
			}
			data, err := json.Marshal(body)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/api/v1/users/refresh", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := utils.RandomString(8)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	return db.User{
		ID:             utils.RandomInt64Range(1, 1000),
		Username:       utils.RandomString(6),
		Email:          utils.RandomString(6) + "@test.com",
		HashedPassword: string(hashedPassword),
	}, password
}

func eqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(arg.HashedPassword), []byte(e.password)); err != nil {
		return false
	}
	return e.arg.Username == arg.Username && e.arg.Email == arg.Email
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches CreateUserParams{Username=%s, Email=%s} with password=%s",
		e.arg.Username, e.arg.Email, e.password)
}
