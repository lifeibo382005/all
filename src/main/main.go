package main

import (
    "os"
    "fmt"
    "runtime"
    "net/http"
    "bufio"
    "time"
    "common"
    "sync"
    "taoke"
    "yiqifa"
    log "code.google.com/p/log4go"
)

func ErrorExit() {
    time.Sleep(time.Second)
    reader := bufio.NewReader(os.Stdin)
    _, _, _ = reader.ReadLine()
    os.Exit(-1)
}

var Cache map[string][]byte = make(map[string][]byte)
var CacheLock sync.RWMutex

func cacheGet(web, account, startTime, endTime string) (ret []byte, ok bool) {
    CacheLock.RLock()
    defer CacheLock.RUnlock()
    st := web + account + startTime + endTime
    ret, ok = Cache[st]
    return
}

func cachePut(web, account, startTime, endTime string, data []byte) {
    CacheLock.Lock()
    defer CacheLock.Unlock()
    st := web + account + startTime + endTime
    Cache[st] = data
}

func cleanAll() {
    CacheLock.Lock()
    defer CacheLock.Unlock()
    Cache = make(map[string][]byte)

    runtime.GC()
}

func cleanCache() {
    go func() {
        for {
            time.Sleep(time.Second * 5)
            cleanAll()
        }
    }()
}

func taokeHandler(w http.ResponseWriter, r *http.Request) {

    account := r.FormValue("account")
    if account == "" {
        fmt.Fprintf(w, "{\"error\":1, \"msg\":\"error, account is nil. eg.http://localhost/taoke?account=account1&startTime=2013-1-1&endTime=2013-3-1\"}")
        return
    }

    startTime := r.FormValue("startTime")
    endTime := r.FormValue("endTime")

    var b []byte
    var e error
    b, ok := cacheGet("taoke", account, startTime, endTime)
    if !ok {
        b, e = taoke.GetTaokeDetail(account, startTime, endTime)
        if e != nil {
            log.Error(e)
            fmt.Fprintf(w, "{\"error\":1, \"msg\":\"%s\"}", e.Error())
            return
        }
        cachePut("taoke", account, startTime, endTime, b)
    }

    fmt.Fprintf(w, "{\"error\":0, \"data\":%s}", string(b))
}

func yiqifaHandler(w http.ResponseWriter, r *http.Request) {

    account := r.FormValue("account")
    if account == "" {
        fmt.Fprintf(w, "{\"error\":1, \"msg\":\"error, account is nil. eg.http://localhost/yiqifa?account=yiqifaaccount1&startTime=2013-1-1&endTime=2013-3-1\"}")
        return
    }

    startTime := r.FormValue("startTime")
    endTime := r.FormValue("endTime")

    var b []byte
    var e error
    b, ok := cacheGet("yiqifa", account, startTime, endTime)
    if !ok {
        b, e = yiqifa.GetCPSDetail(account, startTime, endTime)
        if e != nil {
            log.Error(e)
            fmt.Fprintf(w, "{\"error\":1, \"msg\":\"%s\"}", e.Error())
            return
        }
        cachePut("yiqifa", account, startTime, endTime, b)
    }

    fmt.Fprintf(w, "{\"error\":0, \"data\":%s}", string(b))
}

func run() {
    if err := common.Login("taoke", "http://u.alimama.com","http://u.alimama.com/union/newreport/taobaokeDetail.htm"); err != nil {
        log.Error(err)
        ErrorExit()
    }

    if err := common.Login("yiqifa", "http://www.yiqifa.com/", "http://www.yiqifa.com/"); err != nil {
        log.Error(err)
        ErrorExit()
    }

    port, e := common.Conf.Int("common", "port", 8080)
    if e != nil {
        log.Error(e)
        ErrorExit()
    }

    http.HandleFunc("/taoke", taokeHandler)
    http.HandleFunc("/yiqifa", yiqifaHandler)

    cleanCache()

    for {
        e = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
        if e != nil {
            log.Error(e)
        }

        time.Sleep(time.Second)
    }
}

func main() {
    run()
    ErrorExit()
}
