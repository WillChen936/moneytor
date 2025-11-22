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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateCategory(t *testing.T) {
	category := createRandomCategory()

	testCases := []struct {
		Name          string
		requestBody   gin.H
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "OK",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
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
			Name: "IlleagalTranscationTypeID",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": -1,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name: "InvalidTranscationTypeID",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
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
			Name: "InternalError",
			requestBody: gin.H{
				"name":              category.Name,
				"transactionTypeId": category.TransactionTypeID,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockdb.NewMockStore(ctrl)
		testCase.buildStub(mockStore)

		server := NewServer(mockStore)

		data, err := json.Marshal(testCase.requestBody)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(data))
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		testCase.checkResponse(t, recorder)
	}
}

func TestListCategory(t *testing.T) {
	categories := []db.Category{}
	for i := 0; i < 10; i++ {
		categories = append(categories, createRandomCategory())
	}

	testCases := []struct {
		Name          string
		Queries       map[string]string
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name:    "OK_WithoutParams",
			Queries: map[string]string{},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(1).Return(categories, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name: "OK_WithParams",
			Queries: map[string]string{
				"page_id":   "3",
				"page_size": "10",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListCategoriesParams{
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
			Name: "InvalidPageID",
			Queries: map[string]string{
				"page_id":   "-1",
				"page_size": "10",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCategories(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name: "InternalError",
			Queries: map[string]string{
				"page_id":   "3",
				"page_size": "10",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListCategoriesParams{
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockdb.NewMockStore(ctrl)
		testCase.buildStub(mockStore)

		server := NewServer(mockStore)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/categories", nil)
		require.NoError(t, err)

		queries := request.URL.Query()
		for key, value := range testCase.Queries {
			queries.Add(key, value)
		}
		request.URL.RawQuery = queries.Encode()

		server.router.ServeHTTP(recorder, request)

		testCase.checkResponse(t, recorder)
	}
}

func createRandomCategory() db.Category {
	return db.Category{
		ID:                utils.RandomInt64Range(1, 1000),
		Name:              utils.RandomString(10),
		TransactionTypeID: utils.RandomInt16Range(1, 3),
	}
}
