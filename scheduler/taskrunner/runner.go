package taskrunner

// Runner 死循坏调度任务实现， chan通道通信
type Runner struct {
	Controller controlChan // 控制通道
	Error      controlChan // 控制通道 错误
	Data       dataChan    // 数据 （chan interface 格式） 可以是任意数据
	dataSize   int         // 数据大小
	longLived  bool        // 是否延长（生存时间）
	Dispatcher fn          // 调度函数
	Executor   fn          // 执行者
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		longLived:  longlived,
		dataSize:   size,
		Dispatcher: d, // 准备派遣
		Executor:   e, // 执行者
	}
}

// 开始派遣
func (r *Runner) startDispatch() {
	// 函数异常退出执行
	defer func() {
		// 如果r,没有开启延长生存时间，关闭三个 chan 通道
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for {
		select {
		// 先判断 r.Controller ，把 r.Controller 写入 c
		case c := <-r.Controller:
			// 如果 c 状态是 READY_TO_DISPATCH 准备派遣
			if c == READY_TO_DISPATCH {
				// r 调用派遣函数, 并传入数据
				err := r.Dispatcher(r.Data)
				// 异常处理 如果有错误
				if err != nil {
					// close定义的字符写入 chan（通道错误） r.error
					r.Error <- CLOSE
				} else {
					// READY_TO_EXECUTE  (准备执行状态写入)
					r.Controller <- READY_TO_EXECUTE
				}
			}
			// 如果 c 是准备执行状态
			if c == READY_TO_EXECUTE {
				// 调用执行函数执行
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					// 成功执行后，状态变成准备派遣状态
					r.Controller <- READY_TO_DISPATCH
				}
			}
		// 如果有错误被捕获
		case e := <-r.Error:
			// 状态变成关闭
			if e == CLOSE {
				return
			}
		default:

		}
	}
}

// 初始化运行

func (r *Runner) StartAll() {
	// 派遣状态写入
	r.Controller <- READY_TO_DISPATCH
	// 调用上面实现的状态函数
	r.startDispatch()
}
