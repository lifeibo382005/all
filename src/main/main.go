package main

import (
    "os"
    "fmt"
    "net/http"
    "bufio"
    "time"
    "common"
    "taoke"
    log "code.google.com/p/log4go"
)

func ErrorExit() {
    time.Sleep(time.Second)
    reader := bufio.NewReader(os.Stdin)
    _, _, _ = reader.ReadLine()
    os.Exit(-1)
}


func taokeHandler(w http.ResponseWriter, r *http.Request) {

    account := r.FormValue("account")
    if account == "" {
        fmt.Fprintf(w, "{\"error\":1, \"msg\":\"error, account is nil. eg.http://localhost/taoke?account=account1&startTime=2013-1-1&endTime=2013-3-1\"}")
        return
    }

    startTime := r.FormValue("startTime")
    endTime := r.FormValue("endTime")

    b, e := taoke.GetTaokeDetail(account, startTime, endTime)
    if e != nil {
        log.Error(e)
        fmt.Fprintf(w, "{\"error\":1, \"msg\":\"%s\"}", e.Error())
        return
    }

    fmt.Fprintf(w, "{\"error\":0, \"data\":%s}", string(b))
}


func main() {

    if err := common.Login(); err != nil {
        log.Error(err)
        ErrorExit()
    }

    port, e := common.Conf.Int("taoke", "port", 8080)
    if e != nil {
        log.Error(e)
        ErrorExit()
    }

    http.HandleFunc("/taoke", taokeHandler)

    e = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
    if e != nil {
        log.Error(e)
        ErrorExit()
    }
}