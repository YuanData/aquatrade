package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/YuanData/aquatrade/db/mock"
	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/util"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetTraderAPI(t *testing.T) {
	trader := randomTrader()

	testCases := []struct {
		name          string
		traderID      int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			traderID: trader.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrader(t, recorder.Body, trader)
			},
		},
		{
			name:     "NotFound",
			traderID: trader.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(db.Trader{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			traderID: trader.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidID",
			traderID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/traders/%d", tc.traderID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTrader() db.Trader {
	return db.Trader{
		ID:       util.RandomInt(1, 1000),
		Holder:   util.RandomHolder(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func requireBodyMatchTrader(t *testing.T, body *bytes.Buffer, trader db.Trader) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTrader db.Trader
	err = json.Unmarshal(data, &gotTrader)
	require.NoError(t, err)
	require.Equal(t, trader, gotTrader)
}
