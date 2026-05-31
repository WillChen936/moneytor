package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "moneytor/database/mocks"
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateCategory(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	category := createRandomCategory(userID)

	testCases := []struct {
		name          string
		requestBody   gin.H
		setupAuth     func(t *testing.T, request *http.Request, server *Server)
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					UserID:            userID,
					Name:              category.Name,
					TransactionTypeID: category.TransactionTypeID,
				}
				mockStore.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(category, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidTransactionTypeID",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": -1,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TransactionTypeNotFound",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					UserID:            userID,
					Name:              category.Name,
					TransactionTypeID: category.TransactionTypeID,
				}
				mockStore.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Category{}, db.ErrForeignKeyViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "InternalError",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					UserID:            userID,
					Name:              category.Name,
					TransactionTypeID: category.TransactionTypeID,
				}
				mockStore.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Category{}, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestListCategory(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	categories := make([]db.Category, 10)
	for i := range categories {
		categories[i] = createRandomCategory(userID)
	}

	testCases := []struct {
		name          string
		queries       map[string]string
		setupAuth     func(t *testing.T, request *http.Request, server *Server)
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK_WithoutParams",
			queries: map[string]string{},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(1).Return(categories, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_WithParams",
			queries: map[string]string{
				"pageId":   "3",
				"pageSize": "10",
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListCategoriesParams{
					UserID: userID,
					Limit:  10,
					Offset: 20,
				}
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Eq(arg)).Times(1).Return(categories, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "Unauthorized",
			queries: map[string]string{},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			queries: map[string]string{
				"pageId":   "-1",
				"pageSize": "10",
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			queries: map[string]string{
				"pageId":   "3",
				"pageSize": "10",
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListCategoriesParams{
					UserID: userID,
					Limit:  10,
					Offset: 20,
				}
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Category{}, sql.ErrConnDone)
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

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/api/v1/categories", nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server)

			queries := request.URL.Query()
			for key, value := range testCase.queries {
				queries.Add(key, value)
			}
			request.URL.RawQuery = queries.Encode()

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func createRandomCategory(userID int64) db.Category {
	return db.Category{
		ID:                utils.RandomInt64Range(1, 1000),
		UserID:            userID,
		Name:              utils.RandomString(10),
		TransactionTypeID: utils.RandomInt16Range(1, 3),
	}
}
