package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/yxrrxy/videoHub/kitex_gen/user"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrrxy/videoHub/pkg/errno"
)

type UserController struct {
	client userservice.Client
}

func NewUserController(client userservice.Client) *UserController {
	return &UserController{client: client}
}

// Register 处理用户注册
func (c *UserController) Register(ctx context.Context, req *app.RequestContext) {
	var registerReq user.RegisterRequest
	if err := req.Bind(&registerReq); err != nil {
		req.JSON(consts.StatusBadRequest, errno.ErrBind)
		return
	}

	resp, err := c.client.Register(ctx, &registerReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, err)
		return
	}

	req.JSON(consts.StatusOK, resp)
}

// Login 处理用户登录
func (c *UserController) Login(ctx context.Context, req *app.RequestContext) {
	var loginReq user.LoginRequest
	if err := req.Bind(&loginReq); err != nil {
		req.JSON(consts.StatusBadRequest, errno.ErrBind)
		return
	}

	resp, err := c.client.Login(ctx, &loginReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, err)
		return
	}

	req.JSON(consts.StatusOK, resp)
}

// UploadAvatar 处理头像上传
func (c *UserController) UploadAvatar(ctx context.Context, req *app.RequestContext) {
	file, err := req.FormFile("avatar")
	if err != nil {
		req.JSON(consts.StatusBadRequest, errno.ErrBind)
		return
	}

	fileData, err := file.Data()
	if err != nil {
		req.JSON(consts.StatusInternalServerError, err)
		return
	}

	uploadReq := &user.UploadAvatarRequest{
		AvatarData:  fileData,
		ContentType: file.Header.Get("Content-Type"),
	}

	// 从请求头获取token并设置到context
	token := req.GetHeader("Authorization")
	newCtx := context.WithValue(ctx, "token", token)

	resp, err := c.client.UploadAvatar(newCtx, uploadReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, err)
		return
	}

	req.JSON(consts.StatusOK, resp)
}
