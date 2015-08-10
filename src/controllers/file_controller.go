//error prefix "F5"
package controllers

import (
	"bytes"
	"libs"
)

// 文件 API
type FileController struct {
	BaseController
}

func (f *FileController) Prepare() {
	f.BaseController.Prepare()
}

func (f *FileController) URLMapping() {
	f.Mapping("D", f.D)
	f.Mapping("Get", f.Get)
	f.Mapping("GetUrl", f.GetUrl)
	f.Mapping("Delete", f.Delete)
	f.Mapping("Upload", f.Upload)
}

// @Title 文件下载
// @Description 文件下载
// @Param   id     path    int  true        "文件id"
// @Success 200
// @router /d/:id([0-9]+) [get]
func (f *FileController) D() {
	id, err := f.GetInt64(":id")
	if err != nil || id <= 0 {
		f.Abort("404")
		return
	}
	fileUrl := file_storage.GetFileUrl(id)
	if len(fileUrl) == 0 {
		f.Abort("404")
		return
	}
	f.Redirect(fileUrl, 301)
}

// @Title 获取文件信息
// @Description 提交方式为GET
// @Param   id     path    int  true        "文件id"
// @Success 200 {object} libs.FileNode
// @router /:id([0-9]+) [get]
func (f *FileController) Get() {
	id, err := f.GetInt64(":id")
	if err != nil {
		f.Json(libs.NewError("file_parameter", "F5001", "input parameter error ", ""))
		return
	}
	file := file_storage.GetFileNode(id)
	if file == nil {
		f.Json(libs.NewError("file_not_exist", "F5002", "file not exist ", ""))
		return
	}
	f.Json(file)
}

// @Title 获取文件URI
// @Description 获取文件URI
// @Param   id     path    int  true        "文件id"
// @Success 200
// @router /url/:id([0-9]+) [get]
func (f *FileController) GetUrl() {
	id, _ := f.GetInt64(":id")
	if id <= 0 {
		f.WriteString("")
		return
	}
	uri := file_storage.GetFileUrl(id)
	f.WriteString(uri)
}

// @Title 删除文件
// @Description 提交方式为DELETE,成功删除error_code=F5000
// @Param id path int true "文件id"
// @Success 200 {object} libs.Error
// @router /:id([0-9]+) [delete]
func (f *FileController) Delete() {
	id, err := f.GetInt64(":id")
	if err != nil {
		f.Json(libs.NewError("file_parameter", "F5001", "input parameter error ", ""))
		return
	}
	err = file_storage.DeleteFile(id)
	if err != nil {
		f.Json(libs.NewError("file_delete_fail", "F5003", err.Error(), ""))
		return
	}
	f.Json(libs.NewError("file_delete_success", "F5000", "file deleted", ""))
}

// @Title 上传文件
// @Description 提交方式为POST
// @Param file form string true "Form文件"
// @Success 200 {object} libs.FileNode
// @router /upload [post]
func (f *FileController) Upload() {
	file, fheader, err := f.GetFile("file")
	source, _ := f.GetInt64("source")
	if err != nil {
		f.Json(libs.NewError("file_upload_fail", "F5004", err.Error(), ""))
		return
	}
	file_name := fheader.Filename
	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(file)
	if err != nil {
		f.Json(libs.NewError("file_upload_readData_fail", "F5005", err.Error(), ""))
		return
	}
	node, err := file_storage.SaveFile(buf.Bytes(), file_name, source)
	if err != nil {
		f.Json(libs.NewError("file_upload_saveFile_fail", "F5006", err.Error(), ""))
		return
	}
	f.Json(node)
}
