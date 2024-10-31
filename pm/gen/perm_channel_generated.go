// Code generated by go generate; DO NOT EDIT.
package gen

// PermChannelMap 提供权限字符串到描述的映射
var PermChannelMap = map[string]string{
    "func_channel_read": "频道 - 消息 - 查看",
    "func_channel_text_send": "频道 - 消息 - 文本发送",
    "func_channel_file_send": "频道 - 消息 - 文件发送",
    "func_channel_audio_send": "频道 - 消息 - 音频发送",
    "func_channel_invite": "频道 - 常规 - 邀请加入频道",
    "func_channel_sub_channel_create": "频道 - 常规 - 创建子频道",
    "func_channel_member_remove": "频道 - 频道管理 - 踢人",
    "func_channel_manage_mute": "频道 - 频道管理 - 禁言",
    "func_channel_read_all": "频道 - 特殊 - 查看所有子频道",
    "func_channel_text_send_all": "频道 - 特殊 - 在所有子频道发送文本",
    "func_channel_role_link": "频道 - 成员管理 - 添加角色",
    "func_channel_role_unlink": "频道 - 成员管理 - 移除角色",
    "func_channel_role_link_root": "频道 - 成员管理 - 添加角色 (Root管理员)",
    "func_channel_role_unlink_root": "频道 - 成员管理 - 移除角色 (Root管理员)",
    "func_channel_manage_info": "频道 - 频道设置 - 基础设置",
    "func_channel_manage_role": "频道 - 频道设置 - 权限管理",
    "func_channel_manage_role_root": "频道 - 频道设置 - 权限管理（Root管理员）",
}


// PermChannelArray 提供权限字符串到描述的映射
var PermChannelArray = []map[string]string{
	{"key": "func_channel_read", "desc": "频道 - 消息 - 查看"},
	{"key": "func_channel_text_send", "desc": "频道 - 消息 - 文本发送"},
	{"key": "func_channel_file_send", "desc": "频道 - 消息 - 文件发送"},
	{"key": "func_channel_audio_send", "desc": "频道 - 消息 - 音频发送"},
	{"key": "func_channel_invite", "desc": "频道 - 常规 - 邀请加入频道"},
	{"key": "func_channel_sub_channel_create", "desc": "频道 - 常规 - 创建子频道"},
	{"key": "func_channel_member_remove", "desc": "频道 - 频道管理 - 踢人"},
	{"key": "func_channel_manage_mute", "desc": "频道 - 频道管理 - 禁言"},
	{"key": "func_channel_read_all", "desc": "频道 - 特殊 - 查看所有子频道"},
	{"key": "func_channel_text_send_all", "desc": "频道 - 特殊 - 在所有子频道发送文本"},
	{"key": "func_channel_role_link", "desc": "频道 - 成员管理 - 添加角色"},
	{"key": "func_channel_role_unlink", "desc": "频道 - 成员管理 - 移除角色"},
	{"key": "func_channel_role_link_root", "desc": "频道 - 成员管理 - 添加角色 (Root管理员)"},
	{"key": "func_channel_role_unlink_root", "desc": "频道 - 成员管理 - 移除角色 (Root管理员)"},
	{"key": "func_channel_manage_info", "desc": "频道 - 频道设置 - 基础设置"},
	{"key": "func_channel_manage_role", "desc": "频道 - 频道设置 - 权限管理"},
	{"key": "func_channel_manage_role_root", "desc": "频道 - 频道设置 - 权限管理（Root管理员）"},
}