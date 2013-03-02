package taoke

import (
    "fmt"
    "bytes"
    "common"
    "errors"
    "io/ioutil"
    "encoding/json"
    "github.com/mahonia"
    log "code.google.com/p/log4go"
)

type ItemInfo struct {
    Date string
    Id string
    Count string
    Price string
    State string
    Transaction string
    Commission string
    Income string
}

func GetTaokeDetail(account, startTime, endTime string) (data []byte, err error) {

    log.Info("request: %s, %s, %s", account, startTime, endTime)
    searchurl := fmt.Sprintf("http://www.alimama.com/union/newreport/taobaokeDetail.htm?toPage=1&perPageSize=2000&startTime=%s&endTime=%s", startTime, endTime)

    body, err := common.GetPage(account, searchurl)
    if err != nil {
        return nil, err
    }

    d:=mahonia.NewDecoder("gbk")
    r := d.NewReader(bytes.NewBuffer(body))
    body, _ = ioutil.ReadAll(r)

    log.Info(string(body))

    /* login */

    i := bytes.Index(body, []byte("<title>淘宝联盟-阿里妈妈登录页面</title>"))
    if i != -1 {
        return nil, errors.New("account need login.")
    }


    i = bytes.Index(body, []byte("<table class=\"med-table med-list-s\">"))
    if i == -1 {
        return nil, errors.New("parse taoke detail page failed")
    }

    start := bytes.Index(body[i:], []byte("<tbody>"))
    if start == -1 {
        return nil, errors.New("parse taoke detail page failed")
    }

    i = i + start + len("<tbody>")

    end := bytes.Index(body[i:], []byte("</tbody>"))
    if end == -1 {
        return nil, errors.New("parse taoke detail page failed")
    }

    /* error */
    ei := bytes.Index(body[i:], []byte("<div class=\"med-tip\">")) 
    if ei != -1 {
        return []byte("[]"), nil
    }

    trs := bytes.Split(bytes.TrimSpace(body[i:i+end]), []byte("<tr>"))

    items := make([]ItemInfo, 0)

    for _, tr := range(trs) {
        if len(tr) == 0 {
            continue
        }

        i = bytes.Index(tr, []byte("</tr>"))
        if i == -1 {
            return nil, errors.New("parse taoke detail page failed")
        }
        tr = bytes.TrimSpace(tr[:i])

        tds := bytes.Split(tr, []byte("<td"))

        item := ItemInfo{}

        for index, td := range(tds) {
            if len(td) == 0 {
                continue
            }
            i = bytes.Index(td, []byte("</td>"))
            if i == -1 {
                return nil, errors.New("parse taoke detail page failed")
            }
            td = bytes.TrimSpace(td[:i])

            switch index {
            case 1:
                i = bytes.Index(td, []byte(">"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Date = string(td[i+1:])

            case 2:
                i = bytes.Index(td, []byte("id="))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+3:]

                i = bytes.Index(td, []byte("\""))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Id = string(td[:i])
            case 3:
                i = bytes.Index(td, []byte("2\">"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+3:]

                i = bytes.Index(td, []byte("<"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Count = string(td[:i])
            case 4:
                i = bytes.Index(td, []byte("/i>"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+3:]

                i = bytes.Index(td, []byte("<"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Price = string(td[:i])

            case 5:
                i = bytes.Index(td, []byte("<span class="))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }


                td = td[i:]

                i = bytes.Index(td, []byte(">"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+1:]

                i = bytes.Index(td, []byte("<"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.State = string(td[:i])

            case 6:
                i = bytes.Index(td, []byte("/i>"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+3:]

                i = bytes.Index(td, []byte("<"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Transaction = string(td[:i])
            case 7:
                i = bytes.Index(td, []byte("2\">"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+3:]

                i = bytes.Index(td, []byte("<"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Commission = string(td[:i])
            case 8:
                i = bytes.Index(td, []byte("/i>"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                td = td[i+3:]

                i = bytes.Index(td, []byte("<"))
                if i == -1 {
                    return nil, errors.New("parse taoke detail page failed")
                }

                item.Income = string(td[:i])
            }

        }

        items = append(items, item)
    }

    b, err := json.Marshal(items)
    if err != nil {
        return nil, err
    }

    return b, nil
}
