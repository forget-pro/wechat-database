package wechatdb

import (
	"context"
	"time"

	"github.com/sjzar/chatlog/internal/model"
	"github.com/sjzar/chatlog/internal/wechatdb/datasource"
	"github.com/sjzar/chatlog/internal/wechatdb/repository"

	_ "github.com/mattn/go-sqlite3"
)

// PaginationInfo 分页信息
type PaginationInfo struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

type DB struct {
	path     string
	platform string
	version  int
	account  string
	ds       datasource.DataSource
	repo     *repository.Repository
}

func New(path string, platform string, version int, account string) (*DB, error) {

	w := &DB{
		path:     path,
		platform: platform,
		version:  version,
		account:  account,
	}

	// 初始化，加载数据库文件信息
	if err := w.Initialize(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *DB) Close() error {
	if w.repo != nil {
		return w.repo.Close()
	}
	return nil
}

func (w *DB) Initialize() error {
	var err error
	w.ds, err = datasource.New(w.path, w.platform, w.version)
	if err != nil {
		return err
	}

	w.repo, err = repository.New(w.ds, w.account)
	if err != nil {
		return err
	}

	return nil
}

// GetMessagesResp 获取消息响应
type GetMessagesResp struct {
	Items []*model.Message `json:"items"`
	Total int              `json:"total"`
	Page  *PaginationInfo  `json:"page,omitempty"`
}

func (w *DB) GetMessages(start, end time.Time, talker string, sender string, keyword string, limit, offset int) ([]*model.Message, error) {
	ctx := context.Background()

	// 使用 repository 获取消息
	messages, err := w.repo.GetMessages(ctx, start, end, talker, sender, keyword, limit, offset)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessagesWithPagination 获取消息列表（带分页信息）
func (w *DB) GetMessagesWithPagination(start, end time.Time, talker string, sender string, keyword string, limit, offset int) (*GetMessagesResp, error) {
	ctx := context.Background()

	// 使用 repository 获取消息和总数
	messages, total, err := w.repo.GetMessagesWithTotal(ctx, start, end, talker, sender, keyword, limit, offset)
	if err != nil {
		return nil, err
	}

	var pageInfo *PaginationInfo
	if limit > 0 {
		hasMore := offset+len(messages) < total
		pageInfo = &PaginationInfo{
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		}
	}

	return &GetMessagesResp{
		Items: messages,
		Total: total,
		Page:  pageInfo,
	}, nil
}

type GetContactsResp struct {
	Items []*model.Contact `json:"items"`
	Total int              `json:"total"`
	Page  *PaginationInfo  `json:"page,omitempty"`
}

func (w *DB) GetContacts(key string, limit, offset int) (*GetContactsResp, error) {
	ctx := context.Background()

	contacts, total, err := w.repo.GetContactsWithTotal(ctx, key, limit, offset)
	if err != nil {
		return nil, err
	}

	var pageInfo *PaginationInfo
	if limit > 0 {
		hasMore := offset+len(contacts) < total
		pageInfo = &PaginationInfo{
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		}
	}

	return &GetContactsResp{
		Items: contacts,
		Total: total,
		Page:  pageInfo,
	}, nil
}

type GetChatRoomsResp struct {
	Items []*model.ChatRoom `json:"items"`
	Total int               `json:"total"`
	Page  *PaginationInfo   `json:"page,omitempty"`
}

func (w *DB) GetChatRooms(key string, limit, offset int) (*GetChatRoomsResp, error) {
	ctx := context.Background()

	chatRooms, total, err := w.repo.GetChatRoomsWithTotal(ctx, key, limit, offset)
	if err != nil {
		return nil, err
	}

	var pageInfo *PaginationInfo
	if limit > 0 {
		hasMore := offset+len(chatRooms) < total
		pageInfo = &PaginationInfo{
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		}
	}

	return &GetChatRoomsResp{
		Items: chatRooms,
		Total: total,
		Page:  pageInfo,
	}, nil
}

type GetSessionsResp struct {
	Items []*model.Session `json:"items"`
}

func (w *DB) GetSessions(key string, limit, offset int) (*GetSessionsResp, error) {
	ctx := context.Background()

	// 使用 repository 获取会话列表
	sessions, err := w.repo.GetSessions(ctx, key, limit, offset)
	if err != nil {
		return nil, err
	}

	return &GetSessionsResp{
		Items: sessions,
	}, nil
}

func (w *DB) GetMedia(_type string, key string) (*model.Media, error) {
	return w.repo.GetMedia(context.Background(), _type, key)
}
