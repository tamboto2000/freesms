package freesms

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tamboto2000/htmltojson"
)

const (
	userAgent = "Mozilla/5.0 (Linux; Android 9; Redmi Note 6 Pro Build/PKQ1.180904.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/87.0.4280.141 Mobile Safari/537.36"
	accept    = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	reqWith   = "com.smsGratisSeluruhIndonesia64"
	accLang   = "en-US,en;q=0.9"
	host      = "alpha.payuterus.biz"
)

var opList = []string{"+", "-", "/"}

func baseHeader() http.Header {
	header := make(http.Header)
	header.Set("Host", host)
	header.Set("cache-control", "max-age=0")
	header.Set("upgrade-insecure-request", "1")
	header.Set("user-agent", userAgent)
	header.Set("accept", accept)
	header.Set("x-requested-with", reqWith)
	header.Set("sec-fetch-site", "none")
	header.Set("sec-fetch-mode", "navigate")
	header.Set("sec-fetch-user", "?1")
	header.Set("sec-fetch-dest", "document")
	header.Set("accept-language", "en-US,en;q=0.9")

	return header
}

// Client holds session and perform request.
// DO NOT SEND SMS CONCURRENTLY, create new client for sending new message
type Client struct {
	cookies []*http.Cookie
	raw     []byte
	captcha string
	key     string
	prox    *url.URL
}

// NewClient initiate new client
func NewClient() (*Client, error) {
	cl := new(Client)
	if err := cl.init(); err != nil {
		return nil, err
	}

	return cl, nil
}

// SetProxy set proxy to client.
// Use proxy originated from Indonesia, see sslproxies.org, or scrape it with github.com/tamboto2000/sslproxies
func (cl *Client) SetProxy(ustr string) error {
	u, err := url.Parse(ustr)
	if err != nil {
		return err
	}

	cl.prox = u

	return nil
}

// get session and tokens
func (cl *Client) init() error {
	req, err := http.NewRequest("GET", "https://"+host+"/index.php", nil)
	if err != nil {
		panic(err.Error())
	}

	req.Header = baseHeader()
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

	// find chaptcha and key
	node, err := htmltojson.ParseFromReader(resp.Body)
	if err != nil {
		return err
	}

	node = htmltojson.SearchNode("", "form", "", "", "", node)
	if node == nil {
		return errors.New("unknown error")
	}

	// find captcha
	nodes := htmltojson.SearchAllNode("", "li", "", "", "", node)
	for _, node := range nodes {
		// there's just this weird bug that I need to copy iteration result...
		cpyNode := node
		if htmltojson.SearchNode("", "input", "", "name", "captcha", &cpyNode) != nil {
			span := htmltojson.SearchNode("", "span", "", "", "", &cpyNode)
			if span == nil {
				return errors.New("unknown error")
			}

			str := strings.ReplaceAll(span.Child[0].Data, " ", "")
			str = strings.ReplaceAll(str, "=", "")
			var operand string

			for _, op := range opList {
				if strings.Contains(str, op) {
					operand = op

					break
				}
			}

			if operand == "" {
				return errors.New("unknown error")
			}

			split := strings.Split(str, operand)

			i1, err := strconv.Atoi(split[0])
			if err != nil {
				return err
			}

			i2, err := strconv.Atoi(split[1])
			if err != nil {
				return err
			}

			cl.captcha = strconv.Itoa(i1 + i2)

			break
		}
	}

	// find key
	newNode := htmltojson.SearchNode("", "input", "", "name", "key", node)
	if newNode == nil {
		return errors.New("unknown error")
	}

	for _, attr := range newNode.Attr {
		if attr.Key == "value" {
			cl.key = attr.Val
			break
		}
	}

	// parse cookies
	cl.cookies = resp.Cookies()

	return nil
}
