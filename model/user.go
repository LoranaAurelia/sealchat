package model

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/blake2s"
	"gorm.io/gorm"

	"sealchat/protocol"
	"sealchat/utils"
)

// UserModel 用户表
type UserModel struct {
	StringPKBaseModel
	Nickname string `gorm:"null" json:"nick"` // 昵称
	Avatar   string `json:"avatar"`           // 头像
	Brief    string `json:"brief"`            // 简介
	// Role     string `json:"role"`             // 权限

	Username string `gorm:"index,unique;not null" json:"username"` // 用户名，唯一，非空
	Password string `gorm:"not null" json:"-"`                     // 密码，非空
	Salt     string `gorm:"not null" json:"-"`                     // 盐，非空
	IsBot    bool   `gorm:"null" json:"is_bot"`                    // 是否是机器人

	Disabled    bool              `json:"disabled"`
	AccessToken *AccessTokenModel `gorm:"-" json:"-"`

	RoleIds []string `json:"roleIds" gorm:"-"`
	// Token          string `gorm:"index" json:"token"` // 令牌
	// TokenExpiresAt int64  `json:"expiresAt"`
	// RecentSentAt int64 `json:"recentSentAt"` // 最近发送消息的时间
}

func (*UserModel) TableName() string {
	return "users"
}

func (u *UserModel) ToProtocolType() *protocol.User {
	return &protocol.User{
		ID:     u.ID,
		Nick:   u.Nickname,
		Avatar: u.Avatar,
		IsBot:  u.IsBot,
	}
}

func (u *UserModel) SaveAvatar() {
	db.Model(u).Update("avatar", u.Avatar)
}

func (u *UserModel) SaveInfo() {
	db.Model(u).Select("nickname", "brief").Updates(u)
}

// UserSetDisable 禁用用户函数
func UserSetDisable(userId string, val bool) error {
	return db.Model(&UserModel{}).Where("id = ?", userId).Update("disabled", val).Error
}

// AccessTokenModel access_token表
type AccessTokenModel struct {
	StringPKBaseModel
	UserID    string    `json:"userID" gorm:"not null"`    // 用户ID，非空
	ExpiredAt time.Time `json:"expiredAt" gorm:"not null"` // 过期时间，非空
}

func (*AccessTokenModel) TableName() string {
	return "access_tokens"
}

// 生成随机盐
func generateSalt() string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(salt)
}

// 使用盐对密码进行哈希
func hashPassword(password string, salt string) (string, error) {
	// 将密码和盐拼接起来
	saltedPassword := password + salt

	// 计算哈希值
	hashBytes := blake2s.Sum256([]byte(saltedPassword))
	hash := base64.RawStdEncoding.EncodeToString(hashBytes[:])

	return hash, nil
}

func UserCount() int64 {
	var count int64
	db.Select("id").Find(&UserModel{}).Count(&count)
	return count
}

// 创建用户
func UserCreate(username, password string, nickname string) (*UserModel, error) {
	salt := generateSalt()
	hashedPassword, err := hashPassword(password, salt)
	if err != nil {
		return nil, err
	}
	user := &UserModel{
		Username: username,
		Nickname: nickname,
		Password: hashedPassword,
		Salt:     salt,
	}
	user.ID = utils.NewID()
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 修改密码
func UserUpdatePassword(userID string, newPassword string) error {
	// 查询用户
	var user UserModel
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("UserModel not found")
	}

	// 更新密码
	salt := generateSalt()
	hashedNewPassword, err := hashPassword(newPassword, salt)
	if err != nil {
		return err
	}
	if err := db.Model(&user).Updates(map[string]interface{}{
		"password": hashedNewPassword,
		"salt":     salt,
	}).Error; err != nil {
		return err
	}
	return nil
}

// 登录认证
func UserAuthenticate(username, password string) (*UserModel, error) {
	var user UserModel
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	hashedPassword, err := hashPassword(password, user.Salt)
	if err != nil {
		return nil, err
	}
	if hashedPassword != user.Password {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}

func AcessTokenDeleteAllByUserID(userID string) error {
	return db.Where("user_id = ?", userID).Delete(&AccessTokenModel{}).Error
}

// UserGenerateAccessToken 生成 access_token
func UserGenerateAccessToken(userID string) (string, error) {
	expiredAt := time.Now().Add(time.Duration(15*24) * time.Hour)

	token := utils.NewID()
	accessToken := &AccessTokenModel{
		UserID:    userID,
		ExpiredAt: expiredAt,
	}

	accessToken.ID = token
	signedToken := TokenSign(accessToken.ID, expiredAt)
	if err := db.Create(accessToken).Error; err != nil {
		return "", err
	}
	return signedToken, nil
}

// UserVerifyAccessToken 验证 access_token 是否有效
func UserVerifyAccessToken(tokenString string) (*UserModel, error) {
	// 解析 token
	ret := TokenCheck(tokenString)

	if !ret.HashValid {
		return nil, ErrInvalidToken
	}

	if !ret.TimeValid {
		return nil, ErrTokenExpired
	}

	var accessToken AccessTokenModel
	if err := db.Where("id = ?", ret.Token).Limit(1).Find(&accessToken).Error; err != nil {
		return nil, ErrInvalidToken
	}

	if accessToken.ID == "" {
		return nil, ErrInvalidToken
	}

	now := time.Now()
	if accessToken.ExpiredAt.Compare(now) <= 0 {
		// 二次过期时间校验
		return nil, ErrInvalidToken
	}

	// 查询用户
	var user UserModel
	if err := db.Where("id = ?", accessToken.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.AccessToken = &accessToken
	return &user, nil
}

// UserRefreshAccessToken 刷新 access_token
func UserRefreshAccessToken(tokenID string) (string, error) {
	expiredAt := time.Now().Add(time.Duration(15*24) * time.Hour)

	var accessToken AccessTokenModel
	if err := db.Where("id = ?", tokenID).First(&accessToken).Error; err != nil {
		return "", ErrInvalidToken
	}

	if err := db.Model(&AccessTokenModel{}).Update("expired_at", expiredAt).Error; err != nil {
		return "", fmt.Errorf("update failed")
	}

	signedToken := TokenSign(accessToken.ID, expiredAt)
	return signedToken, nil
}

// UserGetEx 获取用户信息
func UserGetEx(id string) (*UserModel, error) {
	var user UserModel
	result := db.Where("id = ?", id).Limit(1).Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("获取用户信息失败: %v", result.Error)
	}
	return &user, nil
}

func UserGet(id string) *UserModel {
	r, _ := UserGetEx(id)
	return r
}

// UserBotList 查询所有启用的机器人用户
func UserBotList() ([]*UserModel, error) {
	var bots []*UserModel
	err := db.Where("disabled = ? AND is_bot = ?", false, true).Find(&bots).Error
	if err != nil {
		return nil, err
	}
	return bots, nil
}
