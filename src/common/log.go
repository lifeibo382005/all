package common

import (
    log "code.google.com/p/log4go"
)

func init() {
    log.LoadConfiguration("conf/log4go.xml")
}


