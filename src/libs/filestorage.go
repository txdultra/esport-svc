package libs

import (
	"bytes"
	"dbs"
	"errors"
	"fmt"
	"utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	//"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	//"net/textproto"
	"encoding/json"
	//"github.com/golang/groupcache"
	"image"
	//"logs"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	//gp "utils/groupcache"
)

type FileNode struct {
	FileId       int64  `json:"file_id"`
	Size         int64  `json:"size"`
	Path         string `json:"path"`
	ExtName      string `json:"ext"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime"`
	Width        int    `json:"w"`
	Height       int    `json:"h"`
	Source       int64  `json:"source"`
}

type IFileStorage interface {
	SaveFile(data []byte, fileName string, source int64) (*FileNode, error)
	DeleteFile(fileId int64) error
	GetFileUrl(fileId int64) string
	GetPrivateFileUrl(fileId int64) string
	GetFileNode(fileId int64) *FileNode
	GetFile(fileId int64) *File
	ClearFileUrlMap() bool
}

//const (
//	file_groupcache_name = "file-group"
//)

////注册groupcache
//func initToGroupCache() {
//	fs := NewWeedFsFileStorage()
//	gp.RegisterGroupCache(file_groupcache_name, func(_ groupcache.Context, key string, dest groupcache.Sink) error {
//		sps := strings.Split(key, ":")
//		if len(sps) < 2 {
//			return nil
//		}
//		fid, err := strconv.ParseInt(sps[1], 10, 64)
//		if err != nil {
//			logs.Errorf("convert key to fileid fail:%v", err)
//		}
//		fn := fs.GetFile(fid)
//		if fn != nil {
//			data, err := json.Marshal(fn)
//			if err == nil {
//				return dest.SetBytes(data)
//			}
//		}
//		return nil
//	})
//}

type FileStorage struct{}

func (f *FileStorage) FileCacheKey(id int64) string {
	return fmt.Sprintf("mobile_file_id:%d", id)
}

func NewFileStorage() IFileStorage {
	switch file_provider_name {
	case "weedfs":
		return NewWeedFsFileStorage()
	default:
		panic("undefined file storage")
	}
}

//weedfs distributed file system
type WeedFsFileStorage struct {
	FileStorage
	FileUrlMaps   map[int64]string
	VolumeUrlMaps map[string][]weedFsLocation
}

type weedFsAssign struct {
	Count     int    `json:"count"`
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Error     string `json:"error"`
}

type weedFsLocation struct {
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
}
type weedFsLookupResult struct {
	Locations   []weedFsLocation `json:"locations"`
	Error       string           `json:"error"`
	VolumeId    string           `json:"volumeId"`
	RefreshTime time.Time        `json:"refresh_time"`
}

var weedFsLock sync.Locker = &sync.Mutex{}
var weedFsStorage *WeedFsFileStorage
var weedFsScheme, weedFsMasterServer string
var weedFsCdnEnabled bool
var weedFsCdnMapTable map[string]string

func NewWeedFsFileStorage() *WeedFsFileStorage {
	if weedFsStorage == nil {
		weedFsLock.Lock()
		if weedFsStorage == nil {
			weedFsStorage = &WeedFsFileStorage{
				FileUrlMaps:   make(map[int64]string),
				VolumeUrlMaps: make(map[string][]weedFsLocation),
			}
			weedFsScheme = beego.AppConfig.String("weedfs.master.addr.scheme") + "://"
			weedFsMasterServer = beego.AppConfig.String("weedfs.master.addr")
			//file cdn optimiz configuration
			weedFsCdnEnabled, _ = beego.AppConfig.Bool("weedfs.cdn.enabled")
			weedFsCndConfig := beego.AppConfig.String("weedfs.cnd.config")
			wfcc := strings.Split(weedFsCndConfig, ",")
			weedFsCdnMapTable = make(map[string]string)
			for _, wfc := range wfcc {
				if len(wfc) == 0 {
					continue
				}
				wfs := strings.Split(wfc, ">")
				if len(wfs) != 2 {
					continue
				}
				if _, ok := weedFsCdnMapTable[wfs[0]]; ok {
					continue
				}
				weedFsCdnMapTable[wfs[0]] = wfs[1]
			}
		}
		weedFsLock.Unlock()
	}
	return weedFsStorage
}

