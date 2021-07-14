package singleflight

import "sync"

//call 代表正在进行中或已经结束的请求。
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//Group 管理不同key的请求（call）
type Group struct {
	mu sync.Mutex //保护 Group 的成员变量 m 不被并发读写
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call) //延迟实例化
	}
	if c, ok := g.m[key]; ok { // ok说明 key正在被fetch， 所以不需要再执行fn
		g.mu.Unlock()
		c.wg.Wait() //被阻塞说明之前fn未执行完
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
