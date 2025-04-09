package net

import (
	"crypto/cipher"
	"sync/atomic"
)

type Session interface {
	Encryptor

	UID() int64
	SID() int64
	Color() string
	Status() int64
	StartTime() int64

	ClientIP() string
	SetClientIP(ip string)

	CSIndex() int64
	SCIndex() int64
	IncreaseCSIndex() int64
	IncreaseSCIndex() int64
}

type Encryptor interface {
	IsCrypto() bool
	Block() cipher.Block
	Key() []byte
}

var _ Session = (*session)(nil)

type session struct {
	*encryptor

	userId    int64
	serverId  int64
	clientIP  string
	color     string
	status    int64
	startTime int64

	csIndex *indexInfo
	scIndex *indexInfo
}

func DefaultSession() Session {
	return &session{
		encryptor: &encryptor{},
		csIndex:   newIndexInfo(0),
		scIndex:   newIndexInfo(1),
	}
}

func NewSession(userId int64, sid int64, st int64, block cipher.Block, key []byte,
	crypto bool, color string, status int64) Session {
	s := &session{
		encryptor: &encryptor{
			encrypt: crypto,
			block:   block,
			key:     key,
		},
		userId:    userId,
		color:     color,
		status:    status,
		serverId:  sid,
		startTime: st,
		csIndex:   newIndexInfo(0),
		scIndex:   newIndexInfo(1),
	}
	return s
}

func (s *session) IncreaseCSIndex() int64 {
	return s.csIndex.Increase()
}

func (s *session) CSIndex() int64 {
	return s.csIndex.Load()
}

func (s *session) IncreaseSCIndex() int64 {
	return s.scIndex.Increase()
}

func (s *session) SCIndex() int64 {
	return s.scIndex.Load()
}

func (s *session) StartTime() int64 {
	return s.startTime
}

func (s *session) UID() int64 {
	return s.userId
}

func (s *session) SID() int64 {
	return s.serverId
}

func (s *session) Color() string {
	if len(s.color) == 0 {
		return ""
	}
	return s.color
}

func (s *session) Status() int64 {
	return s.status
}

func (s *session) ClientIP() string {
	return s.clientIP
}

func (s *session) SetClientIP(ip string) {
	s.clientIP = ip
}

type indexInfo struct {
	start int64
	index atomic.Int64
}

func newIndexInfo(start int64) *indexInfo {
	i := &indexInfo{
		start: start,
	}
	i.index.Store(start)
	return i
}

func (i *indexInfo) Increase() int64 {
	return i.index.Add(1)
}

func (i *indexInfo) Load() int64 {
	return i.index.Load()
}

var _ Encryptor = (*encryptor)(nil)

type encryptor struct {
	encrypt bool
	block   cipher.Block
	key     []byte
}

func (c *encryptor) IsCrypto() bool {
	return c.encrypt
}

func (c *encryptor) Block() cipher.Block {
	return c.block
}

func (c *encryptor) Key() []byte {
	return c.key
}
