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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount()
	amount := int64(1000)

	categoryIncome := db.Category{
		ID:                utils.RandomInt64Range(1, 1000),
		Name:              utils.RandomString(10),
		TransactionTypeID: TransactionTypeIncome,
	}
	categoryExpense := db.Category{
		ID:                utils.RandomInt64Range(1001, 2000),
		Name:              utils.RandomString(10),
		TransactionTypeID: TransactionTypeExpense,
	}
	categoryTransfer := db.Category{
		ID:                utils.RandomInt64Range(2001, 3000),
		Name:              utils.RandomString(10),
		TransactionTypeID: TransactionTypeTransfer,
	}

	entryIncome := createRandomEntry(account.ID, categoryIncome.ID)
	entryIncome.Amount = amount
	entryExpense := createRandomEntry(account.ID, categoryExpense.ID)
	entryExpense.Amount = amount
	entryTransfer := createRandomEntry(account.ID, categoryTransfer.ID)
	entryTransfer.Amount = amount

	testCases := []struct {
		name          string
		requestBody   gin.H
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
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), categoryIncome.ID).Times(1).Return(categoryIncome, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryIncome.Name,
					Note:       entryIncome.Note,
					AccountID:  entryIncome.FromAccountID,
					CategoryID: entryIncome.CategoryID,
					Amount:     amount,
				}
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.CreateEntryTxResult{Entry: entryIncome, Account: account}, nil)
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
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), categoryExpense.ID).Times(1).Return(categoryExpense, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryExpense.Name,
					Note:       entryExpense.Note,
					AccountID:  entryExpense.FromAccountID,
					CategoryID: entryExpense.CategoryID,
					Amount:     -amount,
				}
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.CreateEntryTxResult{Entry: entryExpense, Account: account}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidRequest",
			requestBody: gin.H{
				"name":   entryIncome.Name,
				"amount": amount,
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
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), categoryIncome.ID).Times(1).Return(db.Category{}, sql.ErrNoRows)
				mockStore.EXPECT().CreateEntryTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), categoryIncome.ID).Times(1).Return(categoryIncome, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryIncome.Name,
					Note:       entryIncome.Note,
					AccountID:  entryIncome.FromAccountID,
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
				"name":       entryTransfer.Name,
				"note":       entryTransfer.Note,
				"accountId":  entryTransfer.FromAccountID,
				"categoryId": entryTransfer.CategoryID,
				"amount":     amount,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), categoryTransfer.ID).Times(1).Return(categoryTransfer, nil)
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
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetCategory(gomock.Any(), categoryIncome.ID).Times(1).Return(categoryIncome, nil)
				arg := db.CreateEntryTxParams{
					Name:       entryIncome.Name,
					Note:       entryIncome.Note,
					AccountID:  entryIncome.FromAccountID,
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockdb.NewMockStore(ctrl)
		testCase.buildStub(mockStore)

		server := NewServer(mockStore)

		data, err := json.Marshal(testCase.requestBody)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/api/v1/entries", bytes.NewReader(data))
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		testCase.checkResponse(t, recorder)
	}
}

func TestGetEntries(t *testing.T) {
	entries := []db.Entry{}
	for i := 0; i < 5; i++ {
		entries = append(entries, createRandomEntry(int64(i+1), int64(i+1)))
	}

	testCases := []struct {
		name          string
		queries       map[string]string
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK_WithoutAccountID",
			queries: map[string]string{},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesParams{
					Limit:  5,
					Offset: 0,
				}
				mockStore.EXPECT().ListEntries(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_WithAccountID",
			queries: map[string]string{
				"account_id": "123",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesByAccountIDParams{
					AccountID: 123,
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
				"page_id":   "2",
				"page_size": "10",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesParams{
					Limit:  10,
					Offset: 10,
				}
				mockStore.EXPECT().ListEntries(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			queries: map[string]string{
				"page_id":   "0",
				"page_size": "5",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListEntries(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			queries: map[string]string{
				"page_id":   "1",
				"page_size": "4",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListEntries(gomock.Any(), gomock.Any()).Times(0)
				mockStore.EXPECT().ListEntriesByAccountID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			queries: map[string]string{},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListEntriesParams{
					Limit:  5,
					Offset: 0,
				}
				mockStore.EXPECT().ListEntries(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil, sql.ErrConnDone)
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
		request, err := http.NewRequest(http.MethodGet, "/api/v1/entries", nil)
		require.NoError(t, err)

		q := request.URL.Query()
		for k, v := range testCase.queries {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()

		server.router.ServeHTTP(recorder, request)

		testCase.checkResponse(t, recorder)
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
