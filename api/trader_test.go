package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/traders/%d", trader.ID)
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
		Account:  util.RandomAccount(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func requireBodyMatchTrader(t *testing.T, body *bytes.Buffer, trader db.Trader) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTrader db.Trader
	err = json.Unmarshal(data, &gotTrader)
	require.NoError(t, err)
	require.Equal(t, trader, gotTrader)
}
