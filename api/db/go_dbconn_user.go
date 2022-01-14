package db

import (
	"database/sql"
	"github.com/fanfansun/goServer/api/model"
	"github.com/fanfansun/goServer/api/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

/**
import 下划线（如：import _ github/demo）的作用：
当导入一个包时，该包下的文件里所有init()函数都会被执行，
然而，有些时候我们并不需要把整个包都导入进来，
仅仅是是希望它执行init()函数而已。这个时候就可以使用 import _ 引用该包。
*/

// 打开一个连接
//func openDBconnenr() {
//	_, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/video_server?useSSL=false&serverTimezone=UTC&nullCatalogMeansCurrent=true")
//	if err != nil {
//		fmt.Printf("connnect fail%s\n", err)
//		panic(err.Error())
//	}
//	fmt.Printf("connnect success")
//}

func AddUserInformation(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)")
	if err != nil {
		return err
	}
	// 向数据库写入数据
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	/*defer 用来err时候结束流*/
	defer stmtIns.Close()
	return nil
}

func GetUserInformation(loginName string) (string, error) {
	//  // 因为查询单条数据时, 可能返回var ErrNoRows = errors.New("sql: no rows in result set")该种错误信息,没有赋值
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name=?")
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%s", err)
		return "", err
	}
	var pwd string
	stmtOut.QueryRow(loginName).Scan(&pwd)
	stmtOut.Close()
	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd=?")
	if err != nil {
		log.Printf("DeleteUser error %s", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

/*视频数据库部分*/

func AddvideoInformation(aid int, name string) (*model.VideoInfo, error) {
	// 创建uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:06") // M D Y : HH : MM :ss

	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info (id, author_id, name, display_ctime) VALUES(?, ?, ?, ?) `)
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}
	res := &model.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	defer stmtIns.Close()
	return res, nil

}

func GetvideoInformation(vid string) (*model.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")

	var aid int
	var dct string
	var name string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()

	res := &model.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}

	return res, nil
}

func Deletevideo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil

}

/*comment 评论部分*/

func AddcommentInformation(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmtIns, err := dbConn.Prepare("INSERT INTO comments (id, video_id, author_id, content) values (?, ?, ?, ?)")
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, to int) ([]*model.Comment, error) {
	stmtOut, err := dbConn.Prepare(`
        SELECT comments.id, users.Login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ? 
		  AND comments.time > FROM_UNIXTIME(?) 
		  AND comments.time <= FROM_UNIXTIME(?)`,
	)

	var res []*model.Comment

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &model.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()

	return res, nil
}
