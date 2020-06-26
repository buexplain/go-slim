package errors

const (
	//客户端错误码
	ClientCode = 400
	//服务端错误码
	ServerCode = 500
)

//将错误标记为服务端错误
func MarkServer(err error) error {
	return Mark(err, ServerCode)
}

//尝试将错误标记为服务端错误
func TryMarkServer(err error) error {
	return TryMark(err, ServerCode)
}

//将错误标记为客户端错误
func MarkClient(err error) error {
	return Mark(err, ClientCode)
}

//尝试将错误标记为客户端错误
func TryMarkClient(err error) error {
	return TryMark(err, ClientCode)
}

//判断是否为服务端错误
func HasMarkerServer(err error) bool {
	return HasMaker(err, ServerCode)
}

//判断是否为客户端错误
func HasMarkerClient(err error) bool {
	return HasMaker(err, ClientCode)
}
