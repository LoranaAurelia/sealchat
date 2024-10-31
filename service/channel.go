package service

import (
	"fmt"
	"strings"

	"github.com/mikespook/gorbac"

	"github.com/samber/lo"

	"sealchat/model"
	"sealchat/pm"
)

// ChannelIdList 获取可见的频道ID，这个函数是下面的修改版，理论上会更精确，等待实际验证。可能有些调用代价，后面可以考虑使用memoize，也可能使用层级权限是更好的方式。
func ChannelIdList(userId string) ([]string, error) {
	// 包括如下内容:
	// 1. 属性为可见的一级频道(即没有父级的频道)
	// 2. 具有明确可看权限的频道(先查频道角色，再根据频道角色验证权限和获取频道id)
	// 3. 补入有权查看的频道的子频道

	roles, err := model.UserRoleMappingListByUserID(userId, "", "channel")
	if err != nil {
		return nil, err
	}

	var rolesCanRead []string
	db := model.GetDB()
	db.Model(&model.RolePermissionModel{}).
		Where("role_id in ? and permission_id in ?", roles, []string{pm.PermFuncChannelRead.ID(), pm.PermFuncChannelReadAll.ID()}).
		Pluck("role_id", &rolesCanRead)

	// 获得1 公开的一级频道
	var idsPublic1 []string
	db.Model(&model.ChannelModel{}).Where("coalesce(root_id, '') = '' and perm_type = ?", "public").
		Pluck("id", &idsPublic1)

	// 这里获得的是2: 具有明确可看权限的频道，包括公开频道和非公开频道
	ids2 := lo.Map(rolesCanRead, func(item string, index int) string {
		return strings.SplitN(item, "-", 3)[1]
	})

	// 将公开一级频道和有权限的频道组合起来
	idsCanRead := append(idsPublic1, ids2...)

	// 值得注意，ids2里可能混合了空中楼阁子频道，也就是说你没有他上级频道的权限
	// 要在之后进行剔除。虽然目前版本不支持2级以上频道，所以理论上不会存在

	// 3.1: 在可访问频道的基础上进一步加入公开的子频道
	var ids3 []string
	db.Model(&model.ChannelModel{}).Where("root_id in ? and perm_type = ?", idsCanRead, "public").
		Pluck("id", &ids3)

	// 3.2
	// 先找出我有“查看全部”权限的的顶级频道
	// 找出这些顶级频道的下属非公开频道
	var rolesCanRead2 []string
	db.Model(&model.RolePermissionModel{}).
		Where("role_id in ? and permission_id in ?", roles, []string{pm.PermFuncChannelReadAll.ID()}).
		Pluck("role_id", &rolesCanRead2)
	ids2x := lo.Map(rolesCanRead2, func(item string, index int) string {
		return strings.SplitN(item, "-", 3)[1]
	})
	var ids32 []string
	db.Model(&model.ChannelModel{}).Where("root_id in ? and perm_type = ?", ids2x, "non-public").
		Pluck("id", &ids32)

	idsCanRead = append(idsCanRead, ids3...)
	idsCanRead = append(idsCanRead, ids32...)

	// 对idsCanRead进行去重
	idsCanRead = lo.Uniq(idsCanRead)

	// 剔除父频道不在可读列表中的频道，但保留顶级频道
	var idsParentNotInCanRead []string
	db.Model(&model.ChannelModel{}).
		Where("id in ? and coalesce(parent_id,'') != '' and parent_id not in ?", idsCanRead, idsCanRead).
		Pluck("id", &idsParentNotInCanRead)

	idsCanRead = lo.Filter(idsCanRead, func(id string, _ int) bool {
		return !lo.Contains(idsParentNotInCanRead, id)
	})

	return idsCanRead, nil
}

// CanReadChannelByUserId 注意性能比较差，后面修改
func CanReadChannelByUserId(userId, channelId string) bool {
	chIds, _ := ChannelIdList(userId)
	return lo.Contains(chIds, channelId)
}

