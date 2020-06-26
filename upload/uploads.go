package upload

import (
	"errors"
	"io"
	"os"
	"strings"
)

type Uploads []*Upload

func (this Uploads) Close() error {
	errs := strings.Builder{}
	for _, file := range this {
		if err := file.Close(); err != nil {
			if errs.Len() > 0 {
				errs.WriteByte('\n')
			}
			errs.WriteString(err.Error())
		}
	}
	if errs.Len() > 0 {
		return errors.New(errs.String())
	}
	return nil
}

//返回文件名称
func (this Uploads) Names() []string {
	result := make([]string, len(this), len(this))
	for k, upload := range this {
		result[k] = upload.Name()
	}
	return result
}

//返回文件后缀
func (this Uploads) Exts() []string {
	result := make([]string, len(this), len(this))
	for k, upload := range this {
		result[k] = upload.Ext()
	}
	return result
}

//返回文件大小
func (this Uploads) Sizes() []int64 {
	result := make([]int64, len(this), len(this))
	for k, upload := range this {
		result[k] = upload.Size()
	}
	return result
}

func (this Uploads) SetFilePerm(perm os.FileMode) Uploads {
	for _, upload := range this {
		upload.SetFilePerm(perm)
	}
	return this
}

func (this Uploads) SetDirPerm(perm os.FileMode) Uploads {
	for _, upload := range this {
		upload.SetDirPerm(perm)
	}
	return this
}

func (this Uploads) SetPathRule(pathRule PathRule) Uploads {
	for _, upload := range this {
		upload.SetPathRule(pathRule)
	}
	return this
}

func (this Uploads) SetNameRule(nameRule NameRule) Uploads {
	for _, upload := range this {
		upload.SetNameRule(nameRule)
	}
	return this
}

func (this Uploads) SetValidateSize(maxSize int64) Uploads {
	for _, upload := range this {
		upload.SetValidateSize(maxSize)
	}
	return this
}

func (this Uploads) SetValidateExt(ext ...string) Uploads {
	for _, upload := range this {
		upload.SetValidateExt(ext...)
	}
	return this
}

func (this Uploads) SetValidateContentType(contentType ...string) Uploads {
	for _, upload := range this {
		upload.SetValidateContentType(contentType...)
	}
	return this
}

func (this Uploads) SetRedraw(isRedraw bool, options interface{}) Uploads {
	for _, upload := range this {
		upload.SetRedraw(isRedraw, options)
	}
	return this
}

func (this Uploads) Validate() error {
	for _, upload := range this {
		if err := upload.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (this Uploads) Save(w io.Writer) (int64, error) {
	var written int64 = 0
	for _, upload := range this {
		n, err := upload.Save(w)
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func (this Uploads) SaveToFile(file string) error {
	errs := strings.Builder{}
	for _, upload := range this {
		if err := upload.SaveToFile(file); err != nil {
			errs.WriteString(err.Error())
			break
		}
	}
	return nil
}

func (this Uploads) SaveToPath(path string) ([]string, error) {
	//写入成功的文件集合
	files := make([]string, 0, len(this))
	//错误集合
	errs := strings.Builder{}
	//开始写入文件
	for _, upload := range this {
		if file, err := upload.SaveToPath(path); err != nil {
			//写入失败，跳出
			errs.WriteString(err.Error())
			break
		} else {
			//写入成功，记录
			files = append(files, file)
		}
	}
	//判断是否存在写入错误
	if errs.Len() > 0 {
		//存在写入错误，开始删除已经写入成功的文件
		for _, file := range files {
			//删除写入成功的文件
			if err := os.Remove(file); err != nil {
				//删除失败记录错误
				errs.WriteByte('\n')
				errs.WriteString(err.Error())
			}
		}
		return nil, errors.New(errs.String())
	}
	return files, nil
}
