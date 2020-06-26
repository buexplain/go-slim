package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	//超出大小
	ErrExceedSize = errors.New("upload file size exceeds the maximum value")
	//禁止的后缀
	ErrDenyExt = errors.New("extensions to upload is not allowed")
	//禁止的 mime type
	ErrDenyContentType = errors.New("Content-Type to upload is not allowed")
	//坏的图片
	ErrInvalidImage = errors.New("invalid image")
)

type Upload struct {
	file   multipart.File
	header *multipart.FileHeader
	//文件名 含扩展
	name string
	//文件扩展
	ext string
	//文件 Content-Type
	contentType string
	//文件md5值
	md5 string
	//创建目录的权限
	dirPerm os.FileMode
	//创建文件的权限
	filePerm os.FileMode
	//保存路径的规则
	pathRule PathRule
	//文件名生成规则
	nameRule NameRule
	//允许的文件大小
	maxSize int64
	//允许的文件扩展
	allowExt []string
	//允许的文件 Content-Type
	allowContentType []string
	//是否重绘图片
	isRedraw bool
	//重绘图片的参数
	redrawOptions interface{}
	//是否校验
	isValidate bool
	//校验结果
	validateErr error
	//最后的存储结果
	result string
}

func New(file multipart.File, header *multipart.FileHeader) *Upload {
	tmp := &Upload{
		file:   file,
		header: header,
	}
	//获得上传文件的名称
	_ = tmp.Name()
	//默认不进行目录创建
	tmp.SetPathRule(PathRuleOriginal)
	//默认原始文件名存储
	tmp.SetNameRule(NameRuleOriginal)
	//默认文件权限
	tmp.SetFilePerm(0666)
	//默认目录的权限
	tmp.SetDirPerm(0777)
	return tmp
}

func (this *Upload) MultipartFile() multipart.File {
	return this.file
}

func (this *Upload) MultipartHeader() *multipart.FileHeader {
	return this.header
}

//关闭打开的文件
func (this *Upload) Close() error {
	return this.file.Close()
}

//返回文件名称
func (this *Upload) Name() string {
	if this.name == "" {
		//得到文件原始名称
		name := strings.Trim(filepath.Base(strings.Trim(this.header.Filename, " ")), " ")
		if name == "." {
			name = ""
		}
		//得到文件扩展
		ext := filepath.Ext(name)
		if ext == "." {
			ext = ""
		} else {
			//去掉扩展中的点
			ext = strings.TrimLeft(ext, ".")
			//强制文件扩展为小写
			ext = strings.ToLower(ext)
			if name[len(name)-len(ext):] != ext {
				name = name[:len(name)-len(ext)] + ext
			}
		}
		this.name = name
		this.ext = ext
	}
	return this.name
}

//返回文件后缀
func (this *Upload) Ext() string {
	_ = this.Name()
	return this.ext
}

//返回文件大小
func (this *Upload) Size() int64 {
	return this.header.Size
}

//返回上传的结果，这个得在SaveToPath方法结束后，没有错误的情况下调用
func (this *Upload) Result() string {
	return this.result
}

//返回文件头内容
func (this *Upload) ContentType() (string, error) {
	if this.contentType == "" {
		if _, err := this.file.Seek(0, 0); err != nil {
			return "", err
		}
		var b []byte = make([]byte, 512)
		if _, err := this.file.Read(b); err != nil {
			return "", err
		}
		this.contentType = http.DetectContentType(b)
	}
	return this.contentType, nil
}

//返回文件md5值
func (this *Upload) MD5() (string, error) {
	if this.md5 == "" {
		hash := md5.New()
		if _, err := this.file.Seek(0, 0); err != nil {
			return "", err
		}
		if _, err := io.Copy(hash, this.file); err != nil {
			return "", err
		}
		this.md5 = hex.EncodeToString(hash.Sum(nil))
	}
	return this.md5, nil
}

//设置创建的文件的权限
func (this *Upload) SetFilePerm(perm os.FileMode) *Upload {
	this.filePerm = perm
	return this
}

//设置创建的文件的权限
func (this *Upload) SetDirPerm(perm os.FileMode) *Upload {
	this.dirPerm = perm
	return this
}

//设置存储路径生成规则
func (this *Upload) SetPathRule(pathRule PathRule) *Upload {
	this.pathRule = pathRule
	return this
}

//设置文件名生成规则
func (this *Upload) SetNameRule(nameRule NameRule) *Upload {
	this.nameRule = nameRule
	return this
}

//校验文件大小
func (this *Upload) SetValidateSize(maxSize int64) *Upload {
	this.maxSize = maxSize
	return this
}

//校验文件后缀
func (this *Upload) SetValidateExt(ext ...string) *Upload {
	if len(ext) > 0 {
		if this.allowExt == nil {
			this.allowExt = make([]string, 0, len(ext))
		}
		this.allowExt = append(this.allowExt, ext...)
	}
	return this
}