// ChannelList 获取可见的频道（等待重构）
func ChannelList(userId string) ([]*model.ChannelModel, error) {
	// 包括如下内容:
	// 1. 属性为可见的一级频道(即没有父级的频道)
	// 2. 具有明确可看权限的频道(先查频道角色，再根据频道角色验证权限和获取频道id)
	// 3. 补入有权查看的频道的子频道

	roles, err := model.UserRoleMappingListByUserID(userId, "", "channel")
	if err != nil {
		return nil, err
	}

	var rolesCanRead []string
	db := model.GetDB()
	db.Model(&model.RolePermissionModel{}).
		Where("role_id in ? and permission_id in ?", roles, []string{pm.PermFuncChannelRead.ID(), pm.PermFuncChannelReadAll.ID()}).
		Pluck("role_id", &rolesCanRead)

	// 这里获得的是2
	ids2 := lo.Map(rolesCanRead, func(item string, index int) string {
		return strings.SplitN(item, "-", 3)[1]
	})

	// 3.1，公开子频道
	var ids3 []string
	db.Model(&model.ChannelModel{}).Where("root_id in ? and perm_type = ?", ids2, "public").
		Pluck("id", &ids3)

	// 3.2
	// 先找出我有“查看全部”权限的的顶级频道
	// 找出这些顶级频道的下属频道
	var rolesCanRead2 []string
	db.Model(&model.RolePermissionModel{}).
		Where("role_id in ? and permission_id in ?", roles, []string{pm.PermFuncChannelReadAll.ID()}).
		Pluck("role_id", &rolesCanRead2)
	ids2x := lo.Map(rolesCanRead2, func(item string, index int) string {
		return strings.SplitN(item, "-", 3)[1]
	})
	var ids32 []string
	db.Model(&model.ChannelModel{}).Where("root_id in ? and perm_type = ?", ids2x, "non-public").
		Pluck("id", &ids32)

	idsAll := append(ids2, ids3...)
	idsAll = append(idsAll, ids32...)
	var items []*model.ChannelModel
	db.Model(&model.ChannelModel{}).Where("id in ? or perm_type = ?", idsAll, "public").
		Group("id").              // 使用Group来去重
		Order("sort_order DESC"). // 先按优先级排序(数字大的在前)
		Order("created_at DESC"). // 同优先级按创建时间降序
		Find(&items)

	return items, nil
}

func ChannelNew(channelID, channelType, channelName string, creatorId string, parentId string) *model.ChannelModel {
	m := model.ChannelPublicNew(channelID, &model.ChannelModel{
		Name:     channelName,
		PermType: channelType,
		ParentID: parentId,
		RootId:   parentId, // TODO: 这个是不准的，但是目前不允许二级以上子频道
	}, creatorId)

	roleCreate(channelID, "owner", "群主", func(roleId string) []gorbac.Permission {
		return []gorbac.Permission{
			pm.PermFuncChannelRead,
			pm.PermFuncChannelTextSend,
			pm.PermFuncChannelFileSend,
			pm.PermFuncChannelAudioSend,
			pm.PermFuncChannelInvite,
			// pm.PermFuncChannelMemberRemove,
			pm.PermFuncChannelSubChannelCreate,
			pm.PermFuncChannelRoleLink,
			pm.PermFuncChannelRoleUnlink,
			pm.PermFuncChannelRoleLinkRoot,
			pm.PermFuncChannelRoleUnlinkRoot,
			pm.PermFuncChannelManageInfo,
			pm.PermFuncChannelManageRole,
			pm.PermFuncChannelManageRoleRoot,
			pm.PermFuncChannelManageMute,
			pm.PermFuncChannelReadAll,
			pm.PermFuncChannelTextSendAll,
		}
	})

	roleCreate(channelID, "admin", "管理员", func(roleId string) []gorbac.Permission {
		return []gorbac.Permission{
			pm.PermFuncChannelRead,
			pm.PermFuncChannelTextSend,
			pm.PermFuncChannelFileSend,
			pm.PermFuncChannelAudioSend,
			pm.PermFuncChannelInvite,
			// pm.PermFuncChannelMemberRemove,
			pm.PermFuncChannelSubChannelCreate,
			pm.PermFuncChannelRoleLink,
			pm.PermFuncChannelRoleUnlink,
			pm.PermFuncChannelReadAll,
			pm.PermFuncChannelManageInfo,
			pm.PermFuncChannelManageRole,
			pm.PermFuncChannelManageMute,
			pm.PermFuncChannelTextSendAll,
		}
	})

	roleCreate(channelID, "ob", "观察者", func(roleId string) []gorbac.Permission {
		return []gorbac.Permission{
			pm.PermFuncChannelRead,
			pm.PermFuncChannelTextSend,
			pm.PermFuncChannelFileSend,
			pm.PermFuncChannelAudioSend,
			pm.PermFuncChannelReadAll,
		}
	})

	roleCreate(channelID, "visitor", "游客", func(roleId string) []gorbac.Permission {
		return []gorbac.Permission{
			pm.PermFuncChannelRead,
			pm.PermFuncChannelTextSend,
		}
	})

	roleCreate(channelID, "member", "成员", func(roleId string) []gorbac.Permission {
		return []gorbac.Permission{
			pm.PermFuncChannelRead,
			pm.PermFuncChannelTextSend,
			pm.PermFuncChannelFileSend,
			pm.PermFuncChannelAudioSend,
			pm.PermFuncChannelInvite,
		}
	})

	roleCreate(channelID, "bot", "机器人", func(roleId string) []gorbac.Permission {
		return []gorbac.Permission{
			pm.PermFuncChannelReadAll,
			pm.PermFuncChannelTextSendAll,
		}
	})

	roleId := fmt.Sprintf("ch-%s-%s", channelID, "owner")
	_ = model.UserRoleMappingCreate(&model.UserRoleMappingModel{
		UserID:   creatorId,
		RoleID:   roleId,
		RoleType: "channel",
	})

	return m
}
