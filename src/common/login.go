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
}


func (tc *TaokeClient) keepalive() {
    go func() {
        for {
            time.Sleep(time.Second * 60)
            _, _ = tc.Get("http://www.alimama.com/")
        }
    }()
}


var HttpClient map[string]*TaokeClient = make(map[string]*TaokeClient)


func Login(site, ustr string) error {

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

        tc := &TaokeClient{http.Client{Jar:jar}}
        HttpClient[account] = tc
        tc.keepalive()
    }

    log.Info("Parse cookie and url successed.")

    return nil
}


func GetPage(account, u string) (body []byte, err error) {

    client, ok := HttpClient[account]
    if !ok {
        return nil, errors.New(fmt.Sprintf("account '%s' notfound", account))
    }

    resp, e := client.Get(u)
    if e != nil {
        return nil, e
    }

    body, err = ioutil.ReadAll(resp.Body)

    return
}
