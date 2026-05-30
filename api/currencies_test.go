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
	userID := utils.RandomInt64Range(1, 1000)
	currencies := []db.Currency{createRandomCurrency()}

	testCases := []struct {
		Name          string
		buildStub     func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "OK",
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCurrencies(gomock.Any()).Times(1).Return(currencies, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name: "InternalError",
			buildStub: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().ListCurrencies(gomock.Any()).Times(1).Return([]db.Currency{}, sql.ErrConnDone)
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

		server, maker := newTestServer(t, mockStore)

		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/currencies", nil)
		require.NoError(t, err)

		addAuthorization(t, request, maker, userID)

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
