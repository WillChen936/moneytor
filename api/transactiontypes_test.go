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
	transactionTypes := make([]db.TransactionType, 3)
	for i := range transactionTypes {
		transactionTypes[i] = createRandomTransactionType()
	}

	testCases := []struct {
		name          string
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListTransactionTypes(gomock.Any()).Times(1).Return(transactionTypes, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListTransactionTypes(gomock.Any()).Times(1).Return([]db.TransactionType{}, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodGet, "/api/v1/transaction-types", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func createRandomTransactionType() db.TransactionType {
	return db.TransactionType{
		ID:   utils.RandomInt16Range(1, 1000),
		Name: utils.RandomString(6),
	}
}
