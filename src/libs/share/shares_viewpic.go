package share

import (
	"bytes"
	"dbs"
	"fmt"
	"image"
	"libs"
	"logs"
	"sync"
	"time"
	"utils"

	"github.com/astaxie/beego/httplib"
	"github.com/disintegration/imaging"
)

var file_storage = libs.NewWeedFsFileStorage()
var picTbls map[string]bool = make(map[string]bool)
var pic_locker *sync.Mutex = new(sync.Mutex)

const (
	share_view_picture_fid  = "share_view_picture_fid:%d"
	share_pic_size_pfxtable = "share_view_pic"
)

type ShareViewPics struct{}

func NewShareViewPics() *ShareViewPics {
	return &ShareViewPics{}
}

//库分表
func (s *ShareViewPics) hash_pic_tbl(fileId int64) string {
	tbl := fmt.Sprintf("%s_%d", share_pic_size_pfxtable, fileId%50)
	if _, ok := picTbls[tbl]; ok {
		return tbl
	}
	pic_locker.Lock()
	defer pic_locker.Unlock()
	if _, ok := picTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(share_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(pfid int(11) NOT NULL,
		  fid int(11) NOT NULL,
		  ts bigint(15) NOT NULL,
		  ps smallint(6) NOT NULL,
		  PRIMARY KEY (pfid,ps),
  		  KEY idx_pfid (pfid)
		  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		picTbls[tbl] = true
	}
	return tbl
}

func (s *ShareViewPics) Create(file *libs.File, picSizes []libs.PIC_SIZE) map[libs.PIC_SIZE]ShareViewPicture {
	sps := make(map[libs.PIC_SIZE]ShareViewPicture)
	data, err := s.fileData(file.Id)
	if err != nil {
		fmt.Println("---------------------------------", err)
		logs.Errorf("share view pic err:%s", err.Error())
		return sps
	}
	csps := []ShareViewPicture{}
	tbl := s.hash_pic_tbl(file.Id)
	o := dbs.NewOrm(share_db)
	for _, spsize := range picSizes {
		size := 0
		switch spsize {
		case libs.PIC_SIZE_MIDDLE:
			size = sns_share_pic_middle_w
			break
		case libs.PIC_SIZE_THUMBNAIL:
			size = sns_share_pic_thumbnail_w
			break
		default:
			size = 0
			break
		}
		sp := ShareViewPicture{}
		sp.ParentFileId = file.Id
		sp.Ts = time.Now().Unix()
		sp.PicSize = spsize
		if spsize != libs.PIC_SIZE_ORIGINAL {
			fid, err := s.resize(data, file, size)
			if err != nil {
				fmt.Println("---------------------------------", err)
				continue
			}
			sp.FileId = fid
		} else {
			sp.FileId = file.Id
		}
		sql := fmt.Sprintf("insert into %s(pfid,fid,ts,ps) values(?,?,?,?)", tbl)
		o.Raw(sql, sp.ParentFileId, sp.FileId, sp.Ts, sp.PicSize).Exec()
		sps[spsize] = sp
		csps = append(csps, sp)
	}
	cache := utils.GetCache()
	cache.Set(fmt.Sprintf(share_view_picture_fid, file.Id), csps, 92*time.Hour)
	return sps
}

func (s *ShareViewPics) Get(fid int64) map[libs.PIC_SIZE]ShareViewPicture {
	cache := utils.GetCache()
	key := fmt.Sprintf(share_view_picture_fid, fid)
	tbl := s.hash_pic_tbl(fid)
	var csps []ShareViewPicture
	toMap := func(arr []ShareViewPicture) map[libs.PIC_SIZE]ShareViewPicture {
		_sps := make(map[libs.PIC_SIZE]ShareViewPicture)
		for _, ar := range arr {
			_sps[ar.PicSize] = ar
		}
		return _sps
	}
	err := cache.Get(key, &csps)
	if err != nil {
		//logs.Errorf("share pic get pictures cache fail:%s", err.Error())
	} else {
		return toMap(csps)
	}
	o := dbs.NewOrm(share_db)
	_, err = o.Raw("select pfid,fid,ts,ps from "+tbl+" where pfid=?", fid).QueryRows(&csps)
	if err != nil {
		return toMap(csps)
	}
	cache.Set(key, csps, 92*time.Hour)
	return toMap(csps)
}

func (s *ShareViewPics) fileData(fid int64) ([]byte, error) {
	fileUrl := file_storage.GetPrivateFileUrl(fid)
	request := httplib.Get(fileUrl)
	request.SetTimeout(1*time.Minute, 1*time.Minute)
	data, err := request.Bytes()
	if err != nil {
		logs.Errorf("share pic get file data fail:%s", err.Error())
		return nil, err
	}
	return data, nil
}

func (s *ShareViewPics) resize(data []byte, file *libs.File, size int) (int64, error) {
	if size <= 0 {
		return file.Id, nil
	}
	if file.Width <= size && file.Height <= size {
		return file.Id, nil
	}
	if file.Width >= file.Height && file.Width > size {
		buf := bytes.NewBuffer(data)
		srcImg, _, err := image.Decode(buf)
		if err != nil {
			logs.Errorf("share pic thumbnail by width's ratio data convert to image fail:%s", err.Error())
			return 0, err
		}
		var dstImage image.Image
		dstImage = imaging.Resize(srcImg, size, 0, imaging.Lanczos)
		fileData, err := utils.ImageToBytes(dstImage, file.OriginalName)
		if err != nil {
			logs.Errorf("share pic thumbnail by width's ratio image to bytes fail:%s", err.Error())
			return 0, err
		}
		ratio := float32(file.Width) / float32(size)
		resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), size, int(float32(file.Height)*ratio), file.ExtName)
		node, err := file_storage.SaveFile(fileData, resizeName, file.Id)
		if err != nil {
			logs.Errorf("share pic thumbnail by width's ratio save image fail:%s", err.Error())
			return 0, err
		}
		return node.FileId, err
	}
	if file.Width < file.Height && file.Height > size {
		buf := bytes.NewBuffer(data)
		srcImg, _, err := image.Decode(buf)
		if err != nil {
			logs.Errorf("share pic thumbnail by height's ratio data convert to image fail:%s", err.Error())
			return 0, err
		}
		var dstImage image.Image
		dstImage = imaging.Resize(srcImg, 0, size, imaging.Lanczos)
		fileData, err := utils.ImageToBytes(dstImage, file.OriginalName)
		if err != nil {
			logs.Errorf("share pic thumbnail by height's ratio image to bytes fail:%s", err.Error())
			return 0, err
		}
		ratio := float32(file.Height) / float32(size)
		resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), int(float32(file.Width)*ratio), size, file.ExtName)
		node, err := file_storage.SaveFile(fileData, resizeName, file.Id)
		if err != nil {
			logs.Errorf("share pic thumbnail by height's ratio save image fail:%s", err.Error())
			return 0, err
		}
		return node.FileId, err
	}
	return file.Id, nil
}
