package temp

import (
	"log"
	"net/http"
	"rs"
	"strings"
	"io"

	"ApiServer/objects"
	"httpTool"
	"es"
)

func put(w http.ResponseWriter, r *http.Request) {
	token:=strings.Split(r.URL.EscapedPath(),"/")[2]
	// log.Println("token: ",token)
	stream,err:=objects.GetRSResumablePutStreamFromToken(token)
	if err!=nil{
		log.Println("Get stream error: ",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	current:=stream.CurrentSize()
	// log.Println("size: ",current)
	if current==-1{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset:=httpTool.GetOffsetFromHeader(r.Header)
	// log.Println("offset: ",offset)
	if offset!=current{
		log.Println("offset与当前位置不符")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	buf:=make([]byte,rs.CHUNK_SIZE*rs.NUM_DATA_SHARES)
	for {
		n,e:=r.Body.Read(buf)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		// log.Println(current,n,stream.Token.Size)
		if current > stream.Token.Size {
			stream.RSPutStream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.CHUNK_SIZE*rs.NUM_DATA_SHARES && current != stream.Token.Size {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		stream.RSPutStream.Write(buf[:n])
		if current == stream.Token.Size{
			// stream.RSPutStream.Flush()
			// log.Println("flush3")
			getStream, e := objects.CreateRSResumableGetStream(stream.Token.Servers, stream.Token.Uuids, stream.Token.Size)
			hash := calculateHash(getStream)
			if hash != stream.Token.Hash {
				stream.RSPutStream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if objects.Exist(hash) {
				stream.RSPutStream.Commit(false)
				log.Println("resumable put done but object already exist")
			} else {
				stream.RSPutStream.Commit(true)
				log.Println("resumable put done")
			}
			e = es.AddVersion(stream.Token.Name, stream.Token.Hash, stream.Token.Size)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}