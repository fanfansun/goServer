package session

import (
	"github.com/fanfansun/goServer/api/db"
	"github.com/fanfansun/goServer/api/model"
	"github.com/fanfansun/goServer/api/utils"
	"sync"
	"time"
)

var sessionMap *sync.Map

// 初始化session_Map
func init() {
	sessionMap = &sync.Map{}
}

// 时间转换为字符串
func nowconvertString() int64 {
	return time.Now().UnixNano() / 100000
}

// 删除过期会话
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	db.DeleteSession(sid)
}

// LoadSessionsFromDB 加载所有 Session 并写入 Map
func LoadSessionsFromDB() {
	re, err := db.RetrieveAllSessions()
	if err != nil {
		return
	}
	re.Range(
		func(k, v interface{}) bool {
			ss := v.(*model.SimpleSession)
			sessionMap.Store(k, ss)
			print("sessionMap:--", sessionMap)
			return true
		})
}

// GenerateNewSessionId 生成新会话ID(令牌)
func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := nowconvertString()
	ttl := ct + 30*60*1000 // 30分钟的令牌
	/**
	 ttl时间比 ct采集时间 > 30 min
	当前的时间是不断向前变化的
	所以当ttl时间 < 时间当前时间  ,意味着ttl过期

	*/
	ss := &model.SimpleSession{Username: un, TTL: ttl}
	sessionMap.Store(id, ss)
	db.InsertSession(id, ttl, un)

	return id
}

// IsSessionExpired 会话是否过期 判断
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		nowCt := nowconvertString()
		// 过期
		if ss.(*model.SimpleSession).TTL < nowCt {
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*model.SimpleSession).Username, false
	}
	return "", true
}
