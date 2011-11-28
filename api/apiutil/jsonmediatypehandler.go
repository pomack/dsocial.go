package apiutil

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "bytes"
    "http"
    "io"
    "io/ioutil"
    "json"
    "strconv"
    "time"
)

type JSONMediaTypeHandler struct {
    obj                 jsonhelper.JSONObject
    lastModified        *time.Time
    etag                string
    writtenStatusHeader bool
}

func NewJSONMediaTypeHandler(obj jsonhelper.JSONObject, lastModified *time.Time, etag string) *JSONMediaTypeHandler {
    return &JSONMediaTypeHandler{
        obj:          obj,
        lastModified: lastModified,
        etag:         etag,
    }
}

func (p *JSONMediaTypeHandler) MediaTypeOutput() string {
    return wm.MIME_TYPE_JSON
}

func (p *JSONMediaTypeHandler) MediaTypeHandleOutputTo(req wm.Request, cxt wm.Context, writer io.Writer, resp wm.ResponseWriter) {
    buf := bytes.NewBufferString("")
    obj := jsonhelper.NewJSONObject()
    enc := json.NewEncoder(buf)
    obj.Set("status", "success")
    obj.Set("result", p.obj)
    err := enc.Encode(obj)
    if err != nil {
        //resp.Header().Set("Content-Type", wm.MIME_TYPE_JSON)
        if !p.writtenStatusHeader {
            resp.WriteHeader(500)
            p.writtenStatusHeader = true
        }
        m := jsonhelper.NewJSONObject()
        w := json.NewEncoder(writer)
        m.Set("status", "error")
        m.Set("message", err.String())
        m.Set("result", nil)
        w.Encode(m)
        return
    }
    //resp.Header().Set("Content-Type", wm.MIME_TYPE_JSON)
    resp.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
    if p.lastModified != nil {
        resp.Header().Set("Last-Modified", p.lastModified.Format(http.TimeFormat))
    }
    if len(p.etag) > 0 {
        resp.Header().Set("ETag", strconv.Quote(p.etag))
    }
    handler := wm.NewPassThroughMediaTypeHandler(wm.MIME_TYPE_JSON, ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())), int64(buf.Len()), p.lastModified)
    handler.MediaTypeHandleOutputTo(req, cxt, writer, resp)
}

func (p *JSONMediaTypeHandler) MediaTypeHandler() wm.MediaTypeHandler {
    return p
}
