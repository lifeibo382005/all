package common

import (
    "fmt"
    "time"
    "errors"
    "strings"
    "io/ioutil"
    "net/url"
    "net/http"
    "github.com/cookiejar"
    log "code.google.com/p/log4go"
)


type TaokeClient struct {
    http.Client
    url string
}


func (tc *TaokeClient) keepalive(sitek string) {
    go func() {
        for {
            time.Sleep(time.Second * 60)
            _, _ = tc.Get("http://www.alimama.com/")
        }
    }()
}


var HttpClient map[string]*TaokeClient = make(map[string]*TaokeClient)


func Login(site, sitek, ustr string) error {

    u, err := url.Parse(ustr)
    if err != nil {
        return err
    }

    accountstr, err := Conf.String(site, "accounts", "")
    if err != nil {
        return err
    }

    if accountstr == "" {
        return errors.New("accounts not found in config.")
    }

    accounts := strings.Split(accountstr, ",")

    for _, account := range(accounts) {
        cookiestr, err := Conf.String(account, "cookies", "")
        if err != nil {
            return err
        }

        if cookiestr == "" {
            return errors.New("Cookies not found in config.")
        }

        log.Info("Read url and cookie from config of %s.", site)

        cos := strings.Split(cookiestr, ";")

        cookies := []*http.Cookie{}

        for _, co := range(cos) {

            in := strings.Index(co, "=")
            if in == -1 {
                return errors.New("Invalid cookies")
            }

            c := &http.Cookie{
                Name:co[:in],
                Value:co[in+1:],
                Raw:co,
            }
            cookies = append(cookies, c)
        }

        jar := cookiejar.NewJar(false)

        jar.SetCookies(u, cookies)

        tc := &TaokeClient{http.Client{Jar:jar}, ustr}
        HttpClient[account] = tc
        tc.keepalive(sitek)
    }

    log.Info("Parse cookie and url successed.")

    return nil
}


func GetPage(account, u string) (body []byte, err error) {

    client, ok := HttpClient[account]
    if !ok {
        return nil, errors.New(fmt.Sprintf("account '%s' notfound", account))
    }

    req, err := http.NewRequest("GET", u, nil)
    req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1312.57 Safari/537.17")
    resp, e := client.Do(req)
    if e != nil {
        return nil, e
    }

    body, err = ioutil.ReadAll(resp.Body)

    return
}
