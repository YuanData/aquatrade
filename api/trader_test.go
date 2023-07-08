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
	"time"

	mockdb "github.com/YuanData/aquatrade/db/mock"
	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/token"
	"github.com/YuanData/aquatrade/util"
	"github.com/gin-gonic/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetTraderAPI(t *testing.T) {
	user, _ := randomMember(t)
	trader := randomTrader(user.Membername)

	testCases := []struct {
		name          string
		traderID      int64
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Generator)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			traderID: trader.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
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
			name:     "UnauthorizedMember",
			traderID: trader.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "NoAuthorization",
			traderID: trader.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "NotFound",
			traderID: trader.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},

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
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
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
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
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

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateTraderAPI(t *testing.T) {
	user, _ := randomMember(t)
	trader := randomTrader(user.Membername)

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
				"currency": trader.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTraderParams{
					Holder:   trader.Holder,
					Currency: trader.Currency,
					Balance:  0,
				}

				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrader(t, recorder.Body, trader)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"currency": trader.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"currency": trader.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"currency": "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/traders"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListTradersAPI(t *testing.T) {
	user, _ := randomMember(t)

	n := 5
	traders := make([]db.Trader, n)
	for i := 0; i < n; i++ {
		traders[i] = randomTrader(user.Membername)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Generator)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTradersParams{
					Holder: user.Membername,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(traders, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTraders(t, recorder.Body, traders)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			url := "/traders"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomTrader(owner string) db.Trader {
	return db.Trader{
		ID:       util.RandomInt(1, 1000),
		Holder:   owner,
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

func requireBodyMatchTraders(t *testing.T, body *bytes.Buffer, traders []db.Trader) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTraders []db.Trader
	err = json.Unmarshal(data, &gotTraders)
	require.NoError(t, err)
	require.Equal(t, traders, gotTraders)
}
