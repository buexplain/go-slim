package slim

import "net/http"

type Session interface {
	//从session中读取一个条目
	Get(k interface{}) interface{}
	GetString(k interface{}) string
	GetInt(k interface{}) int
	GetFloat64(k interface{}) float64
	GetFloat32(k interface{}) float32
	//从session中读取一个条目，并移除它
	Pull(k interface{}) interface{}
	PullString(k interface{}) string
	PullInt(k interface{}) int
	PullFloat64(k interface{}) float64
	PullFloat32(k interface{}) float32
	//给session增加一个条目
	Set(k, v interface{})
	//从session中移除一个条目
	Del(k interface{})
	//检查session里是否有此条目
	Has(k interface{}) bool
	//返回session id
	ID() string
	//返回session name
	Name() string
	//重新生成一个session id
	Regenerate()
	//销毁session
	Destroy()
	//session落地
	Save(r *http.Request, w http.ResponseWriter) error
}

type SessionHandler interface {
	Get(r *Request) (Session, error)
}
