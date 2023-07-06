package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/YuanData/aquatrade/db/mock"
	db "github.com/YuanData/aquatrade/db/sqlc"
	"github.com/YuanData/aquatrade/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateMemberParamsMatcher struct {
	arg      db.CreateMemberParams
	password string
}

func (e eqCreateMemberParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateMemberParams)
	if !ok {
		return false
	}

	err := util.VerifyPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateMemberParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateMemberParams(arg db.CreateMemberParams, password string) gomock.Matcher {
	return eqCreateMemberParamsMatcher{arg, password}
}

func TestCreateMemberAPI(t *testing.T) {
	member, password := randomMember(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"membername": member.Membername,
				"password":   password,
				"full_name":  member.FullName,
				"email":      member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateMemberParams{
					Membername: member.Membername,
					FullName:   member.FullName,
					Email:      member.Email,
				}
				store.EXPECT().
					CreateMember(gomock.Any(), EqCreateMemberParams(arg, password)).
					Times(1).
					Return(member, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMember(t, recorder.Body, member)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"membername": member.Membername,
				"password":   password,
				"full_name":  member.FullName,
				"email":      member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Member{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateMembername",
			body: gin.H{
				"membername": member.Membername,
				"password":   password,
				"full_name":  member.FullName,
				"email":      member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Member{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidMembername",
			body: gin.H{
				"membername": "invalid-member#1",
				"password":   password,
				"full_name":  member.FullName,
				"email":      member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"membername": member.Membername,
				"password":   password,
				"full_name":  member.FullName,
				"email":      "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"membername": member.Membername,
				"password":   "123",
				"full_name":  member.FullName,
				"email":      member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/members"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func randomMember(t *testing.T) (member db.Member, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	member = db.Member{
		Membername:     util.RandomHolder(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomHolder(),
		Email:          util.RandomEmail(),
	}
	return
}
func requireBodyMatchMember(t *testing.T, body *bytes.Buffer, member db.Member) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotMember db.Member
	err = json.Unmarshal(data, &gotMember)

	require.NoError(t, err)
	require.Equal(t, member.Membername, gotMember.Membername)
	require.Equal(t, member.FullName, gotMember.FullName)
	require.Equal(t, member.Email, gotMember.Email)
	require.Empty(t, gotMember.HashedPassword)
}
