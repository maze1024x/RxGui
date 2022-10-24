package core

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"bytes"
	"errors"
	"strings"
	"strconv"
	"net/url"
	"net/http"
	"sync/atomic"
	"rxgui/standalone/util"
)


type Request struct {
	Method       RequestMethod
	Endpoint     *url.URL
	AuthToken    string
	BodyContent  [] byte
}
func (req *Request) Observe(lg Logger) Observable {
	return Observable(func(pub DataPublisher) {
		go sendRequest(req, pub, lg)
	})
}

type RequestMethod string
const (
	GET RequestMethod = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"
	SUBSCRIBE = "SUBSCRIBE"
)
func (m RequestMethod) ToHttpMethod() (string, error) {
	switch m {
	case GET:    return http.MethodGet, nil
	case POST:   return http.MethodPost, nil
	case PUT:    return http.MethodPut, nil
	case DELETE: return http.MethodDelete, nil
	default:     return "", errors.New("unsupported method: " + string(m))
	}
}

type RequestPipe struct {
	name            string
	sourceFile      *os.File
	sinkFile        *os.File
	readError       error
	commandQueue    chan func()
	nextRequestId   int64
	activeRequests  map[int64] pipeRespChan
}
type pipeRespChan chan(func() ([] byte, error))
func pipeResp(content ([] byte)) func() ([] byte, error) {
	return func() ([] byte, error) { return content, nil }
}
func pipeRespError(err error) func() ([] byte, error) {
	return func() ([] byte, error) { return nil, err }
}
var singletonStdioPipe = (*RequestPipe)(nil)
func stdioPipe() *RequestPipe {
	if singletonStdioPipe == nil {
		singletonStdioPipe = createRequestPipe("stdio", os.Stdin, os.Stdout)
	}
	return singletonStdioPipe
}
func (req *Request) Pipe() (*RequestPipe, error, bool) {
	var scheme = req.Endpoint.Scheme
	var host = req.Endpoint.Host
	if scheme == "pipe" {
		if host == "stdio" {
			return stdioPipe(), nil, true
		} else {
			var err = errors.New("unknown pipe: " + host)
			return nil, err, true
		}
	}
	return nil, nil, false
}
func createRequestPipe(name string, source *os.File, sink *os.File) *RequestPipe {
	var pipe = &RequestPipe {
		name:           name,
		sourceFile:     source,
		sinkFile:       sink,
		readError:      nil,
		commandQueue:   make(chan func(), 256),
		nextRequestId:  0,
		activeRequests: make(map[int64] pipeRespChan),
	}
	go (func() {
		for k := range pipe.commandQueue {
			k()
		}
	})()
	pipe.commandQueue <- func() {
		var read = func() error {
			var r = bufio.NewReader(pipe.sourceFile)
			for {
				var line, _, err = util.WellBehavedFscanln(r)
				if err != nil { return err }
				if line == "" {
					continue
				}
				var t = strings.Split(line, " ")
				if !(2 < len(t)) { goto invalid }
				{ var kind = t[0]
				var id_str = t[1]
				var length_str = t[2]
				var id, err1 = strconv.ParseInt(id_str, 10, 64)
				if err1 != nil { goto invalid }
				var length, err2 = strconv.Atoi(length_str)
				if err2 != nil { goto invalid }
				var content = make([] byte, length)
				if length > 0 {
					var _, err = io.ReadFull(r, content)
					if err != nil { return err }
				}
				switch kind {
				case "OK":
					pipe.commandQueue <- func() {
						if resp, ok := pipe.activeRequests[id]; ok {
							resp <- pipeResp(content)
						} else {
							// no-op
						}
					}
					continue
				case "ERR":
					pipe.commandQueue <- func() {
						if resp, ok := pipe.activeRequests[id]; ok {
							var msg = util.WellBehavedDecodeUtf8(content)
							var err = errors.New(msg)
							resp <- pipeRespError(err)
						} else {
							// no-op
						}
					}
					continue
				default:
					goto invalid
				}}
				invalid:
				return errors.New("invalid response header: " + line)
			}
		}
		go (func() {
			var err = read()
			if err != nil {
				var err = fmt.Errorf("pipe read error: %w", err)
				pipe.commandQueue <- func() {
					pipe.readError = err
					for _, resp := range pipe.activeRequests {
						resp <- pipeRespError(err)
					}
				}
			}
		})()
	}
	return pipe
}
func (pipe *RequestPipe) addRequest(req *Request) (pipeRespChan, int64, func()) {
	var resp = make(pipeRespChan, 256)
	var id = atomic.AddInt64(&(pipe.nextRequestId), 1)
	pipe.commandQueue <- func() {
		pipe.activeRequests[id] = resp
		var method = req.Method
		var path = req.Endpoint.Path
		if path == "" { path = "/" }
		var token = strconv.Quote(req.AuthToken)
		var length = len(req.BodyContent)
		var write = func() error {
			if pipe.readError != nil {
				return pipe.readError
			}
			var _, err = fmt.Fprintf(pipe.sinkFile,
				"REQ %d %s %s %s %d\n", id, method, path, token, length)
			if err != nil { return err }
			if length > 0 {
				var _, err = pipe.sinkFile.Write(req.BodyContent)
				if err != nil { return err }
			}
			return nil
		}
		var err = write()
		if err != nil {
			resp <- pipeRespError(err)
		}
	}
	var remove = func() {
		pipe.commandQueue <- func() {
			delete(pipe.activeRequests, id)
		}
	}
	return resp, id, remove
}
func (pipe *RequestPipe) cancelRequest(id int64) {
	pipe.commandQueue <- func() {
		var _, err = fmt.Fprintf(pipe.sinkFile,
			"CANCEL %d\n", id)
		if err != nil {
			var err = fmt.Errorf(
				"%s: request %d: error sending cancel signal: %w",
				pipe.name, id, err,
			)
			println(err.Error())
		}
	}
}

