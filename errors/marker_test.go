package errors

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

//测试对nil进行标记的情况
func TestMarkNil(t *testing.T) {
	var err error
	markErr := Mark(err, 1)
	if markErr != nil {
		t.Fatal("TestMarkNil fatal")
	}
	markErr = Mark(fmt.Errorf("more info: %w", err), 1)
	if errors.Is(markErr, err) {
		t.Fatal("TestMarkNil Is nil fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkNil HasMaker nil fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkNil IsMarker nil fatal")
	}
}

//测试对错误进行标记
func TestMark(t *testing.T) {
	//测试哨兵类型的错误
	markErr := Mark(io.EOF, 1)
	if errors.Is(markErr, io.EOF) == false {
		t.Fatal("TestMark Is io.EOF fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMark HasMaker io.EOF fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMark IsMarker io.EOF fatal")
	}

	//测试自定义类型的错误
	rawErr := errors.New("test raw error")
	markErr = Mark(rawErr, 1)
	if errors.Is(markErr, rawErr) == false {
		t.Fatal("TestMark Is errors.New fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMark HasMaker errors.New fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMark IsMarker errors.New fatal")
	}

	//测试实现了 error 接口的错误
	pathErr := &os.PathError{Op: "op", Path: "Path", Err: errors.New("Err")}
	markErr = Mark(pathErr, 1)
	if errors.Is(markErr, pathErr) == false {
		t.Fatal("TestMark Is os.PathError fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMark HasMaker os.PathError fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMark IsMarker os.PathError fatal")
	}

	var rawPathErr *os.PathError
	if errors.As(markErr, &rawPathErr) {
		fmt.Printf("%+v\n", rawPathErr)
	} else {
		t.Fatal("TestMark As os.PathError fatal")
	}
}

//测试对被包裹一次的错误进行标记
func TestMarkErrorf(t *testing.T) {
	//测试哨兵类型的错误
	markErr := Mark(fmt.Errorf("more info: %w", os.ErrExist), 1)
	if errors.Is(markErr, os.ErrExist) == false {
		t.Fatal("TestMarkErrorf Is os.ErrExist fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorf HasMaker os.ErrExist fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorf IsMarker os.ErrExist fatal")
	}

	//测试自定义类型的错误
	rawErr := errors.New("test raw error")
	markErr = Mark(fmt.Errorf("more info: %w", rawErr), 1)
	if errors.Is(markErr, rawErr) == false {
		t.Fatal("TestMarkErrorf Is errors.New fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorf HasMaker errors.New fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorf IsMarker errors.New fatal")
	}

	//测试实现了 error 接口的错误
	pathErr := &os.PathError{Op: "op", Path: "Path", Err: errors.New("Err")}
	markErr = Mark(fmt.Errorf("more info: %w", pathErr), 1)
	if errors.Is(markErr, pathErr) == false {
		t.Fatal("TestMarkErrorf Is os.PathError fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorf HasMaker os.PathError fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorf IsMarker os.PathError fatal")
	}

	var rawPathErr *os.PathError
	if errors.As(markErr, &rawPathErr) {
		fmt.Printf("%+v\n", rawPathErr)
	} else {
		t.Fatal("TestMarkErrorf As os.PathError fatal")
	}
}

//测试对被标记的错误进行包裹情况
func TestErrorfMark(t *testing.T) {
	//测试哨兵类型的错误
	markErr := fmt.Errorf("more info: %w", Mark(os.ErrExist, 1))
	if errors.Is(markErr, os.ErrExist) == false {
		t.Fatal("TestErrorfMark Is os.ErrExist fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestErrorfMark HasMaker os.ErrExist fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestErrorfMark IsMarker os.ErrExist fatal")
	}

	//测试自定义类型的错误
	rawErr := errors.New("test raw error")
	markErr = fmt.Errorf("more info: %w", Mark(rawErr, 1))
	if errors.Is(markErr, rawErr) == false {
		t.Fatal("TestErrorfMark Is errors.New fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestErrorfMark HasMaker errors.New fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestErrorfMark IsMarker errors.New fatal")
	}

	//测试实现了 error 接口的错误
	pathErr := &os.PathError{Op: "op", Path: "Path", Err: errors.New("Err")}
	markErr = fmt.Errorf("more info: %w", Mark(pathErr, 1))
	if errors.Is(markErr, pathErr) == false {
		t.Fatal("TestErrorfMark Is os.PathError fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestErrorfMark HasMaker os.PathError fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestErrorfMark IsMarker os.PathError fatal")
	}

	var rawPathErr *os.PathError
	if errors.As(markErr, &rawPathErr) {
		fmt.Printf("%+v\n", rawPathErr)
	} else {
		t.Fatal("TestErrorfMark As os.PathError fatal")
	}
}

