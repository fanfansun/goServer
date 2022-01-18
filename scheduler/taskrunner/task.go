package taskrunner

import (
	"errors"
	"github.com/fanfansun/goServer/scheduler/db"
	"log"
	"os"
	"sync"
)

func deleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + vid)

	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
		return err
	}

	return nil
}

// VideoClearDispatcher
// 发布视频删除任务
func VideoClearDispatcher(dc dataChan) error {
	//  返回视频删除表里最近n条记录
	res, err := db.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}
	// 没有可以删除的视频
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}

	for _, id := range res {
		dc <- id
	}

	return nil
}

// VideoClearExecutor 执行删除视频的任务
func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error

forloop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				// 硬盘删除视频
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				// 数据库删除记录
				if err := db.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})

	return err
}
