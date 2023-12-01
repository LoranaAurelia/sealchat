package api

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"sealchat/model"
	"sealchat/protocol"
)

func (ctx *ChatContext) TagCheck(ChannelID, msgId, text string) {
	root := protocol.ElementParse(text)
	db := model.GetDB()

	root.Traverse(func(el *protocol.Element) {
		switch el.Type {
		case "at":
			mention := model.MentionModel{
				StringPKBaseModel: model.StringPKBaseModel{
					ID: gonanoid.Must(),
				},
				SenderId:    ctx.User.ID,
				LocPostType: "channel",
				LocPostID:   ChannelID,
				RelatedType: "message",
				RelatedID:   msgId,
			}
			if el.Attrs["role"] == "all" {
				mention.ReceiverId = "all"
				db.Create(&mention)
			} else {
				if id, exists := el.Attrs["id"]; exists {
					mention.ReceiverId = id.(string)
					db.Create(&mention)
				}
			}
		}
	})
	//fmt.Println("xxx", root.ToString())
}