//测试对错误的多次标记
func TestMarkMore(t *testing.T) {
	//测试哨兵类型的错误
	markErr := Mark(Mark(io.EOF, 1), 2)
	if errors.Is(markErr, io.EOF) == false {
		t.Fatal("TestMarkMore Is io.EOF fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 2) == false {
		t.Fatal("TestMarkMore HasMaker io.EOF fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 2 {
		t.Fatal("TestMarkMore IsMarker io.EOF fatal")
	}

	//测试自定义类型的错误
	rawErr := errors.New("test raw error")
	markErr = Mark(Mark(rawErr, 1), 1)
	if errors.Is(markErr, rawErr) == false {
		t.Fatal("TestMarkMore Is errors.New fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkMore HasMaker errors.New fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkMore IsMarker errors.New fatal")
	}

	//测试实现了 error 接口的错误
	pathErr := &os.PathError{Op: "op", Path: "Path", Err: errors.New("Err")}
	markErr = Mark(Mark(pathErr, 1), 2)
	if errors.Is(markErr, pathErr) == false {
		t.Fatal("TestMarkMore Is os.PathError fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 2) == false {
		t.Fatal("TestMarkMore HasMaker os.PathError fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 2 {
		t.Fatal("TestMarkMore IsMarker os.PathError fatal")
	}

	var rawPathErr *os.PathError
	if errors.As(markErr, &rawPathErr) {
		fmt.Printf("%+v\n", rawPathErr)
	} else {
		t.Fatal("TestMarkMore As os.PathError fatal")
	}
}

//测试对被包裹一次的错误进行多次标记的情况
func TestMarkErrorfMore(t *testing.T) {
	//测试哨兵类型的错误
	markErr := Mark(Mark(fmt.Errorf("more info: %w", os.ErrExist), 1), 2)
	if errors.Is(markErr, os.ErrExist) == false {
		t.Fatal("TestMarkErrorfMore Is os.ErrExist fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 2) == false {
		t.Fatal("TestMarkErrorfMore HasMaker os.ErrExist fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 2 {
		t.Fatal("TestMarkErrorfMore IsMarker os.ErrExist fatal")
	}

	//测试自定义类型的错误
	rawErr := errors.New("test raw error")
	markErr = Mark(Mark(fmt.Errorf("more info: %w", rawErr), 1), 1)
	if errors.Is(markErr, rawErr) == false {
		t.Fatal("TestMarkErrorfMore Is errors.New fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorfMore HasMaker errors.New fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorfMore IsMarker errors.New fatal")
	}

	//测试实现了 error 接口的错误
	pathErr := &os.PathError{Op: "op", Path: "Path", Err: errors.New("Err")}
	markErr = Mark(Mark(fmt.Errorf("more info: %w", pathErr), 1), 1)
	if errors.Is(markErr, pathErr) == false {
		t.Fatal("TestMarkErrorfMore Is os.PathError fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorfMore HasMaker os.PathError fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorfMore IsMarker os.PathError fatal")
	}

	var rawPathErr *os.PathError
	if errors.As(markErr, &rawPathErr) {
		fmt.Printf("%+v\n", rawPathErr)
	} else {
		t.Fatal("TestMarkErrorfMore As os.PathError fatal")
	}
}

//测试反复包裹与标记的情况
func TestMarkErrorfMarkErrorf(t *testing.T) {
	//测试哨兵类型的错误被包裹后的错误
	markErr := Mark(fmt.Errorf("mark: %w", Mark(fmt.Errorf("more info: %w", os.ErrExist), 1)), 2)
	if errors.Is(markErr, os.ErrExist) == false {
		t.Fatal("TestMarkErrorfMarkErrorf Is os.ErrExist fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 2) == false {
		t.Fatal("TestMarkErrorfMarkErrorf HasMaker os.ErrExist fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 2 {
		t.Fatal("TestMarkErrorfMarkErrorf IsMarker os.ErrExist fatal")
	}

	//测试自定义类型的错误被包裹后的错误
	rawErr := errors.New("test raw error")
	markErr = Mark(fmt.Errorf("mark: %w", Mark(fmt.Errorf("more info: %w", rawErr), 1)), 1)
	if errors.Is(markErr, rawErr) == false {
		t.Fatal("TestMarkErrorfMarkErrorf Is errors.New fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorfMarkErrorf HasMaker errors.New fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorfMarkErrorf IsMarker errors.New fatal")
	}

	//测试实现了 error 接口的错误被包裹后的错误
	pathErr := &os.PathError{Op: "op", Path: "Path", Err: errors.New("Err")}
	markErr = Mark(fmt.Errorf("mark: %w", Mark(fmt.Errorf("more info: %w", pathErr), 1)), 1)
	if errors.Is(markErr, pathErr) == false {
		t.Fatal("TestMarkErrorfMarkErrorf Is os.PathError fatal")
	} else {
		fmt.Printf("%+v\n", markErr)
	}
	if HasMaker(markErr, 1) == false {
		t.Fatal("TestMarkErrorfMarkErrorf HasMaker os.PathError fatal")
	}
	if m := IsMarker(markErr); m == nil || m.Code() != 1 {
		t.Fatal("TestMarkErrorfMarkErrorf IsMarker os.PathError fatal")
	}

	var rawPathErr *os.PathError
	if errors.As(markErr, &rawPathErr) {
		fmt.Printf("%+v\n", rawPathErr)
	} else {
		t.Fatal("TestMarkErrorfMarkErrorf As os.PathError fatal")
	}
}
