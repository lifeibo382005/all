package yiqifa

import (
    "fmt"
    "errors"
    "common"
    "archive/zip"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "github.com/mahonia"
    log "code.google.com/p/log4go"
)

func GetCPSDetail(account, startTime, endTime string) (data []byte, err error) {
    log.Info("request: %s, %s, %s", account, startTime, endTime)

    searchurl := fmt.Sprintf("http://www.yiqifa.com/earner/earnerExportCpsEffectOriList.do?schStartDate=&schEndDate=&back=&effectDateOrderby=&balanceDateOrderby=&commissionOrderby=&orderNoOrderby=&productNoOrderby=&sysWebsiteCommisionOrderby=&pageNumber=1&pageSize=10&searchOption=orderNo&startDate=%s&endDate=%s&startConfirmDate=&endConfirmDate=&websiteId=&campaignType=&campaignName=&schCampaignId=0&searchOptionValue=&confirmStatus=&dataSourceType=&perSize=10&perSize2=10", startTime, endTime)

    body, err := common.GetPage(account, searchurl)
    if err != nil {
        log.Info(err)
        return nil, err
    }

    r, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
    if err != nil {

        d:=mahonia.NewDecoder("gbk")
        r := d.NewReader(bytes.NewBuffer(body))
        body, _ = ioutil.ReadAll(r)

        if bytes.Index(body, []byte("会员登录")) != -1 {
            return nil, errors.New("account need login.")
        }

        /* login failed */
        log.Error(string(body))
        return nil, errors.New("fetch failed.")
    }

    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            log.Info(err)
        }

        body, err = ioutil.ReadAll(rc)

        d:=mahonia.NewDecoder("gbk")
        r := d.NewReader(bytes.NewBuffer(body))
        body, _ = ioutil.ReadAll(r)

        rc.Close()
    }

    lines := bytes.Split(body, []byte("\n"))
    lines = lines[:len(lines)-2]
    items := make([][]string, len(lines))
    for i, line := range(lines) {
        cols := bytes.Split(line, []byte(","))
        items[i] = make([]string, len(cols))
        for j, col := range(cols) {
            items[i][j] = string(col[1:len(col)-1])
        }
    }

    data, err = json.Marshal(items)
    if err != nil {
        return nil, err
    }

    return data, nil
}
