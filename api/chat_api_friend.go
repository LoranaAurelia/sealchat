package api

import (
	"sealchat/model"
	"sealchat/utils"
)

func apiFriendChannelList(ctx *ChatContext, data *struct {
	Next string `json:"next"`
}) (any, error) {
	items, _ := model.FriendChannelList(ctx.User.ID)

	return &struct {
		Data []*model.ChannelModel `json:"data"`
		Next string                `json:"next"`
	}{
		Data: items,
	}, nil
}

// 好友申请列表
func apiFriendRequestList(ctx *ChatContext, data *struct {
	Next string `json:"next"`
}) (any, error) {
	items, _ := model.FriendRequestListByReceiverID(ctx.User.ID)

	utils.QueryOneToManyMap(model.GetDB(), items, func(i *model.FriendRequestModel) []string {
		return []string{i.SenderID, i.ReceiverID}
	}, func(i *model.FriendRequestModel, x []*model.UserModel) {
		i.UserInfoSender = x[0]
		i.UserInfoReceiver = x[1]
	}, "")

	return &struct {
		Data []*model.FriendRequestModel `json:"data"`
		Next string                      `json:"next"`
	}{
		Data: items,
	}, nil
}

// 好友申请列表
func apiFriendRequestSenderList(ctx *ChatContext, data *struct {
	Next string `json:"next"`
}) (any, error) {
	items, _ := model.FriendRequestListBySenderID(ctx.User.ID)

	utils.QueryOneToManyMap(model.GetDB(), items, func(i *model.FriendRequestModel) []string {
		return []string{i.SenderID, i.ReceiverID}
	}, func(i *model.FriendRequestModel, x []*model.UserModel) {
		i.UserInfoSender = x[0]
		i.UserInfoReceiver = x[1]
	}, "")

	return &struct {
		Data []*model.FriendRequestModel `json:"data"`
		Next string                      `json:"next"`
	}{
		Data: items,
	}, nil
}

func apiFriendRequestCreate(ctx *ChatContext, data *struct {
	SenderID   string `json:"senderId"`   // 发送者
	ReceiverID string `json:"receiverId"` // 接收者
	Note       string `json:"note"`       // 申请理由
}) (any, error) {
	err := model.FriendRequestCreate(&model.FriendRequestModel{
		SenderID:   data.SenderID,
		ReceiverID: data.ReceiverID,
		Note:       data.Note,
	})
	status := 0
	if err != nil {
		status = -1
	}
	return &struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "ok",
		Status:  status,
	}, err
}

func apiFriendRequestApprove(ctx *ChatContext, data *struct {
	MessageId string `json:"message_id"` // 请求 ID
	Approve   bool   `json:"approve"`    // 是否通过
	Comment   string `json:"comment"`    // 备注信息
}) (any, error) {
	req, _ := model.FriendRequestGetByID(data.MessageId)
	if req.ID != "" {
		if req.ReceiverID != ctx.User.ID {
			// 无权审批
			return false, nil
		}
		if req.Status != "pending" {
			return false, nil
		}

		if model.FriendRequestSetApprove(req.ID, data.Approve) {
			// 建立联系
			ok, _ := model.FriendRelationFriendApprove(req.SenderID, req.ReceiverID)

			if ok {
				ch, _ := model.ChannelPrivateGet(req.SenderID, req.ReceiverID)
				if ch.ID == "" {
					model.ChannelPrivateNew(req.SenderID, req.ReceiverID)
				}
			}
			return ok, nil
		}

	}

	return false, nil
}

func apiFriendDelete(ctx *ChatContext, data *struct {
	UserId string `json:"user_id"` // 用户 ID
}) (any, error) {
	ok := model.FriendRelationDelete(data.UserId, ctx.User.ID)
	return ok, nil
}
