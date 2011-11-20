package apiutil

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "http"
    "io"
    "json"
    "os"
    "time"
)

func OutputErrorMessage(writer io.Writer, message string, result interface{}, statusCode int, headers http.Header) (int, http.Header, os.Error) {
    if statusCode == 0 {
        statusCode = 500
    }
    if headers == nil {
        headers = make(http.Header)
    }
    //headers.Set("Content-Type", wm.MIME_TYPE_JSON)
    m := jsonhelper.NewJSONObject()
    w := json.NewEncoder(writer)
    m.Set("status", "error")
    m.Set("message", message)
    m.Set("result", result)
    w.Encode(m)
    return statusCode, headers, nil
}

func OutputJSONObject(writer io.Writer, obj jsonhelper.JSONObject, lastModified *time.Time, etag string, statusCode int, headers http.Header) (int, http.Header, os.Error) {
    if statusCode == 0 {
        statusCode = 200
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
    w := json.NewEncoder(writer)
    m.Set("status", "success")
    m.Set("result", obj)
    w.Encode(m)
    return statusCode, headers, nil
}