//校验文件头
func (this *Upload) SetValidateContentType(contentType ...string) *Upload {
	if len(contentType) > 0 {
		if this.allowContentType == nil {
			this.allowContentType = make([]string, 0, len(contentType))
		}
		this.allowContentType = append(this.allowContentType, contentType...)
	}
	return this
}

//重绘图片
func (this *Upload) SetRedraw(isRedraw bool, options interface{}) *Upload {
	this.isRedraw = isRedraw
	this.redrawOptions = options
	return this
}

func (this *Upload) Validate() error {
	//判断是否已经校验
	if this.isValidate {
		return this.validateErr
	}

	//设置已经校验了
	this.isValidate = true

	//初始化校验结果
	this.validateErr = nil

	//验证大小
	if this.validateErr == nil && this.maxSize > 0 && this.header.Size > this.maxSize {
		this.validateErr = ErrExceedSize
	}

	//验证后缀
	if this.validateErr == nil && this.allowExt != nil {
		var isError bool = true
		for _, v := range this.allowExt {
			if strings.EqualFold(this.ext, v) {
				isError = false
				break
			}
		}
		if isError {
			this.validateErr = ErrDenyExt
		}
	}

	//验证 Content-Type
	if this.validateErr == nil && this.allowContentType != nil {
		if contentType, err := this.ContentType(); err != nil {
			this.validateErr = err
		} else {
			var isError bool = true
			for _, v := range this.allowContentType {
				if strings.EqualFold(contentType, v) {
					isError = false
					break
				}
			}
			if isError {
				this.validateErr = ErrDenyContentType
			}
		}
	}

	return this.validateErr
}

func (this *Upload) Save(w io.Writer) (int64, error) {
	//校验文件
	if err := this.Validate(); err != nil {
		return 0, err
	}

	//拨回偏移量
	if _, err := this.file.Seek(0, 0); err != nil {
		return 0, err
	}

	//重绘图片
	if this.isRedraw {
		var img image.Image
		var imgGIF *gif.GIF
		var err error
		switch this.ext {
		case "jpg", "jpeg":
			if img, err = jpeg.Decode(this.file); err != nil {
				return 0, ErrInvalidImage
			} else {
				if o, ok := this.redrawOptions.(*jpeg.Options); ok {
					err = jpeg.Encode(w, img, o)
				} else {
					err = jpeg.Encode(w, img, o)
				}
			}
		case "png":
			if img, err = png.Decode(this.file); err != nil {
				return 0, ErrInvalidImage
			} else {
				err = png.Encode(w, img)
			}
		case "gif":
			if imgGIF, err = gif.DecodeAll(this.file); err != nil {
				return 0, ErrInvalidImage
			} else {
				err = gif.EncodeAll(w, imgGIF)
			}
		}
		if img != nil || imgGIF != nil {
			if err != nil {
				return 0, ErrInvalidImage
			}
			//如果图片中包含非图片所需字节，如恶意代码，或者是.jpg图片，那么此时返回的写入字节是有误差的
			return this.header.Size, nil
		}
	}

	//写入文件
	return io.Copy(w, this.file)
}

func (this *Upload) SaveToFile(file string) error {
	//校验文件
	if err := this.Validate(); err != nil {
		return err
	}
	//创建文件
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, this.filePerm)
	if err != nil {
		return err
	}

	//写入文件
	if _, saveErr := this.Save(f); saveErr == nil {
		//写入成功
		if closeErr := f.Close(); closeErr == nil {
			//关闭成功
			return nil
		} else { //关闭失败
			//移除文件
			if removeErr := os.Remove(file); removeErr == nil {
				//删除成功
				return closeErr
			} else {
				//删除失败
				return fmt.Errorf("closeErr: %s removeErr: %s", closeErr, removeErr)
			}
		}
	} else { //写入失败
		errs := strings.Builder{}
		errs.WriteString("saveErr: ")
		errs.WriteString(saveErr.Error())

		//关闭文件
		if closeErr := f.Close(); closeErr != nil {
			errs.WriteString(" closeErr: ")
			errs.WriteString(closeErr.Error())
		}

		//移除文件
		if removeErr := os.Remove(file); removeErr != nil {
			errs.WriteString(" removeErr: ")
			errs.WriteString(removeErr.Error())
		}

		return errors.New(errs.String())
	}
}

func (this *Upload) SaveToPath(path string) (string, error) {
	//校验文件
	if err := this.Validate(); err != nil {
		return "", err
	}

	//生成存储路径
	if tmp, err := CreatePath(this.pathRule); err != nil {
		return "", err
	} else {
		if tmp != "" {
			path = filepath.Join(path, tmp)
		}
	}

	//创建路径
	if err := os.MkdirAll(path, this.dirPerm); err != nil {
		return "", err
	}

	//创建文件名
	var file string
	if name := CreateName(this.nameRule); name == "" {
		file = filepath.Join(path, this.name)
	} else {
		file = filepath.Join(path, name+"."+this.ext)
	}

	//写入文件
	if err := this.SaveToFile(file); err != nil {
		return "", err
	} else {
		file = filepath.ToSlash(file)
		//持有当前的文件路径
		this.result = file
		return file, err
	}
}
