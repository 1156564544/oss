package objects

import (
	"ApiServer/locate"
	"crypto/sha256"
	"io"
	"encoding/base64"
	"net/url"

	"rs"
)

// 计算hash值
func calculateHash(r io.Reader)string {
	h := sha256.New()
	io.Copy(h, r)
	hash:= base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.PathEscape(hash)
}

// 判断object是否存在
func Exist(hash string)bool{
	return len(locate.Locate(hash)) >= rs.NUM_DATA_SHARES
}

// 计算单个chunk的大小
func getPersharedSize(size int64)int64{
	return (size + int64(rs.NUM_DATA_SHARES-1)) / int64(rs.NUM_DATA_SHARES)
}

// 计算一个round读写的数据量
func getRoundSize()int{
	return rs.CHUNK_SIZE*(rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES)
}