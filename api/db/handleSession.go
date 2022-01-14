package db

import (
	"database/sql"
	"github.com/fanfansun/goServer/api/model"
	"log"
	"strconv"
	"sync"
)

func InsertSession(sid string, ttl int64, uname string) error {
	// TTL 转为 str 类型
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare("INSERT INTO sessions (session_id, TTL, login_name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

// 从数据库中检索某个会话
func RetrieveSession(sid string) (*model.SimpleSession, error) {
	ss := &model.SimpleSession{}
	stmtOut, err := dbConn.Prepare("SELECT TTL, login_name FROM sessions WHERE session_id=?")
	if err != nil {
		return nil, err
	}

	var ttl string
	var uname string
	stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// TTL 转为 int64
	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.Username = uname
	} else {
		return nil, err
	}

	defer stmtOut.Close()
	return ss, nil
}

// RetrieveAllSessions 检索所有会话
/**
  @*sync.Map 把session 放到 map队列中返回
  map 充当了 Catce
*/
func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() {
		var id string
		var ttlstr string
		var login_name string
		if err := rows.Scan(&id, &ttlstr, &login_name); err != nil {
			log.Printf("retrive sessions error: %s", err)
			break
		}

		if ttl, err1 := strconv.ParseInt(ttlstr, 10, 64); err1 == nil {
			ss := &model.SimpleSession{Username: login_name, TTL: ttl}
			m.Store(id, ss)
			log.Printf(" session id: %s, ttl: %d", id, ss.TTL)
		}

	}

	return m, nil
}

// 从数据库中删除 session
func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}

	return nil
}
