package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// ReadVideoDeletionRecord 读取视频删除记录
func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmtOut, err2 := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")
	var ids []string
	if err2 != nil {
		return ids, err
	}
	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %v", err)
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	defer stmtOut.Close()
	return ids, nil

}

// DelVideoDeletionRecord 删除视频删除记录
func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}

	defer stmtDel.Close()
	return nil
}
