package apiutil

import (
    wm "github.com/pomack/webmachine.go/webmachine"
)

func UserIdFromRequestUrl(req wm.Request) string {
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen >= 6 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "u" {
            return path[5]
        }
    }
    if pathLen >= 3 {
        if path[0] == "" && path[1] == "u" {
            return path[2]
        }
    }
    return ""
}

