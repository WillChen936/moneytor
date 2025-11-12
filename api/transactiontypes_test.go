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

func TestListTransactionType(t *testing.T) {
	transactionTypes := []db.TransactionType{}
	for i := 0; i < 1; i++ {
		transactionTypes = append(transactionTypes, createRandomTranscationType())
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
				mockQuerier.EXPECT().ListTransactionTypes(gomock.Any()).Times(1).Return(transactionTypes, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name: "InternalError",
			buildStub: func(mockQuerier *mockdb.MockQuerier) {
				mockQuerier.EXPECT().ListTransactionTypes(gomock.Any()).Times(1).Return([]db.TransactionType{}, sql.ErrConnDone)
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
		request, err := http.NewRequest(http.MethodGet, "/api/v1/transaction-types", nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		testCase.checkResponse(t, recorder)
	}
}

func createRandomTranscationType() db.TransactionType {
	return db.TransactionType{
		ID:   utils.RandomInt16Range(1, 1000),
		Name: utils.RandomString(6),
	}
}
