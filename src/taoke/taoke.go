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
    Name string
    ShopId string
    ShopName string
    Count string
    Price string
    State string
    Transaction string
    Commission string
    Income string
}

func GetTaokeDetail(account, startTime, endTime string) (data []byte, err error) {

    log.Info("request: %s, %s, %s", account, startTime, endTime)

    items := make([]ItemInfo, 0)
    page := 1
    for {
        have := false

        searchurl := fmt.Sprintf("http://u.alimama.com/union/newreport/taobaokeDetail.htm?toPage=%d&perPageSize=20&startTime=%s&endTime=%s", page, startTime, endTime)


        log.Error(searchurl)

        body, err := common.GetPage(account, searchurl)
        if err != nil {
            return nil, err
        }

        i := bytes.Index(body, []byte("charset=GBK"))
        if i != -1 {
            d:=mahonia.NewDecoder("gbk")
            r := d.NewReader(bytes.NewBuffer(body))
            body, _ = ioutil.ReadAll(r)
        }

        /* login */

        i = bytes.Index(body, []byte("<title>阿里妈妈-阿里妈妈登录页面</title>"))
        if i != -1 {
            return nil, errors.New("account need login.")
        }

        /* when parse error, log page */
        defer func() {
            if data == nil {
                log.Error(string(body))
            }
        }()

        i = bytes.Index(body, []byte("<table class=\"med-table med-list-s\">"))
        if i == -1 {
            return nil, errors.New("1parse taoke detail page failed")
        }

        start := bytes.Index(body[i:], []byte("<tbody>"))
        if start == -1 {
            return nil, errors.New("2parse taoke detail page failed")
        }

        i = i + start + len("<tbody>")

        end := bytes.Index(body[i:], []byte("</tbody>"))
        if end == -1 {
            return nil, errors.New("3parse taoke detail page failed")
        }

        /* error */
        ei := bytes.Index(body[i:], []byte("<div class=\"med-tip\">")) 
        if ei != -1 {
            break
        }

        trs := bytes.Split(bytes.TrimSpace(body[i:i+end]), []byte("<tr>"))

        for _, tr := range(trs) {
            if len(tr) == 0 {
                continue
            }

            i = bytes.Index(tr, []byte("</tr>"))
            if i == -1 {
                return nil, errors.New("4parse taoke detail page failed")
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
                    return nil, errors.New("5parse taoke detail page failed")
                }
                td = bytes.TrimSpace(td[:i])

                switch index {
                case 1:
                    i = bytes.Index(td, []byte(">"))
                    if i == -1 {
                        return nil, errors.New("6parse taoke detail page failed")
                    }

                    item.Date = string(td[i+1:])

                case 2:
                    i = bytes.Index(td, []byte("id="))
                    if i == -1 {
                        return nil, errors.New("7parse taoke detail page failed")
                    }

                    td = td[i+3:]

                    i = bytes.Index(td, []byte("\""))
                    if i == -1 {
                        return nil, errors.New("8parse taoke detail page failed")
                    }

                    //
                    item.Id = string(td[:i])

                    td = td[i+2:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("8parse taoke detail page failed")
                    }

                    //
                    item.Name = string(td[:i])

                    td = td[i:]

                    i = bytes.Index(td, []byte("oid="))
                    if i == -1 {
                        return nil, errors.New("8parse taoke detail page failed")
                    }

                    td = td[i+4:]

                    i = bytes.Index(td, []byte("\""))
                    if i == -1 {
                        return nil, errors.New("8parse taoke detail page failed")
                    }

                    item.ShopId = string(td[:i])

                    td = td[i:]

                    i = bytes.Index(td, []byte(">"))
                    if i == -1 {
                        return nil, errors.New("8parse taoke detail page failed")
                    }

                    td = td[i+1:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("8parse taoke detail page failed")
                    }

                    item.ShopName = string(td[:i])

                case 3:
                    i = bytes.Index(td, []byte("2\">"))
                    if i == -1 {
                        return nil, errors.New("9parse taoke detail page failed")
                    }

                    td = td[i+3:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("10parse taoke detail page failed")
                    }

                    item.Count = string(td[:i])
                case 4:
                    i = bytes.Index(td, []byte("/i>"))
                    if i == -1 {
                        return nil, errors.New("11parse taoke detail page failed")
                    }

                    td = td[i+3:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("12parse taoke detail page failed")
                    }

                    item.Price = string(td[:i])

                case 5:
                    i = bytes.Index(td, []byte("<span"))
                    if i == -1 {
                        log.Info(string(td))
                        return nil, errors.New("13parse taoke detail page failed")
                    }


                    td = td[i:]

                    i = bytes.Index(td, []byte(">"))
                    if i == -1 {
                        return nil, errors.New("14parse taoke detail page failed")
                    }

                    td = td[i+1:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("15parse taoke detail page failed")
                    }

                    item.State = string(td[:i])

                case 6:
                    continue

                case 7:
                    i = bytes.Index(td, []byte("/i>"))
                    if i == -1 {
                        return nil, errors.New("16parse taoke detail page failed")
                    }

                    td = td[i+3:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("17parse taoke detail page failed")
                    }

                    item.Transaction = string(td[:i])
                case 8:
                    i = bytes.Index(td, []byte("2\">"))
                    if i == -1 {
                        return nil, errors.New("18parse taoke detail page failed")
                    }

                    td = td[i+3:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("19parse taoke detail page failed")
                    }

                    item.Commission = string(td[:i])
                case 9:
                    continue
                case 10:
                    continue
                case 11:
                    i = bytes.Index(td, []byte("/i>"))
                    if i == -1 {
                        return nil, errors.New("20parse taoke detail page failed")
                    }

                    td = td[i+3:]

                    i = bytes.Index(td, []byte("<"))
                    if i == -1 {
                        return nil, errors.New("21parse taoke detail page failed")
                    }

                    item.Income = string(td[:i])
                }
            }

            have = true

            items = append(items, item)
        }

        if !have {
            break
        }

        page++
    }

    data, err = json.Marshal(items)
    if err != nil {
        return nil, err
    }

    return data, nil
}
