package handler

import (
	"context"
	"gorm.io/gorm"
	"goshop_srvs/user_srv/global"
	"goshop_srvs/user_srv/model"
	"goshop_srvs/user_srv/proto"
)

// 只需要实现user.proto里的方法即可
// 再生成的pb.go文件中找 interface 关键字，找到那几个API

type UserServer struct{}

// 指向 UserServer 结构体的指针;第二个参数PageInfo 消息类型的指针;返回值是一个指向 UserListResponse 消息类型的指针和一个错误
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// 创建响应消息
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	// gorm分页
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
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

func ModelToResponse(user model.User) proto.UserInfoResponse {
	// 在grpc的message字段中，如果有默认值不能随便赋值nil，容易出错
	// 要搞清哪些有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}
