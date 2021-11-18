package handler

import (
	"context"
	"mxshop_srvs/user_srv/proto"
	"testing"
)

func TestUserServer_CreteUser(t *testing.T) {
	ctx := context.Background()
	s := &UserServer{}
	createList := []proto.CreteUserInfo{
		{
			NickName: "user1",
			PassWord: "user1",
			Mobile:   "11111111111",
		},
		{
			NickName: "user2",
			PassWord: "user2",
			Mobile:   "12222222222",
		},
		{
			NickName: "user3",
			PassWord: "user3",
			Mobile:   "13333333333",
		},
	}

	for i := range createList {
		_, _ = s.CreteUser(ctx, &createList[i])
	}

}

func TestUserServer_GetUserList(t *testing.T) {
	ctx := context.Background()
	s := &UserServer{}

	pageArr := []proto.PageInfo{
		{
			PageNum:  1,
			PageSize: 1,
		},
		{
			PageNum:  1,
			PageSize: 2,
		},
		{
			PageNum:  1,
			PageSize: 10,
		},
		{
			PageNum:  2,
			PageSize: 10,
		},
	}

	for i := range pageArr {
		_, err := s.GetUserList(ctx, &pageArr[i])
		if err != nil {
			panic(err)
		}
	}
}

func TestUserServer_GetUserId(t *testing.T) {
	ctx := context.Background()
	s := &UserServer{}
	_, err := s.GetUserId(ctx, &proto.IdRequest{Id: 1})
	if err != nil {
		panic(err)
	}
}

func TestUserServer_GetUserMobile(t *testing.T) {
	ctx := context.Background()
	s := &UserServer{}
	_, err := s.GetUserMobile(ctx, &proto.MobileRequest{Mobile: "11111111111"})
	if err != nil {
		panic(err)
	}
}