func (f *WeedFsFileStorage) SaveFile(data []byte, fileName string, source int64) (*FileNode, error) {
	if len(fileName) == 0 {
		return nil, errors.New("文件名称不能为空")
	}
	defer func() { //异常处理
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	ext := strings.ToLower(path.Ext(fileName))
	mimeType, _ext := "", ""
	if len(ext) > 0 {
		_ext = ext[1:]
		//如果是图片，读取图片格式
		if utils.IsImage(_ext) {
			imgfmt := utils.ImageFormat(data)
			if len(imgfmt) > 0 {
				ext = "." + imgfmt
				_ext = imgfmt
			}
		}
		mimeType = mime.TypeByExtension(ext)
	}
	weedfs_master_addr := f.scheme() + f.fileMasterServer() + "/dir/assign?ts=" + strconv.Itoa(int(time.Now().Unix()))
	assign := &weedFsAssign{}
	_req := httplib.Get(weedfs_master_addr)
	_req.SetTimeout(file_connect_timeout, file_rwconnect_timeout)
	err := _req.ToJson(assign)
	if err != nil {
		return nil, errors.New("远程WeedFs分配错误:" + err.Error())
	}
	_temp := strings.Split(assign.Fid, ",")
	if len(_temp) != 2 {
		return nil, errors.New("远程WeedFs分配错误,Fid编码错误")
	}
	post_addr := f.scheme() + assign.Url + "/" + assign.Fid

	body_buf := &bytes.Buffer{}
	body_writer := multipart.NewWriter(body_buf)
	//关键的一步操作
	fileWriter, err := body_writer.CreateFormFile("uploadfile", fileName)
	fileWriter.Write(data)
	content_type := body_writer.FormDataContentType()
	if err = body_writer.Close(); err != nil {
		return nil, errors.New("初始化提交文件错误:" + err.Error())
	}

	rsp, err := http.Post(post_addr, content_type, body_buf)
	defer func() {
		if !rsp.Close {
			rsp.Body.Close()
		}
	}()

	if err != nil {
		return nil, errors.New("上传文件到文件服务器失败:" + err.Error())
	}
	rd, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.New("上传文件到文件服务器,返回错误:" + err.Error())
	}
	json_data, err := utils.NewJson(rd)
	if err != nil {
		return nil, errors.New("转换返回json数据时错误:" + err.Error())
	}
	size, _ := json_data.Get("size").Int64()
	if size < 0 {
		return nil, errors.New("文件未完成上传")
	}
	w, h := 0, 0
	//如果是图片,获取图片大小
	if utils.IsImage(_ext) {
		buf := bytes.NewBuffer(data)
		img, _, err := image.Decode(buf)
		if err == nil {
			bounds := img.Bounds()
			w = bounds.Max.X
			h = bounds.Max.Y
		}
	}

	volume := _temp[0]

	file := File{
		FileName:     assign.Fid,
		OriginalName: fileName,
		Volume:       volume,
		ExtName:      _ext,
		Size:         size,
		PostTime:     time.Now(),
		MimeType:     mimeType,
		Height:       h,
		Width:        w,
		Source:       source,
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(&file)
	if err != nil {
		return nil, errors.New("保持文件信息时错误:" + err.Error())
	}
	file.Id = id
	cache := utils.GetCache()
	cache.Set(f.FileCacheKey(id), file, 240*time.Hour)
	fileUrl := ""
	if len(_ext) > 0 {
		fileUrl = f.scheme() + assign.PublicUrl + "/" + assign.Fid + "." + file.ExtName
	} else {
		fileUrl = f.scheme() + assign.PublicUrl + "/" + assign.Fid
	}
	fn := &FileNode{
		FileId:       id,
		Size:         size,
		Path:         fileUrl,
		ExtName:      _ext,
		FileName:     assign.Fid,
		OriginalName: fileName,
		MimeType:     mimeType,
		Height:       h,
		Width:        w,
		Source:       source,
	}

	volUrl := assign.PublicUrl
	if weedFsCdnEnabled {
		_cdnUrl := f.getCdnUrl(volUrl)
		if len(_cdnUrl) > 0 {
			volUrl = strings.Replace(fileUrl, assign.PublicUrl, _cdnUrl, -1)
		}
	}
	//加入文件hash表
	f.FileUrlMaps[id] = volUrl
	return fn, nil
}

func (f *WeedFsFileStorage) GetFileNode(fileId int64) *FileNode {
	file := f.GetFile(fileId)
	if file == nil {
		return nil
	}
	fn := &FileNode{
		FileId:       file.Id,
		Size:         file.Size,
		Path:         f.GetFileUrl(file.Id),
		ExtName:      file.ExtName,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		MimeType:     file.MimeType,
		Width:        file.Width,
		Height:       file.Height,
		Source:       file.Source,
	}
	return fn
}

func (f *WeedFsFileStorage) GetFile(fileId int64) *File {
	var file File
	cache := utils.GetCache()
	err := cache.Get(f.FileCacheKey(fileId), &file)
	if err != nil {
		file.Id = fileId
		o := dbs.NewDefaultOrm()
		err := o.Read(&file)
		if err == orm.ErrNoRows || err == orm.ErrMissPK {
			return nil
		}
		cache.Set(f.FileCacheKey(fileId), file, 240*time.Hour)
	}
	return &file
}

func (f *WeedFsFileStorage) DeleteFile(fileId int64) error {
	file := &File{}
	file.Id = fileId
	o := dbs.NewDefaultOrm()
	err := o.Read(file)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		return err
	}
	volLocations, err := f.volumeUrl(file.Volume)
	if err != nil && len(volLocations) > 0 {
		return err
	}
	volUrl := volLocations[0].Url
	delUrl := f.scheme() + volUrl + "/" + file.FileName
	req, err := http.NewRequest("DELETE", delUrl, nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err == nil {
		file.IsDeleted = true
		o.Update(file)
		delete(f.FileUrlMaps, file.Id)
		cache := utils.GetCache()
		cache.Delete(f.FileCacheKey(fileId))
	}
	return err
}

func (f *WeedFsFileStorage) GetPrivateFileUrl(fileId int64) string {
	file := f.GetFile(fileId)
	if file == nil {
		return ""
	}
	volLocations, err := f.volumeUrl(file.Volume)
	if err != nil || len(volLocations) == 0 {
		return ""
	}
	volUrl := volLocations[0].Url
	fileUrl := ""
	if len(file.ExtName) > 0 {
		fileUrl = f.scheme() + volUrl + "/" + file.FileName + "." + file.ExtName
	} else {
		fileUrl = f.scheme() + volUrl + "/" + file.FileName
	}
	return fileUrl
}

func (f *WeedFsFileStorage) GetFileUrl(fileId int64) string {
	url, ok := f.FileUrlMaps[fileId]
	if ok {
		return url
	}

	file := f.GetFile(fileId)
	if file == nil {
		//将错误的也加入文件hash表
		f.FileUrlMaps[fileId] = ""
		return ""
	}
	volLocations, err := f.volumeUrl(file.Volume)
	if err != nil || len(volLocations) == 0 {
		return ""
	}
	volUrl := volLocations[0].PublicUrl
	if weedFsCdnEnabled {
		_cdnUrl := f.getCdnUrl(volUrl)
		if len(_cdnUrl) > 0 {
			volUrl = _cdnUrl
		}
	}
	fileUrl := ""
	if len(file.ExtName) > 0 {
		fileUrl = f.scheme() + volUrl + "/" + file.FileName + "." + file.ExtName
	} else {
		fileUrl = f.scheme() + volUrl + "/" + file.FileName
	}

	f.FileUrlMaps[fileId] = fileUrl
	return fileUrl
}

func (f *WeedFsFileStorage) ClearFileUrlMap() bool {
	f.VolumeUrlMaps = make(map[string][]weedFsLocation)
	f.FileUrlMaps = make(map[int64]string)
	return true
}

func (f *WeedFsFileStorage) volumeUrl(volumeId string) ([]weedFsLocation, error) {
	locs, ok := f.VolumeUrlMaps[volumeId]
	if ok {
		return locs, nil
	}
	weedfs_master_addr := f.fileMasterServer()
	masterUrl := f.scheme() + weedfs_master_addr + "/dir/lookup?volumeId=" + volumeId
	_req := httplib.Get(masterUrl)
	_req.SetTimeout(file_connect_timeout, file_rwconnect_timeout)
	bs, err := _req.Bytes()
	if err != nil {
		return nil, errors.New("文件master错误:" + err.Error())
	}
	var ret weedFsLookupResult
	err = json.Unmarshal(bs, &ret)
	if err != nil {
		return nil, errors.New("文件master返回数据错误:" + err.Error())
	}
	if len(ret.Locations) > 0 {
		f.VolumeUrlMaps[volumeId] = ret.Locations
	}
	return ret.Locations, nil
}

func (f *WeedFsFileStorage) scheme() string {
	return weedFsScheme
}

func (f *WeedFsFileStorage) fileMasterServer() string {
	return weedFsMasterServer
}

func (f *WeedFsFileStorage) getCdnUrl(url string) string {
	domain, ok := weedFsCdnMapTable[url]
	if ok {
		return domain
	}
	return ""
}
