package http

// HandleFunc 业务处理函数接口
type HandleFunc func(controller *GRegisterController) (err error)

// Filter 过滤器（拦截器），根据dispatch处理流程进行上下文拦截处理
type Filter func(controller *GRegisterController, f HandleFunc) (err error)

// NoopFilter 空Filter实现
func NoopFilter(controller *GRegisterController, f HandleFunc) (err error) {
	return f(controller)
}

// Chain 链式过滤器
type Chain []Filter

// Handle 链式过滤器递归处理流程
func (fc *Chain) Handle(controller *GRegisterController, f HandleFunc) (err error) {

	n := len(*fc)

	//无Filter,执行空Filter
	if n == 0 {
		return NoopFilter(controller, f)
	}

	//一个Filter,直接处理
	if n == 1 {
		return (*fc)[0](controller, f)
	}

	//多个Filter,递归执行
	lastI := n - 1
	return func(controller *GRegisterController, f HandleFunc) error {

		var (
			chainFunc HandleFunc
			curI      int
		)
		chainFunc = func(controller *GRegisterController) error {
			if curI == lastI {
				return f(controller)
			}
			curI++
			err := (*fc)[curI](controller, chainFunc)
			curI--
			return err
		}
		return (*fc)[0](controller, chainFunc)
	}(controller, f)
}

func (fc *Chain) AddFilter(f []Filter) {
	if f == nil {
		return
	}
	*fc = append(*fc, f...)
}
