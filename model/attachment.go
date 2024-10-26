package model

import (
	"encoding/hex"
	"encoding/json"

	"gorm.io/gorm"

	"sealchat/utils"
)

type ByteArray []byte

func (m ByteArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(m))
}

type AttachmentModel struct {
	StringPKBaseModel
	Hash      ByteArray `gorm:"index,size:100" json:"hash"` // hash是32byte
	Filename  string    `json:"filename"`
	Size      int64     `gorm:"index" json:"size"`
	UserID    string    `json:"userId" gorm:"index"`
	ChannelID string    `json:"channel_id"` // 上传的频道ID

	Extra string `json:"extra,omitempty"` // 额外标记
	Note  string `json:"note"`            // 另一个额外标记

	RootID       string `json:"rootId,omitempty" gorm:"index"`   // 相关的对象，用于后期检查不再使用的文件。使用这个给大的所属对象，例如项目中上传的所有附件，rootId都设置为项目id
	RootIDType   string `json:"rootIdType,omitempty"`            //
	ParentID     string `json:"parentId,omitempty" gorm:"index"` // 也是相关的对象，相当于第二槽位
	ParentIDType string `json:"parentIdType,omitempty"`          //

	IsTemp        bool   `json:"isTemp,omitempty" gorm:"index"` // 临时文件标记，先上传上来，无问题转正，有问题自动删除
	CreatorName   string `json:"creatorName,omitempty"`         // 上传者的名字
	CreatorAvatar string `json:"creatorAvatar,omitempty"`
}

func (*AttachmentModel) TableName() string {
	return "attachments"
}

func AttachmentCreate(at *AttachmentModel) (tx *gorm.DB, item *AttachmentModel) {
	db := GetDB()
	at.ID = utils.NewID()
	return db.Create(at), at
}

func AttachmentSetConfirm(ids []string, data map[string]any) (tx *gorm.DB) {
	item := &AttachmentModel{}
	m := map[string]any{
		"is_temp": false,
	}

	for key, value := range data {
		if value != "" {
			switch key {
			case "postIdType":
				m["post_id_type"] = value
			case "postId":
				m["post_id"] = value
			case "relatedPostIDType":
				m["related_post_id_type"] = value
			case "relatedPostID":
				m["related_post_id"] = value
			case "extra":
				m["extra"] = value
			case "note":
				m["note"] = value
			case "note2":
				m["note2"] = value
			case "isTemp":
				m["is_temp"] = value
			}
		}
	}

	q := db.Model(&item).
		Where("id IN (?)", ids).
		Updates(m)

	return q
}

// AttachmentsSetDelete 删除附件(注意，删除文件需要另外处理，id与hash为多对一关系)
func AttachmentsSetDelete(attachmentIdList []string) int64 {
	if len(attachmentIdList) > 0 {
		ret := db.Unscoped().Delete(&AttachmentModel{}, "id IN (?)", attachmentIdList)
		return ret.RowsAffected
	}
	return 0
}
