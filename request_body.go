package slim

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/buexplain/go-slim/errors"
	"io/ioutil"
	"net/http"
)

func (this *Request) Body() ([]byte, error) {
	this.r.Body = http.MaxBytesReader(this.ctx.w.Raw(), this.r.Body, this.ctx.app.bodyMaxBytes)
	b, err := ioutil.ReadAll(this.r.Body)
	if err != nil {
		return nil, errors.MarkClient(err)
	}
	return b, nil
}

func (this *Request) JSON(v interface{}) error {
	if !this.IsJSON() {
		return errors.MarkClient(fmt.Errorf("must set Content-Type: application/json header"))
	}

	this.r.Body = http.MaxBytesReader(this.ctx.w.Raw(), this.r.Body, this.ctx.app.bodyMaxBytes)

	b, err := ioutil.ReadAll(this.r.Body)
	if err != nil {
		return errors.MarkClient(err)
	}

	if len(b) == 0 {
		return errors.MarkClient(fmt.Errorf("JSON payload is empty"))
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return errors.MarkClient(err)
	}

	return nil
}

func (this *Request) XML(v interface{}) error {
	if !this.IsXML() {
		return errors.MarkClient(fmt.Errorf("must set Content-Type: application/xml header"))
	}

	this.r.Body = http.MaxBytesReader(this.ctx.w.Raw(), this.r.Body, this.ctx.app.bodyMaxBytes)

	b, err := ioutil.ReadAll(this.r.Body)
	if err != nil {
		return errors.MarkClient(err)
	}

	if len(b) == 0 {
		return errors.MarkClient(fmt.Errorf("XML payload is empty"))
	}

	err = xml.Unmarshal(b, v)
	if err != nil {
		return errors.MarkClient(err)
	}

	return nil
}
