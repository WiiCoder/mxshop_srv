package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
	"mxshop_srvs/user_srv/utils"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct{}

func ModelToResponse(user model.User) *proto.UserInfoResponse {
	response := &proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     user.Role,
	}

	if user.Birthday != nil {
		response.BirthDay = uint64(user.Birthday.Unix())
	}

	return response
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(_ context.Context, pageInfo *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(pageInfo.PageNum), int(pageInfo.PageSize))).Find(&users)

	for _, user := range users {
		userInfoResponse := ModelToResponse(user)
		rsp.Data = append(rsp.Data, userInfoResponse)
	}
	return rsp, nil
}

func (s *UserServer) GetUserMobile(_ context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	response := ModelToResponse(user)
	return response, nil
}

func (s *UserServer) GetUserId(_ context.Context, request *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.First(&user, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	response := ModelToResponse(user)
	return response, nil
}

func (s *UserServer) CreteUser(_ context.Context, request *proto.CreteUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = request.Mobile
	user.NickName = request.NickName

	options := &utils.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodePwd := utils.Encode(request.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodePwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResponse := ModelToResponse(user)
	return userInfoResponse, nil
}

func (s *UserServer) UpdateUser(_ context.Context, request *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户不存在")
	}

	unix := time.Unix(int64(request.BirthDay), 0)
	user.Birthday = &unix
	user.Gender = request.Gender
	user.NickName = request.NickName

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (s *UserServer) CheckPassword(_ context.Context, request *proto.CheckPasswordInfo) (*proto.CheckResponse, error) {
	options := &utils.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	split := strings.Split(request.EncryptedPassword, "$")
	verify := utils.Verify(request.Password, split[2], split[3], options)
	return &proto.CheckResponse{Success: verify}, nil
}
