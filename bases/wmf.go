package bases

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type CT_WmfCfg struct {
	Uri string
}

func New_CT_WmfCfg() *CT_WmfCfg {
	return &CT_WmfCfg{}
}

var WmfConfiguration *CT_WmfCfg

func checkWmfCfg() {
	if WmfConfiguration == nil {
		log.Panic("没有实例化office-parser的wmf配置，无法进行wmf图片转换，请检查`OfficeParserQiniuCfg`变量")
	}
}

func (w *CT_WmfCfg) convert(filepath string, imageName string) {
	checkWmfCfg()

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	_, err := bodyWriter.CreateFormFile("file", fmt.Sprintf("%s.wmf", imageName))
	if err != nil {
		log.Panicf("创建失败,err=%s", err)
	}

	fh, err := os.Open(filepath)
	if err != nil {
		log.Panicf("文件打开失败,err=%s", err)
	}

	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	reqReader := io.MultiReader(bodyBuf, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		log.Panicf("Error Stating file: %s", filepath)
	}

	req, err := http.NewRequest("POST", w.Uri, reqReader)
	if err != nil {
		log.Panicf("文件传输,err=%s", err)
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("地址请求失败,err=%s", err)
	}
}

func WmfConvert(filepath string, imageName string) {
	WmfConfiguration.convert(filepath, imageName)
}
