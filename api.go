package geektime

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultArticlesURL      = "https://time.geekbang.org/serv/v1/column/articles"
	defaultArticleURL       = "https://time.geekbang.org/serv/v1/article"
	defaultIntroURL         = "https://time.geekbang.org/serv/v1/column/intro"
	defaultVideoPlayAuthURL = "https://time.geekbang.org/serv/v3/source_auth/video_play_auth"
	defaultPlayListURL      = "https://vod.cn-beijing.aliyuncs.com/"
)

type course struct {
	Title        string `json:"column_title"`
	ArticleCount int    `json:"article_count"`
}

func (c *course) String() string {
	return fmt.Sprintf("course:[title:%s,count:%d]", c.Title, c.ArticleCount)
}

func fetchCourse(introURL, courseID, cookie string) (c course, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		introURL,
		strings.NewReader(fmt.Sprintf(
			`{"cid":"%s","with_groupbuy":true}`, courseID,
		)),
	)
	if err != nil {
		return
	}

	fillHeaders(cookie, req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	err = check(resp)
	if err != nil {
		return
	}

	data, err := readDataAndCloseResp(resp)
	if err != nil {
		return
	}

	var d struct {
		Course course `json:"data"`
		Code   int    `json:"code"`
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	return d.Course, nil
}

type article struct {
	ID           int    `json:"id"`
	VideoID      string `json:"video_id"`
	Title        string `json:"article_title"`
	AuditM3U8URL string `json:"audio_url"`
	MP3URL       string `json:"audio_download_url"`
}

func (a article) String() string {
	return fmt.Sprintf(
		"article:[id:%d,videoID:%s,title:%s,auditM3U8URL:%s,mp3URL:%s]",
		a.ID, a.VideoID, a.Title, a.AuditM3U8URL, a.MP3URL,
	)
}

func fetchArticles(articlesURL, courseID, cookie string) ([]article, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		articlesURL,
		strings.NewReader(fmt.Sprintf(
			`{"cid":"%s","size":500,"prev":0,"order":"earliest","sample":false}`, courseID,
		)),
	)
	if err != nil {
		return nil, err
	}

	fillHeaders(cookie, req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = check(resp)
	if err != nil {
		return nil, err
	}

	data, err := readDataAndCloseResp(resp)
	if err != nil {
		return nil, err
	}

	var d struct {
		Data struct {
			List []article `json:"list"`
		} `json:"data"`
		Code int `json:"code"`
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	if d.Code != 0 {
		return nil, errors.New("fetch acticles error:" + string(data))
	}

	return d.Data.List, nil
}

type resolution string

const (
	ld resolution = "ld"
	sd            = "sd"
	hd            = "hd"

	defaultResolution = sd
)

type video struct {
	Size    int64  `json:"size"`
	Name    string `json:"article_title"`
	VideoID string `json:"video_id"`
}

func fetchVideoOfArticle(articleURL, courseID, cookie string, articleID int) (v video, err error) {
	req, err := http.NewRequest(
		http.MethodPost,
		articleURL,
		strings.NewReader(fmt.Sprintf(`{"id":%d}`, articleID)),
	)
	if err != nil {
		return
	}

	fillHeaders(cookie, req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	err = check(resp)
	if err != nil {
		return
	}

	data, err := readDataAndCloseResp(resp)
	if err != nil {
		return
	}

	var d struct {
		Data struct {
			Videos  map[resolution]video `json:"video_preview"`
			Name    string               `json:"article_title"`
			VideoID string               `json:"video_id"`
		} `json:"data"`
		Code int `json:"code"`
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	if d.Code != 0 {
		return video{}, errors.New("fetch acticles error:" + string(data))
	}

	v, ok := d.Data.Videos[defaultResolution]

	if ok {
		v.Name, v.VideoID = d.Data.Name, d.Data.VideoID
		return
	}

	for _, v = range d.Data.Videos {
		v.Name, v.VideoID = d.Data.Name, d.Data.VideoID
		return
	}

	return video{}, errors.New("no video")
}

type videoPlayAuth struct {
	SecurityToken   string `json:"SecurityToken"`
	AuthInfo        string `json:"AuthInfo"`
	AccessKeyID     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Region          string `json:"Region"`
	PlayDomain      string `json:"PlayDomain"`
	CustomerID      int64  `json:"CustomerId"`
	VideoMeta       struct {
		Status   string  `json:"Status"`
		VideoID  string  `json:"VideoId"`
		Title    string  `json:"Title"`
		CoverURL string  `json:"CoverURL"`
		Duration float64 `json:"Duration"`
	} `json:"VideoMeta"`
}

func fetchVideoPlayAuth(
	videoPlayAuthURL, cookie string, articleID, sourceType int, videoID string,
) (
	auth videoPlayAuth, err error,
) {
	req, err := http.NewRequest(
		http.MethodPost,
		videoPlayAuthURL,
		strings.NewReader(fmt.Sprintf(
			`{"source_type":%d,"aid":%d,"video_id":"%s"}`, sourceType, articleID, videoID,
		)),
	)
	if err != nil {
		return
	}

	fillHeaders(cookie, req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	err = check(resp)
	if err != nil {
		return
	}

	data, err := readDataAndCloseResp(resp)
	if err != nil {
		return
	}

	var d struct {
		Data struct {
			Auth string `json:"play_auth"`
		} `json:"data"`
		Code int `json:"code"`
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	if d.Code != 0 {
		err = errors.New("play auth error:" + string(data))
		return
	}

	data, err = base64.StdEncoding.DecodeString(d.Data.Auth)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &auth)
	if err != nil {
		return
	}

	return
}

func fillHeaders(cookie string, header http.Header) {
	header.Set("Pragma", "no-cache")
	header.Set("Sec-Fetch-Site", "same-origin")
	header.Set("Origin", "https://time.geekbang.org")
	header.Set("Accept-Encoding", "gzip, deflate, br")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36")
	header.Set("Sec-Fetch-Mode", "cors")
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json, text/plain, */*")
	header.Set("Cache-Control", "no-cache")
	header.Set("Cookie", cookie)
}

type playList struct {
	Code      string `json:"Code"`
	RequestID string `json:"RequestId"`
	VideoBase struct {
		TranscodeMode string    `json:"TranscodeMode"`
		CreationTime  time.Time `json:"CreationTime"`
		CoverURL      string    `json:"CoverURL"`
		Status        string    `json:"Status"`
		MediaType     string    `json:"MediaType"`
		VideoID       string    `json:"VideoId"`
		Duration      string    `json:"Duration"`
		OutputType    string    `json:"OutputType"`
		Title         string    `json:"Title"`
	} `json:"VideoBase"`
	PlayInfoList struct {
		PlayInfo []struct {
			Format           string    `json:"Format"`
			Plaintext        string    `json:"Plaintext"`
			PreprocessStatus string    `json:"PreprocessStatus"`
			StreamType       string    `json:"StreamType"`
			ModificationTime time.Time `json:"ModificationTime"`
			Specification    string    `json:"Specification"`
			Height           int       `json:"Height"`
			PlayURL          string    `json:"PlayURL"`
			EncryptType      string    `json:"EncryptType"`
			Rand             string    `json:"Rand"`
			NarrowBandType   string    `json:"NarrowBandType"`
			CreationTime     time.Time `json:"CreationTime"`
			Status           string    `json:"Status"`
			JobID            string    `json:"JobId"`
			Duration         string    `json:"Duration"`
			Encrypt          int       `json:"Encrypt"`
			Width            int       `json:"Width"`
			Fps              string    `json:"Fps"`
			Bitrate          string    `json:"Bitrate"`
			Size             int       `json:"Size"`
			Definition       string    `json:"Definition"`
		} `json:"PlayInfo"`
	} `json:"PlayInfoList"`
}

func fetchPlayList(playListURL, videoID string, auth videoPlayAuth) (playList playList, err error) {
	query, err := buildQuery(http.MethodGet, videoID, auth)
	if err != nil {
		return
	}
	resp, err := http.Get(playListURL + "?" + query)
	if err != nil {
		return
	}

	err = check(resp)
	if err != nil {
		return
	}

	data, err := readDataAndCloseResp(resp)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &playList)
	if playList.Code != "" {
		err = errors.New(playList.Code)
	}
	return
}

func check(resp *http.Response) error {
	v := resp.Header.Get("X-GEEK-WARN")
	if v != "" {
		return errors.New(v)
	}
	return nil
}

func buildQuery(method, videoID string, auth videoPlayAuth) (string, error) {
	params := url.Values{
		"AccessKeyId":      []string{auth.AccessKeyID},
		"Action":           []string{"GetPlayInfo"},
		"AuthInfo":         []string{auth.AuthInfo},
		"AuthTimeout":      []string{"7200"},
		"Channel":          []string{"HTML5"},
		"Definition":       []string{""},
		"Format":           []string{"JSON"},
		"Formats":          []string{""},
		"PlayConfig":       []string{"{}"},
		"PlayerVersion":    []string{"2.8.2"},
		"ReAuthInfo":       []string{"{}"},
		"SecurityToken":    []string{auth.SecurityToken},
		"SignatureMethod":  []string{"HMAC-SHA1"},
		"SignatureNonce":   []string{newUUID().String()},
		"SignatureVersion": []string{"1.0"},
		"StreamType":       []string{"video"},
		"Version":          []string{"2017-03-21"},
		"VideoId":          []string{videoID},
	}

	query := params.Encode()
	sign := sign(auth.AccessKeySecret, buildStringToSign(method, query))

	return query + "&Signature=" + url.QueryEscape(sign), nil
}

func sign(accessKeySecret, stringToSign string) string {
	return shaHMAC1(stringToSign, accessKeySecret+"&")
}

func shaHMAC1(source, secret string) string {
	key := []byte(secret)
	hmac := hmac.New(sha1.New, key)
	hmac.Write([]byte(source))
	signedBytes := hmac.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func buildStringToSign(method, query string) string {
	query = strings.Replace(query, "+", "%20", -1)
	query = strings.Replace(query, "*", "%2A", -1)
	query = strings.Replace(query, "%7E", "~", -1)
	query = url.QueryEscape(query)
	return method + "&%2F&" + query
}

func download(url, outDir string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	data, err := readDataAndCloseResp(resp)
	if err != nil {
		return "", err
	}

	path := outDir + string(os.PathSeparator) + filename(url)
	err = writeFile(path, data)

	return path, err
}

func filename(url string) string {
	return url[strings.LastIndexByte(url, '/')+1:]
}

func writeFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	f.Write(data)
	return f.Close()
}
