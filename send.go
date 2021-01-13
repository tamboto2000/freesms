package freesms

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SendMsg send message
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

	resp, err := http.DefaultClient.Do(req)
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

	if strings.Contains(string(raw), "Mohon Tunggu 15 Menit Lagi Untuk Pengiriman Pesan Yang Sama") {
		return errors.New("wait for 15 minutes for sending the same message, or try to use proxy")
	}

	return nil
}