func sendRequest(req *Request, pub DataPublisher, lg Logger) {
	lg.LogRequest(req)
	if pipe, err, ok := req.Pipe(); ok {
		if err != nil { pub.AsyncThrow(err); return }
		var resp, id, remove = pipe.addRequest(req)
		defer remove()
		if req.Method == SUBSCRIBE {
			var yield, complete = pub.AsyncGenerate()
			var throw = pub.AsyncThrow
			loop: for {
				select {
				case resp_ := <- resp:
					var resp, err = resp_()
					if err != nil {
						throw(err)
						break loop
					}
					if len(resp) > 0 {
						yield(ObjBytes(resp))
						continue loop
					} else {
						complete()
						break loop
					}
				case <- pub.AsyncContext().Done():
					pipe.cancelRequest(id)
					break loop
				}
			}
		} else {
			select {
			case resp_ := <- resp:
				var resp, err = resp_()
				if err != nil { pub.AsyncThrow(err); return }
				pub.AsyncReturn(ObjBytes(resp))
			case <- pub.AsyncContext().Done():
				pipe.cancelRequest(id)
			}
		}
	} else {
		var ctx = pub.AsyncContext()
		var method, err = req.Method.ToHttpMethod()
		if err != nil { pub.AsyncThrow(err); return }
		var endpoint = req.Endpoint.String()
		var body = bytes.NewReader(req.BodyContent)
		var token = req.AuthToken
		{ var req, err = http.NewRequestWithContext (
			ctx, method, endpoint, body,
		)
		if err != nil { pub.AsyncThrow(err); return }
		if token != "" {
			req.Header.Set("X-Auth-Token", token)
		}
		{ var res, err = http.DefaultClient.Do(req)
		if err != nil { pub.AsyncThrow(err); return }
		defer (func() {
			_ = res.Body.Close()
		})()
		var status = res.StatusCode
		var ok = (200 <= status && status < 300)
		if !(ok) {
			var err = errors.New(fmt.Sprintf("HTTP %d", status))
			{ pub.AsyncThrow(err); return }
		}
		{ var binary, err = io.ReadAll(res.Body)
		if err != nil { pub.AsyncThrow(err); return }
		pub.AsyncReturn(ObjBytes(binary)) }}}
	}
}


