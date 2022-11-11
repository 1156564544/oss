package objects

import (
	"ApiServer/locate"
	"crypto/sha256"
	"io"
	"encoding/base64"

	"rs"
)

// 计算hash值
func calculateHash(r io.Reader)string {
	h := sha256.New()
	io.Copy(h, r)
	hash:= base64.StdEncoding.EncodeToString(h.Sum(nil))
	// 这里是因为我算出来的hash值最后面有一个等号hao，所以我把它去掉了，客户端提交的hash值也要做同样的处理
	// return hash[:len(hash)-1]
	return hash
}

// 判断object是否存在
func Exist(hash string)bool{
	return len(locate.Locate(hash)) >= rs.NUM_DATA_SHARES
}

// 计算单个chunk的大小
func getPersharedSize(size int64)int64{
	return (size + int64(rs.NUM_DATA_SHARES-1)) / int64(rs.NUM_DATA_SHARES)
}