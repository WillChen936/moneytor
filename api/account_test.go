package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "moneytor/database/mocks"
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mockdb.NewMockQuerier(ctrl)
	server := NewServer(mockQuerier)

	testName := utils.RandomString(10)
	testCurrencyID := int16(1)
	testBalance, err := decimal.NewFromString("1000.00")
	require.NoError(t, err)

	requestBody := createAccountRequest{
		Name:       testName,
		CurrencyID: testCurrencyID,
		InitialBalance: Decimal{
			Decimal: testBalance,
		},
	}

	params := db.CreateAccountParams{
		Name:       testName,
		CurrencyID: testCurrencyID,
		Balance:    testBalance,
	}

	stubs := db.Account{
		ID:         utils.RandomInt64Range(1, 1000),
		Name:       utils.RandomString(10),
		CurrencyID: testCurrencyID,
		Balance:    testBalance,
		CreatedAt:  time.Now(),
		UpdatedAt: pgtype.Timestamptz{
			Valid: false,
		},
	}

	mockQuerier.EXPECT().CreateAccount(gomock.Any(), eqCreateAccountParams(params)).Times(1).Return(stubs, nil)

	data, err := json.Marshal(requestBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts", bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestCreateAccountFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mockdb.NewMockQuerier(ctrl)
	server := NewServer(mockQuerier)

	testName := utils.RandomString(10)
	testCurrencyID := int16(-1)
	testBalance, err := decimal.NewFromString("1000.00")
	require.NoError(t, err)

	requestBody := createAccountRequest{
		Name:       testName,
		CurrencyID: testCurrencyID,
		InitialBalance: Decimal{
			Decimal: testBalance,
		},
	}

	mockQuerier.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)

	data, err := json.Marshal(requestBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/api/v1/accounts", bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
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

	return e.arg.Name == arg.Name &&
		e.arg.CurrencyID == arg.CurrencyID &&
		e.arg.Balance.Equal(arg.Balance)
}

func (e eqCreateAccountParamsMatcher) String() string {
	return fmt.Sprintf("is equal to CreateAccountParams{Name=%s, CurrencyID=%d, Balance=%s}",
		e.arg.Name, e.arg.CurrencyID, e.arg.Balance.String())
}
