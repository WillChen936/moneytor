package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
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

func TestCreateAccount(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	account := createRandomAccountForUser(userID)

	testCases := []struct {
		name          string
		requestBody   gin.H
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"name":       account.Name,
				"currencyId": account.CurrencyID,
				"balance":    account.Balance,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					UserID:     userID,
					Name:       account.Name,
					CurrencyID: account.CurrencyID,
					Balance:    account.Balance,
				}
				mockStore.EXPECT().CreateAccount(gomock.Any(), eqCreateAccountParams(arg)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "IllegalCurrencyID",
			requestBody: gin.H{
				"name":       account.Name,
				"currencyId": -1,
				"balance":    account.Balance,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ForeignKeyViolation",
			requestBody: gin.H{
				"name":       account.Name,
				"currencyId": math.MaxInt16,
				"balance":    account.Balance,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					UserID:     userID,
					Name:       account.Name,
					CurrencyID: math.MaxInt16,
					Balance:    account.Balance,
				}
				mockStore.EXPECT().CreateAccount(gomock.Any(), eqCreateAccountParams(arg)).Times(1).Return(db.Account{}, db.ErrForeignKeyViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "InternalError",
			requestBody: gin.H{
				"name":       account.Name,
				"currencyId": account.CurrencyID,
				"balance":    account.Balance,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			requestBody: gin.H{
				"name":       account.Name,
				"currencyId": account.CurrencyID,
				"balance":    account.Balance,
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mockdb.NewMockStore(ctrl)
		testCase.buildStub(mockStore)

		server := newTestServer(t, mockStore)

		data, err := json.Marshal(testCase.requestBody)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts", bytes.NewReader(data))
		require.NoError(t, err)

		if testCase.name != "Unauthorized" {
			addAuthorization(t, request, server, userID)
		}

		server.router.ServeHTTP(recorder, request)
		testCase.checkResponse(t, recorder)
	}
}

func TestListAccounts(t *testing.T) {
	userID := utils.RandomInt64Range(1, 1000)
	accounts := make([]db.Account, 10)
	for i := range accounts {
		accounts[i] = createRandomAccountForUser(userID)
	}

	testCases := []struct {
		name          string
		queries       map[string]string
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK_WithoutParams",
			queries: map[string]string{},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return(accounts, nil)
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
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					UserID: userID,
					Limit:  10,
					Offset: 20,
				}
				mockStore.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			queries: map[string]string{
				"pageId":   "-1",
				"pageSize": "10",
			},
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
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
			buildStub: func(mockStore *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					UserID: userID,
					Limit:  10,
					Offset: 20,
				}
				mockStore.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Account{}, sql.ErrConnDone)
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

		server := newTestServer(t, mockStore)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/accounts", nil)
		require.NoError(t, err)

		addAuthorization(t, request, server, userID)

		queries := request.URL.Query()
		for key, value := range testCase.queries {
			queries.Add(key, value)
		}
		request.URL.RawQuery = queries.Encode()

		server.router.ServeHTTP(recorder, request)
		testCase.checkResponse(t, recorder)
	}
}

func createRandomAccountForUser(userID int64) db.Account {
	return db.Account{
		ID:         utils.RandomInt64Range(1, 1000),
		UserID:     userID,
		Name:       utils.RandomString(10),
		CurrencyID: utils.RandomInt16Range(1, 10),
		Balance:    utils.RandomInt64Range(1, 100),
	}
}

func eqCreateAccountParams(arg db.CreateAccountParams) gomock.Matcher {
	return eqCreateAccountParamsMatcher{arg}
}

type eqCreateAccountParamsMatcher struct {
	arg db.CreateAccountParams
}

func (e eqCreateAccountParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateAccountParams)
	if !ok {
		return false
	}
	return e.arg.UserID == arg.UserID &&
		e.arg.Name == arg.Name &&
		e.arg.CurrencyID == arg.CurrencyID &&
		e.arg.Balance == arg.Balance
}

func (e eqCreateAccountParamsMatcher) String() string {
	return fmt.Sprintf("is equal to CreateAccountParams{UserID=%d, Name=%s, CurrencyID=%d, Balance=%d}",
		e.arg.UserID, e.arg.Name, e.arg.CurrencyID, e.arg.Balance)
}
