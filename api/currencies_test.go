package api

import (
	"database/sql"
	mockdb "moneytor/database/mocks"
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListCurrencies(t *testing.T) {
	transactionTypes := []db.Currency{}
	for i := 0; i < 1; i++ {
		transactionTypes = append(transactionTypes, createRandomCurrency())
	}

	testCases := []struct {
		Name          string
		Queries       map[string]string
		buildStub     func(mockQuerier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "OK",
			buildStub: func(mockQuerier *mockdb.MockQuerier) {
				mockQuerier.EXPECT().ListCurrencies(gomock.Any()).Times(1).Return(transactionTypes, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name: "InternalError",
			buildStub: func(mockQuerier *mockdb.MockQuerier) {
				mockQuerier.EXPECT().ListCurrencies(gomock.Any()).Times(1).Return([]db.Currency{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := mockdb.NewMockQuerier(ctrl)
		testCase.buildStub(mockQuerier)

		server := NewServer(mockQuerier)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/currencies", nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		testCase.checkResponse(t, recorder)
	}
}

func createRandomCurrency() db.Currency {
	return db.Currency{
		ID:   utils.RandomInt16Range(1, 100),
		Code: utils.RandomString(3),
	}
}
