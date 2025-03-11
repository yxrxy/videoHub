package controller

import (
	"context"
	"io"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/yxrrxy/videoHub/kitex_gen/user"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/response"
	"github.com/yxrrxy/videoHub/pkg/utils"
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
	if err := req.BindAndValidate(&registerReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, errno.ErrInvalidParam.ErrMsg))
		return
	}

	resp, err := c.client.Register(ctx, &registerReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	userInfo, err := c.client.GetUserInfo(ctx, &user.UserInfoRequest{UserId: resp.UserId})
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	data := map[string]interface{}{
		"id":            strconv.FormatInt(userInfo.User.Id, 10),
		"username":      userInfo.User.Username,
		"avatar_url":    userInfo.User.AvatarUrl,
		"created_at":    utils.FormatTimestamp(userInfo.User.CreatedAt),
		"updated_at":    utils.FormatTimestamp(userInfo.User.UpdatedAt),
		"deleted_at":    utils.FormatTimestamp(userInfo.User.DeletedAt),
		"token":         resp.Token,
		"refresh_token": resp.RefreshToken,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}

// Login 处理用户登录
func (c *UserController) Login(ctx context.Context, req *app.RequestContext) {
	var loginReq user.LoginRequest
	if err := req.BindAndValidate(&loginReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, errno.ErrInvalidParam.ErrMsg))
		return
	}

	resp, err := c.client.Login(ctx, &loginReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	userInfo, err := c.client.GetUserInfo(ctx, &user.UserInfoRequest{UserId: resp.UserId})
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	data := map[string]interface{}{
		"id":            strconv.FormatInt(userInfo.User.Id, 10),
		"username":      userInfo.User.Username,
		"avatar_url":    userInfo.User.AvatarUrl,
		"created_at":    utils.FormatTimestamp(userInfo.User.CreatedAt),
		"updated_at":    utils.FormatTimestamp(userInfo.User.UpdatedAt),
		"deleted_at":    utils.FormatTimestamp(userInfo.User.DeletedAt),
		"token":         resp.Token,
		"refresh_token": resp.RefreshToken,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}

// UploadAvatar 处理头像上传
func (c *UserController) UploadAvatar(ctx context.Context, req *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		req.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	file, err := req.FormFile("avatar")
	if err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, errno.ErrInvalidParam.ErrMsg))
		return
	}

	f, err := file.Open()
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	rpcCtx := pkgcontext.WithUserID(ctx, userID)

	uploadReq := &user.UploadAvatarRequest{
		UserId:      userID,
		AvatarData:  fileData,
		ContentType: file.Header.Get("Content-Type"),
	}

	resp, err := c.client.UploadAvatar(rpcCtx, uploadReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	data := map[string]interface{}{
		"avatar_url": resp.AvatarUrl,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}

func (c *UserController) GetUserInfo(ctx context.Context, req *app.RequestContext) {
	userID := req.Query("user_id")
	if userID == "" {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, errno.ErrInvalidParam.ErrMsg))
		return
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, errno.ErrInvalidParam.ErrMsg))
		return
	}

	userInfo, err := c.client.GetUserInfo(ctx, &user.UserInfoRequest{UserId: userIDInt})
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	data := map[string]interface{}{
		"id":         strconv.FormatInt(userInfo.User.Id, 10),
		"username":   userInfo.User.Username,
		"avatar_url": userInfo.User.AvatarUrl,
		"created_at": utils.FormatTimestamp(userInfo.User.CreatedAt),
		"updated_at": utils.FormatTimestamp(userInfo.User.UpdatedAt),
		"deleted_at": utils.FormatTimestamp(userInfo.User.DeletedAt),
	}

	req.JSON(consts.StatusOK, response.Success(data))
}
