package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/YuanData/aquatrade/db/mock"
	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/token"
	"github.com/YuanData/aquatrade/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPaymentAPI(t *testing.T) {
	amount := int64(10)

	member1, _ := randomMember(t)
	member2, _ := randomMember(t)
	member3, _ := randomMember(t)

	trader1 := randomTrader(member1.Membername)
	trader2 := randomTrader(member2.Membername)
	trader3 := randomTrader(member3.Membername)

	trader1.Currency = util.AUD
	trader2.Currency = util.AUD
	trader3.Currency = util.JPY

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Generator)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(1).Return(trader2, nil)

				arg := db.PaymentTxParams{
					FromTraderID: trader1.ID,
					ToTraderID:   trader2.ID,
					Amount:       amount,
				}
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnauthorizedMember",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member2.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "FromTraderNotFound",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(db.Trader{}, sql.ErrNoRows)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToTraderNotFound",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(1).Return(db.Trader{}, sql.ErrNoRows)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromTraderCurrencyMismatch",
			body: gin.H{
				"from_trader_id": trader3.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member3.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader3.ID)).Times(1).Return(trader3, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToTraderCurrencyMismatch",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader3.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader3.ID)).Times(1).Return(trader3, nil)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       "XYZ",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         -amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetTraderError",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(1).Return(db.Trader{}, sql.ErrConnDone)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "PaymentTxError",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"amount":         amount,
				"currency":       util.AUD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(1).Return(trader2, nil)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(1).Return(db.PaymentTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/payments"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
