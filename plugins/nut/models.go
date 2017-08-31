package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// RoleAdmin admin role
	RoleAdmin = "admin"
	// RoleRoot root role
	RoleRoot = "root"
	// UserTypeEmail email user
	UserTypeEmail = "email"

	// DefaultResourceType default resource type
	DefaultResourceType = "-"
	// DefaultResourceID default resourc id
	DefaultResourceID = 0
)

// User user
type User struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	UID             string     `json:"uid" gorm:"column:uid"`
	Password        []byte     `json:"-"`
	ProviderID      string     `json:"-"`
	ProviderType    string     `json:"providerType"`
	Home            string     `json:"home"`
	Logo            string     `json:"logo"`
	SignInCount     uint       `json:"signInCount"`
	LastSignInAt    *time.Time `json:"lastSignInAt"`
	LastSignInIP    string     `json:"lastSignInIp"`
	CurrentSignInAt *time.Time `json:"currentSignInAt"`
	CurrentSignInIP string     `json:"currentSignInIp"`
	ConfirmedAt     *time.Time `json:"confirmedAt"`
	LockedAt        *time.Time `json:"lockedAt"`

	Logs []Log `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (User) TableName() string {
	return "users"
}

// IsConfirm is confirm?
func (p *User) IsConfirm() bool {
	return p.ConfirmedAt != nil
}

// IsLock is lock?
func (p *User) IsLock() bool {
	return p.LockedAt != nil
}

//SetGravatarLogo set logo by gravatar
func (p *User) SetGravatarLogo() {
	buf := md5.Sum([]byte(strings.ToLower(p.Email)))
	p.Logo = fmt.Sprintf("https://gravatar.com/avatar/%s.png", hex.EncodeToString(buf[:]))
}

//SetUID generate uid
func (p *User) SetUID() {
	p.UID = uuid.New().String()
}

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

// Attachment attachment
type Attachment struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Length       int64  `json:"length"`
	MediaType    string `json:"mediaType"`
	ResourceID   uint   `json:"resourceId"`
	ResourceType string `json:"resourceType"`

	UserID uint `json:"userId"`
	User   User `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Attachment) TableName() string {
	return "attachments"
}

// IsPicture is picture?
func (p *Attachment) IsPicture() bool {
	return strings.HasPrefix(p.MediaType, "image/")
}

// Log log
type Log struct {
	ID uint `gorm:"primary_key" json:"id"`

	Message string `json:"message"`
	Type    string `json:"type"`
	IP      string `json:"ip"`

	UserID uint `json:"userId"`
	User   User `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (Log) TableName() string {
	return "logs"
}

func (p Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// Policy policy
type Policy struct {
	ID       uint `gorm:"primary_key" json:"id"`
	StartUp  time.Time
	ShutDown time.Time

	UserID uint
	User   User
	RoleID uint
	Role   Role

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// TableName table name
func (Policy) TableName() string {
	return "policies"
}

// Role role
type Role struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	Name         string `json:"name"`
	ResourceID   uint   `json:"resourceId"`
	ResourceType string `json:"resourceType"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Role) TableName() string {
	return "roles"
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// Vote vote
type Vote struct {
	ID           uint `gorm:"primary_key" json:"id"`
	Point        int
	ResourceID   uint
	ResourceType string

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Vote) TableName() string {
	return "votes"
}

// LeaveWord leave-word
type LeaveWord struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (LeaveWord) TableName() string {
	return "leave_words"
}

// Link link
type Link struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Loc       string    `json:"loc"`
	Href      string    `json:"href"`
	Label     string    `json:"label"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Link) TableName() string {
	return "links"
}

// Card card
type Card struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Loc       string    `json:"loc"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Href      string    `json:"href"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Card) TableName() string {
	return "cards"
}

// FriendLink friend_links
type FriendLink struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Title     string    `json:"title"`
	Home      string    `json:"home"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (FriendLink) TableName() string {
	return "friend_links"
}