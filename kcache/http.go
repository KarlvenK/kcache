package kcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_kcache/"

//HTTPPool implements PeerPicker for a pool of HTTP peers
type HTTPPool struct {
	//this peer's base URL, eg "https://example.net:8000"
	self     string //用来记录自己的地址，包括主机名/IP 和端口。
	basePath string //作为节点间通讯地址的前缀，默认是 /_kcache/
	//那么 http://example.com/_kcache/ 开头的请求，就用于节点间的访问。
	//因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 /api 作为前缀。
}

/*
在标准库中，http.Handler 接口的定义如下：

package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

所以实现了 ServeHTTP(...) 的struct 就是Handler
*/

//NewHTTPPool initializes an HTTP pool of peers
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	//  /<basePath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName, key := parts[0], parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusInternalServerError)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(view.ByteSlice())
}
