package freesms

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SendMsg send message.
// Note that there will be no any JSON, XML, or any other returned response, is err == nil then
// the message is successfully sent to receiver.
// Minimum chars is 15, maximum 122
func (cl *Client) SendMsg(phone, msg string) error {
	sendBody := make(url.Values)
	sendBody.Set("nohp", phone)
	sendBody.Set("pesan", msg)
	sendBody.Set("captcha", cl.captcha)
	sendBody.Set("key", cl.key)

	req, err := http.NewRequest("POST", "https://"+host+"/send.php", strings.NewReader(sendBody.Encode()))
	if err != nil {
		return err
	}

	req.Header = baseHeader()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "https://alpha.payuterus.biz")
	req.Header.Set("referer", "https://alpha.payuterus.biz/index.php")

	for _, c := range cl.cookies {
		req.AddCookie(c)
	}

	httpCl := http.DefaultClient
	if cl.prox != nil {
		httpCl.Transport = &http.Transport{Proxy: http.ProxyURL(cl.prox)}
	}

	resp, err := httpCl.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		return errors.New(resp.Status)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	rawString := string(raw)
	if strings.Contains(rawString, "Untuk Pengiriman Pesan Yang Sama") {
		split := strings.Split(rawString, "<br>Mohon Tunggu ")
		split = strings.Split(split[1], " Menit Lagi")

		return errors.New("wait for " + split[0] + " minutes for sending the same message")
	}

	return nil
}
