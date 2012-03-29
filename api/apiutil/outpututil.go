package apiutil

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "http"
    "io"
    "json"
    "os"
    "time"
)

type jsonWriter struct {
    obj jsonhelper.JSONObject
}

func newJSONWriter(obj jsonhelper.JSONObject) *jsonWriter {
    return &jsonWriter{obj: obj}
}

func (p *jsonWriter) WriteTo(writer io.Writer) (n int64, err os.Error) {
    w := json.NewEncoder(writer)
    err = w.Encode(p.obj)
    return
}

func (p *jsonWriter) String() string {
    b, err := json.Marshal(p.obj)
    if err != nil {
        return err.String()
    }
    return string(b)
}

func OutputErrorMessage(message string, result interface{}, statusCode int, headers http.Header) (int, http.Header, io.WriterTo) {
    if statusCode == 0 {
        statusCode = http.StatusInternalServerError
    }
    if headers == nil {
        headers = make(http.Header)
    }
    //headers.Set("Content-Type", wm.MIME_TYPE_JSON)
    m := jsonhelper.NewJSONObject()
    m.Set("status", "error")
    m.Set("message", message)
    m.Set("result", result)
    return statusCode, headers, newJSONWriter(m)
}

func OutputJSONObject(obj jsonhelper.JSONObject, lastModified *time.Time, etag string, statusCode int, headers http.Header) (int, http.Header, io.WriterTo) {
    if statusCode == 0 {
        statusCode = http.StatusOK
    }
    if headers == nil {
        headers = make(http.Header)
    }
    //headers.Set("Content-Type", wm.MIME_TYPE_JSON)
    if lastModified != nil {
        headers.Set("Last-Modified", lastModified.Format(http.TimeFormat))
    }
    if len(etag) > 0 {
        headers.Set("ETag", etag)
    }
    m := jsonhelper.NewJSONObject()
    m.Set("status", "success")
    m.Set("result", obj)
    return statusCode, headers, newJSONWriter(m)
}

func AddNoCacheHeaders(headers http.Header) http.Header {
    if headers == nil {
        headers = make(http.Header)
    }
    headers.Set("Pragma", "no-cache")
    headers.Set("Cache-Control", "no-cache")
    headers.Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
    return headers
}
