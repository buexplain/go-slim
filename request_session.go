package slim

import "fmt"

func (this *Request) Session() Session {
	if this.session == nil {
		if s, err := this.ctx.app.sessionHandler.Get(this); err != nil {
			panic(fmt.Errorf("get session error: %w", err).Error())
		} else {
			this.session = s
		}
	}
	return this.session
}
