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
	currencies := make([]db.Currency, 5)
	for i := range currencies {
		currencies[i] = createRandomCurrency()
	}

	testCases := []struct {
		name          string
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCurrencies(gomock.Any()).Times(1).Return(currencies, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCurrencies(gomock.Any()).Times(1).Return([]db.Currency{}, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodGet, "/api/v1/currencies", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func createRandomCurrency() db.Currency {
	return db.Currency{
		ID:   utils.RandomInt16Range(1, 100),
		Code: utils.RandomString(3),
	}
}
