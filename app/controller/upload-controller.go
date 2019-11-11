package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"hash"
	"im/common"
	"im/dao"
	"im/model"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type UploadController struct {
	Session *sessions.Session
	Ctx iris.Context
}

func (c *UploadController) Post(){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		file multipart.File
		info *multipart.FileHeader
		h hash.Hash
		attachment *model.Attachments

		mime string
		mimeSplit []string
		hashname string
		out *os.File
		timeStr string
		docPath string
		exist bool
		written int64
		fileLimit = "image,video,audio"

		err error
	)

	// Get the file from the request.
	file, info, err = c.Ctx.FormFile("file")
	if err != nil {
		goto ERR
	}
	defer file.Close()

	//sha1值计算
	h = sha1.New()
	written, err = io.Copy(h, file)
	fmt.Println("sha1计算时copy:",written)
	if err != nil {
		goto ERR
	}

	hashname=hex.EncodeToString(h.Sum(nil))
	attachment=dao.NewAttachmentsDao().GetFileBySha1(user.AppsId,hashname)
	if attachment!=nil{
		goto SUCCESS
	}

	//格式校验
	mime=info.Header.Get("Content-Type")
	mimeSplit=strings.Split(mime,"/")
	if len(mimeSplit)!=2{
		err=errors.New("文件类型未知")
		goto ERR
	}

	//文件类型校验
	if !strings.Contains(fileLimit,mimeSplit[0]){
		err=errors.New("不支持的文件类型")
		goto ERR
	}

	timeStr=time.Now().Format("20060102")
	hashname+="."+mimeSplit[1]
	docPath="./uploads/"+info.Header.Get("Content-Type")+"/"+timeStr
	fmt.Println(docPath)
	exist,err=common.PathExists(docPath)
	if !exist{
		err=os.MkdirAll(docPath,0755)
		if err!=nil{
			goto ERR
		}
	}

	//创建并复制
	out, err = os.OpenFile(docPath+"/"+hashname, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		goto ERR
	}
	defer out.Close()

	file.Seek(0,0)
	written,err=io.Copy(out, file)
	if err!=nil{
		goto ERR
	}

	//TODO 通过内网将资源上传至oss


	//TODO 插入附件记录

	goto SUCCESS


ERR:
	c.Ctx.StatusCode(iris.StatusInternalServerError)
	c.Ctx.JSON(common.SendCry("Error while uploading: <b>" + err.Error() + "</b>"))
	return

SUCCESS:
	c.Ctx.JSON(common.SendSmile(1))
	return
}