package oauth

import (
	"http"
	"os"
	"io"
	"net"
//	"encoding/base64"
	"strings"
	"bufio"
	"crypto/tls"
)



// get Taken from the golang source modifed to allow headers to be passed and no redirection allowed
func get(url_ string,  headers map[string]string, timeout int64) (r *http.Response, err os.Error) {

	req, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	r, err = send(req, timeout)
	if err != nil { return }
	return
}

// post taken from Golang modified to allow Headers to be pased
func post(url_ string, headers map[string]string, body io.Reader, timeout int64) (r *http.Response, err os.Error) {
	req, err := http.NewRequest("POST", url_, nopCloser{body})
	if err != nil {
		return
	}
    req.ProtoMajor = 1
    req.ProtoMinor = 1
    req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}
    req.TransferEncoding = []string{"chunked"}

    return send(req, timeout)
}

// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//

// From the http package - modified to allow Headers to be sent to the Post method
type nopCloser struct {
    io.Reader
}

func (nopCloser) Close() os.Error { return nil }

type readClose struct {
    io.Reader
    io.Closer
}

func send(req *http.Request, timeout int64) (resp *http.Response, err os.Error) {
    if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
        return nil, nil
    }

    addr := req.URL.Host
    if !hasPort(addr) {
        addr += ":" + req.URL.Scheme
    }
    /*info := req.URL.Userinfo
    if len(info) > 0 {
        enc := base64.URLEncoding
        encoded := make([]byte, enc.EncodedLen(len(info)))
        enc.Encode(encoded, []byte(info))
        if req.Header == nil {
            req.Header = make(map[string]string)
        }
        req.Header["Authorization"] = "Basic " + string(encoded)
    }
*/
    var conn io.ReadWriteCloser
    if req.URL.Scheme == "http" {
		conn_, err_ := net.Dial("tcp", addr)
		if conn_ != nil && timeout > 0 {
			conn_.SetTimeout(timeout)
		}
		conn, err = conn_, err_
    } else { // https
		conn_, err_ := tls.Dial("tcp", addr, nil)
		if conn_ != nil && timeout > 0 {
			conn_.SetTimeout(timeout)
		}
		conn, err = conn_, err_
    }
    if err != nil {
        return nil, err
    }

    err = req.Write(conn)
    if err != nil {
        conn.Close()
        return nil, err
    }

    reader := bufio.NewReader(conn)
    resp, err = http.ReadResponse(reader, req)
    if err != nil {
        conn.Close()
        return nil, err
    }

    resp.Body = readClose{resp.Body, conn}

    return
}


func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }
