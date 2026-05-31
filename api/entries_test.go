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

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateEntry(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	account := createRandomAccount(userID)
	amount := int64(1000)

	categoryIncome := createRandomCategory(userID)
	categoryIncome.TransactionTypeID = TransactionTypeIncome
	categoryExpense := createRandomCategory(userID)
	categoryExpense.TransactionTypeID = TransactionTypeExpense

	entryIncome := createRandomEntry(account.ID, categoryIncome.ID)
	entryIncome.Amount = amount
	entryExpense := createRandomEntry(account.ID, categoryExpense.ID)
	entryExpense.Amount = amount

	testCases := []struct {
		name          string
		requestBody   gin.H
		setupAuth     func(t *testing.T, request *http.Request, server *Server)
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK_Income",
			requestBody: gin.H{
				"name":       entryIncome.Name,
				"note":       entryIncome.Note,
				"accountId":  entryIncome.FromAccountID,
				"categoryId": entryIncome.CategoryID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), db.GetCategoryParams{
					ID:     categoryIncome.ID,
					UserID: userID,
				}).Times(1).Return(categoryIncome, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryIncome.Name,
					Note:       entryIncome.Note,
					FromAccountID: entryIncome.FromAccountID,
					CategoryID: entryIncome.CategoryID,
					Amount:     amount,
				}
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.CreateEntryTxResult{Entry: entryIncome, FromAccount: account}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_Expense",
			requestBody: gin.H{
				"name":       entryExpense.Name,
				"note":       entryExpense.Note,
				"accountId":  entryExpense.FromAccountID,
				"categoryId": entryExpense.CategoryID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), db.GetCategoryParams{
					ID:     categoryExpense.ID,
					UserID: userID,
				}).Times(1).Return(categoryExpense, nil)
				arg := db.CreateEntryTxParams{
					Name:          entryExpense.Name,
					Note:          entryExpense.Note,
					FromAccountID: entryExpense.FromAccountID,
					CategoryID:    entryExpense.CategoryID,
					Amount:        -amount,
				}
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.CreateEntryTxResult{Entry: entryExpense, FromAccount: account}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			requestBody: gin.H{
				"name":       entryIncome.Name,
				"note":       entryIncome.Note,
				"accountId":  entryIncome.FromAccountID,
				"categoryId": entryIncome.CategoryID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidRequest",
			requestBody: gin.H{
				"name":   entryIncome.Name,
				"amount": amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetCategoryError",
			requestBody: gin.H{
				"name":       entryIncome.Name,
				"note":       entryIncome.Note,
				"accountId":  entryIncome.FromAccountID,
				"categoryId": categoryIncome.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), db.GetCategoryParams{
					ID:     categoryIncome.ID,
					UserID: userID,
				}).Times(1).Return(db.Category{}, sql.ErrNoRows)
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ForeignKeyViolation",
			requestBody: gin.H{
				"name":       entryIncome.Name,
				"note":       entryIncome.Note,
				"accountId":  entryIncome.FromAccountID,
				"categoryId": entryIncome.CategoryID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), db.GetCategoryParams{
					ID:     categoryIncome.ID,
					UserID: userID,
				}).Times(1).Return(categoryIncome, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryIncome.Name,
					Note:       entryIncome.Note,
					FromAccountID: entryIncome.FromAccountID,
					CategoryID: entryIncome.CategoryID,
					Amount:     amount,
				}
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.CreateEntryTxResult{}, db.ErrForeignKeyViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "InvalidTransactionTypeID",
			requestBody: gin.H{
				"name":       utils.RandomString(10),
				"accountId":  account.ID,
				"categoryId": utils.RandomInt64Range(1, 1000),
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				transferCategory := createRandomCategory(userID)
				transferCategory.TransactionTypeID = TransactionTypeTransfer
				mockStore.EXPECT().GetCategory(gomock.Any(), gomock.Any()).
					Times(1).Return(transferCategory, nil)
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			requestBody: gin.H{
				"name":       entryIncome.Name,
				"note":       entryIncome.Note,
				"accountId":  entryIncome.FromAccountID,
				"categoryId": entryIncome.CategoryID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), db.GetCategoryParams{
					ID:     categoryIncome.ID,
					UserID: userID,
				}).Times(1).Return(categoryIncome, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryIncome.Name,
					Note:       entryIncome.Note,
					FromAccountID: entryIncome.FromAccountID,
					CategoryID: entryIncome.CategoryID,
					Amount:     amount,
				}
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.CreateEntryTxResult{}, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodPost, "/api/v1/entries", bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestListEntries(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	accountID := utils.RandomInt64Range(1, 1000)
	entries := make([]db.Entry, 5)
	for i := range entries {
		entries[i] = createRandomEntry(accountID, int64(i+1))
	}

	testCases := []struct {
		name          string
		queries       map[string]string
		setupAuth     func(t *testing.T, request *http.Request, server *Server)
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			queries: map[string]string{
				"accountId": fmt.Sprintf("%d", accountID),
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesByAccountIDParams{
					AccountID: accountID,
					Limit:     5,
					Offset:    0,
				}
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_WithPagination",
			queries: map[string]string{
				"accountId": fmt.Sprintf("%d", accountID),
				"pageId":    "2",
				"pageSize":  "10",
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesByAccountIDParams{
					AccountID: accountID,
					Limit:     10,
					Offset:    10,
				}
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "Unauthorized",
			queries: map[string]string{"accountId": fmt.Sprintf("%d", accountID)},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:    "MissingAccountID",
			queries: map[string]string{},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			queries: map[string]string{
				"accountId": fmt.Sprintf("%d", accountID),
				"pageId":    "0",
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			queries: map[string]string{
				"accountId": fmt.Sprintf("%d", accountID),
			},
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, request, server, authorizationTypeBearer, userID, time.Minute)
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesByAccountIDParams{
					AccountID: accountID,
					Limit:     5,
					Offset:    0,
				}
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodGet, "/api/v1/entries", nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server)

			q := request.URL.Query()
			for k, v := range testCase.queries {
				q.Add(k, v)
			}
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func createRandomEntry(accountID, categoryID int64) db.Entry {
	return db.Entry{
		Name:          utils.RandomString(10),
		Note:          utils.RandomString(10),
		FromAccountID: accountID,
		ToAccountID:   pgtype.Int8{Valid: false},
		CategoryID:    categoryID,
		Amount:        utils.RandomInt64Range(100, 10000),
	}
}
