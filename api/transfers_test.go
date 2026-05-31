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

func TestCreateTransfer(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	fromAccount := createRandomAccount(userID)
	toAccount := createRandomAccount(userID)
	amount := int64(1000)

	categoryTransfer := createRandomCategory(userID)
	categoryTransfer.TransactionTypeID = TransactionTypeTransfer

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
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     fromAccount.ID,
					UserID: userID,
				}).Times(1).Return(fromAccount, nil)
				mockStore.EXPECT().GetCategory(gomock.Any(), db.GetCategoryParams{
					ID:     categoryTransfer.ID,
					UserID: userID,
				}).Times(1).Return(categoryTransfer, nil)
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     toAccount.ID,
					UserID: userID,
				}).Times(1).Return(toAccount, nil)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.CreateTransferTxResult{
					FromAccount: fromAccount,
					ToAccount:   toAccount,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "AmountTooLarge",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        int64(1000000000000),
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "SameAccount",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   fromAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     fromAccount.ID,
					UserID: userID,
				}).Times(1).Return(db.Account{}, db.ErrRecordNotFound)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "GetCategoryNotFound",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     fromAccount.ID,
					UserID: userID,
				}).Times(1).Return(fromAccount, nil)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Category{}, db.ErrRecordNotFound)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NonTransferCategory",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				expenseCategory := createRandomCategory(userID)
				expenseCategory.TransactionTypeID = TransactionTypeExpense
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     fromAccount.ID,
					UserID: userID,
				}).Times(1).Return(fromAccount, nil)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).
					Times(1).Return(expenseCategory, nil)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     fromAccount.ID,
					UserID: userID,
				}).Times(1).Return(fromAccount, nil)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).
					Times(1).Return(categoryTransfer, nil)
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     toAccount.ID,
					UserID: userID,
				}).Times(1).Return(db.Account{}, db.ErrRecordNotFound)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			requestBody: gin.H{
				"name":          utils.RandomString(10),
				"fromAccountId": fromAccount.ID,
				"toAccountId":   toAccount.ID,
				"categoryId":    categoryTransfer.ID,
				"amount":        amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     fromAccount.ID,
					UserID: userID,
				}).Times(1).Return(fromAccount, nil)
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).
					Times(1).Return(categoryTransfer, nil)
				mockStore.EXPECT().GetAccount(gomock.Any(), db.GetAccountParams{
					ID:     toAccount.ID,
					UserID: userID,
				}).Times(1).Return(toAccount, nil)
				mockStore.EXPECT().CreateTransferTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.CreateTransferTxResult{}, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodPost, "/api/v1/transfers", bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}
