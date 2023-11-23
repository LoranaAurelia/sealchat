package model

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"sealchat/protocol"
)

type MemberModel struct {
	StringPKBaseModel
	Nickname  string `gorm:"null" json:"nick"`           // 昵称
	ChannelID string `gorm:"not null" json:"channel_id"` // 频道ID
	UserID    string `json:"user_id" gorm:"null"`        // 用户ID
}

func (*MemberModel) TableName() string {
	return "members"
}

func (u *MemberModel) ToProtocolType() *protocol.GuildMember {
	return &protocol.GuildMember{
		Nick: u.Nickname,
	}
}

func MemberGetByUserIDAndChannelID(userId string, channelId string, defaultName string) (*MemberModel, error) {
	db := GetDB()
	var member MemberModel
	err := db.Where("user_id = ? AND channel_id = ?", userId, channelId).First(&member).Error
	if err != nil {
		// 未找到记录，尝试创建新的记录
		x := MemberModel{StringPKBaseModel: StringPKBaseModel{ID: gonanoid.Must()}, UserID: userId, ChannelID: channelId, Nickname: defaultName}
		err = db.Create(&x).Error
		return &x, err
	}
	return &member, nil
}