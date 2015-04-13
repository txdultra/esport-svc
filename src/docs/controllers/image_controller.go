package controllers

import (
	"bytes"
	"fmt"
	"image"
	"libs"
	"strings"
	"time"
	"utils"

	"github.com/astaxie/beego/httplib"
	"github.com/disintegration/imaging"
)

// 图片 API
type ImageController struct {
	BaseController
}

func (f *ImageController) Prepare() {
}

func (f *ImageController) URLMapping() {
	f.Mapping("Resize", f.Resize)
	f.Mapping("Crop", f.Crop)
}

// @Title 图片修改大小
// @Description 图片大小(高度宽度不能同时为0);Fit,Thumbnail模式高宽不能为0
// @Param   fid    path    int  true        "文件id"
// @Param   w     path    int  true        "宽度(可为0)"
// @Param   h     path    int  true        "高度(可为0)"
// @Param   r     path    string  false    "模式(Resize=R,Fit=F,Thumbnail=T),默认Resize"
// @Param   ro    path    string  false    "模式(90度=90,180度=180,270度=270,FlipH=fh,FlipV=fv),默认0"
// @Param   m     path    string  false        "模式(暂时忽略)"
// @Success 200  {object} libs.FileNode
// @router /resize [get]
func (f *ImageController) Resize() {
	fileId, _ := f.GetInt64("fid")
	w, _ := f.GetInt("w")
	h, _ := f.GetInt("h")
	r := f.GetString("r")
	ro := f.GetString("ro")
	if fileId <= 0 {
		f.Json(libs.NewError("img_resize_fileid", "F6001", "fileid don't less equal zero", ""))
		return
	}
	if (w == 0 && h == 0) || (w < 0 || h < 0) {
		f.Json(libs.NewError("img_resize_w_h", "F6002", "width and height not equal to zero at the same", ""))
		return
	}
	if (w == 0 || h == 0) && (strings.ToLower(r) == "f" || strings.ToLower(r) == "t") {
		f.Json(libs.NewError("img_resize_w_h", "F6002", "Fit,Thumbnail mode width and height not equal to zero at the same", ""))
		return
	}
	file := file_storage.GetFile(fileId)
	if file == nil {
		f.Json(libs.NewError("img_resize_file_notexist", "F6003", "file not exist", ""))
		return
	}
	if file.IsDeleted {
		f.Json(libs.NewError("img_resize_file_deleted", "F6004", "file deleted", ""))
		return
	}
	if file.ExtName == "gif" {
		f.Json(libs.NewError("img_resize_gif_notsupport", "F6005", "gif photo not support resize", ""))
		return
	}
	if !utils.IsImage(file.ExtName) {
		f.Json(libs.NewError("img_resize_file_isnot_image", "F6005", "file isn't image file", ""))
		return
	}
	fileUrl := file_storage.GetFileUrl(file.Id)
	request := httplib.Get(fileUrl)
	request.SetTimeout(1*time.Minute, 1*time.Minute)
	data, err := request.Bytes()
	if err != nil {
		f.Json(libs.NewError("img_resize_fileread_fail", "F6006", err.Error(), ""))
		return
	}
	buf := bytes.NewBuffer(data)
	srcImg, _, err := image.Decode(buf)
	if err != nil {
		f.Json(libs.NewError("img_resize_img_decode_fail", "F6009", err.Error(), ""))
		return
	}
	var dstImage image.Image
	switch strings.ToLower(r) {
	case "f":
		dstImage = imaging.Fit(srcImg, int(w), int(h), imaging.Lanczos)
	case "t":
		dstImage = imaging.Thumbnail(srcImg, int(w), int(h), imaging.Lanczos)
	default:
		dstImage = imaging.Resize(srcImg, int(w), int(h), imaging.Lanczos)
	}

	switch ro {
	case "90":
		dstImage = imaging.Rotate90(dstImage)
	case "180":
		dstImage = imaging.Rotate180(dstImage)
	case "270":
		dstImage = imaging.Rotate270(dstImage)
	case "fh":
		dstImage = imaging.FlipH(dstImage)
	case "fv":
		dstImage = imaging.FlipV(dstImage)
	default:
		break
	}

	fileData, err := utils.ImageToBytes(dstImage, file.OriginalName)
	if err != nil {
		f.Json(libs.NewError("img_resize_convert_fail", "F6007", err.Error(), ""))
		return
	}
	resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), w, h, file.ExtName)
	node, err := file_storage.SaveFile(fileData, resizeName, file.Id)
	if err != nil {
		f.Json(libs.NewError("img_resize_saveFile_fail", "F6008", err.Error(), ""))
		return
	}
	f.Json(node)
}

