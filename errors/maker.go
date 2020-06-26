package errors

import "errors"

type MrKErr struct {
	err  error
	code int
}

//给错误标记上一个code码
func Mark(err error, code int) error {
	if err == nil {
		return nil
	}
	var m *MrKErr
	if errors.As(err, &m) {
		if m.code != code {
			m.code = code
		}
		return err
	}
	return &MrKErr{err: err, code: code}
}

//尝试给错误标记上一个code码，如果错误已经被标记，则不在标记
func TryMark(err error, code int) error {
	if err == nil {
		return nil
	}
	var m *MrKErr
	if errors.As(err, &m) {
		return err
	}
	return &MrKErr{err: err, code: code}
}

func (this MrKErr) Code() int {
	return this.code
}

func (this MrKErr) Error() string {
	return this.err.Error()
}

func (e *MrKErr) Unwrap() error {
	return e.err
}

func (this MrKErr) Is(target error) bool {
	return errors.Is(this.err, target)
}

func (this MrKErr) As(target interface{}) bool {
	return errors.As(this.err, &target)
}

func IsMarker(err error) *MrKErr {
	if err == nil {
		return nil
	}
	var m *MrKErr
	if errors.As(err, &m) {
		return m
	}
	return nil
}

func HasMaker(err error, code int) bool {
	if err == nil {
		return false
	}
	var m *MrKErr
	if errors.As(err, &m) && m.code == code {
		return true
	}
	return false
}
