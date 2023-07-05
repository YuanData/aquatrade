package api

import (
	"fmt"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
		Times(1).
		Return(trader, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/traders/%d", trader.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomTrader() db.Trader {
	return db.Trader{
		ID:       util.RandomInt(1, 1000),
		Account:  util.RandomAccount(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