// @Title 图片截取
// @Description normal正常截取,center从中心窃取(只需传x1,y1 对应宽高)
// @Param   fid    path    int  true    "文件id"
// @Param   x1     path    int  true    "x1轴"
// @Param   y1     path    int  true    "y1轴"
// @Param   x2     path    int  true    "x2轴"
// @Param   y2     path    int  true    "y2轴"
// @Param   m     path    string  false "模式:normal,center,默认normal"
// @Success 200  {object} libs.FileNode
// @router /crop [get]
func (f *ImageController) Crop() {
	fileId, _ := f.GetInt64("fid")
	x1, _ := f.GetInt("x1")
	y1, _ := f.GetInt("y1")
	x2, _ := f.GetInt("x2")
	y2, _ := f.GetInt("y2")
	m := f.GetString("m")
	if fileId <= 0 {
		f.Json(libs.NewError("img_crop_fileid", "F6011", "fileid don't less equal zero", ""))
		return
	}
	if (x1 <= 0 || y1 <= 0 || x2 <= 0 || y2 <= 0) && (strings.ToLower(m) == "center" && x1 <= 0 && y1 <= 0) {
		f.Json(libs.NewError("img_crop_x_y", "F6012", "x,y don't less equal zero", ""))
		return
	}
	file := file_storage.GetFile(fileId)
	if file == nil {
		f.Json(libs.NewError("img_crop_file_notexist", "F6013", "file not exist", ""))
		return
	}
	if file.IsDeleted {
		f.Json(libs.NewError("img_crop_file_deleted", "F6014", "file deleted", ""))
		return
	}
	if file.ExtName == "gif" {
		f.Json(libs.NewError("img_crop_gif_notsupport", "F6015", "gif photo not support resize", ""))
		return
	}
	if !utils.IsImage(file.ExtName) {
		f.Json(libs.NewError("img_crop_file_isnot_image", "F6015", "file isn't image file", ""))
		return
	}
	fileUrl := file_storage.GetFileUrl(file.Id)
	request := httplib.Get(fileUrl)
	request.SetTimeout(1*time.Minute, 1*time.Minute)
	data, err := request.Bytes()
	if err != nil {
		f.Json(libs.NewError("img_crop_fileread_fail", "F6016", err.Error(), ""))
		return
	}

	buf := bytes.NewBuffer(data)
	srcImg, _, err := image.Decode(buf)
	if err != nil {
		f.Json(libs.NewError("img_crop_img_decode_fail", "F6019", err.Error(), ""))
		return
	}
	var dstImage image.Image
	switch strings.ToLower(m) {
	case "center":
		dstImage = imaging.CropCenter(srcImg, int(x1), int(y2))
	default:
		fmt.Println(x1, y1, x2, y2)
		dstImage = imaging.Crop(srcImg, image.Rect(int(x1), int(y1), int(x2), int(y2)))
	}

	fileData, err := utils.ImageToBytes(dstImage, file.OriginalName)
	if err != nil {
		f.Json(libs.NewError("img_crop_convert_fail", "F6017", err.Error(), ""))
		return
	}
	resizeName := fmt.Sprintf("%s-c%dx%dx%dx%d.%s", utils.FileName(file.OriginalName), x1, y1, x2, y2, file.ExtName)
	node, err := file_storage.SaveFile(fileData, resizeName, file.Id)
	if err != nil {
		f.Json(libs.NewError("img_crop_saveFile_fail", "F6018", err.Error(), ""))
		return
	}
	f.Json(node)
}
