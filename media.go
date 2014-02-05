package weixin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Upload media from local file
func (wx *Weixin) UploadMediaFromFile(mediaType string, fp string) (string, error) {
	file, err := os.Open(fp)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return wx.UploadMedia(mediaType, filepath.Base(fp), file)
}

// Download media and save to local file
func (wx *Weixin) DownloadMediaToFile(mediaId string, fp string) error {
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()
	return wx.DownloadMedia(mediaId, file)
}

// Upload media with media
func (wx *Weixin) UploadMedia(mediaType string, filename string, reader io.Reader) (string, error) {
	return uploadMedia(wx.tokenChan, mediaType, filename, reader)
}

// Download media with media
func (wx *Weixin) DownloadMedia(mediaId string, writer io.Writer) error {
	return downloadMedia(wx.tokenChan, mediaId, writer)
}

func uploadMedia(c chan accessToken, mediaType string, filename string, reader io.Reader) (string, error) {
	reqURL := weixinFileURL + "/upload?type=" + mediaType + "&access_token="
	for i := 0; i < 3; i++ {
		token := <-c
		if time.Since(token.expires).Seconds() < 0 {
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			fileWriter, err := bodyWriter.CreateFormFile("filename", filename)
			if err != nil {
				return "", err
			}
			if _, err = io.Copy(fileWriter, reader); err != nil {
				return "", err
			}
			contentType := bodyWriter.FormDataContentType()
			bodyWriter.Close()
			r, err := http.Post(reqURL+token.token, contentType, bodyBuf)
			if err != nil {
				return "", err
			}
			defer r.Body.Close()
			reply, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return "", err
			}
			var result struct {
				response
				Type      string `json:"type"`
				MediaId   string `json:"media_id"`
				CreatedAt int64  `json:"created_at"`
			}
			err = json.Unmarshal(reply, &result)
			if err != nil {
				return "", err
			} else {
				switch result.ErrorCode {
				case 0:
					return result.MediaId, nil
				case 42001: // access_token timeout and retry
					continue
				default:
					return "", errors.New(fmt.Sprintf("WeiXin upload[%d]: %s", result.ErrorCode, result.ErrorMessage))
				}
			}
		}
	}
	return "", errors.New("WeiXin upload media too many times")
}

func downloadMedia(c chan accessToken, mediaId string, writer io.Writer) error {
	reqURL := weixinFileURL + "/get?media_id=" + mediaId + "&access_token="
	for i := 0; i < 3; i++ {
		token := <-c
		if time.Since(token.expires).Seconds() < 0 {
			r, err := http.Get(reqURL + token.token)
			if err != nil {
				return err
			}
			defer r.Body.Close()
			if r.Header.Get("Content-Type") != "text/plain" {
				_, err := io.Copy(writer, r.Body)
				return err
			} else {
				reply, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return err
				}
				var result response
				if err := json.Unmarshal(reply, &result); err != nil {
					return err
				} else {
					switch result.ErrorCode {
					case 0:
						return nil
					case 42001: // access_token timeout and retry
						continue
					default:
						return errors.New(fmt.Sprintf("WeiXin download[%d]: %s", result.ErrorCode, result.ErrorMessage))
					}
				}
			}
		}
	}
	return errors.New("WeiXin download media too many times")
}

// Upload media from local file
func (w responseWriter) UploadMediaFromFile(mediaType string, filepath string) (string, error) {
	return w.wx.UploadMediaFromFile(mediaType, filepath)
}

// Download media and save to local file
func (w responseWriter) DownloadMediaToFile(mediaId string, filepath string) error {
	return w.wx.DownloadMediaToFile(mediaId, filepath)
}

// Upload media with reader
func (w responseWriter) UploadMedia(mediaType string, filename string, reader io.Reader) (string, error) {
	return w.wx.UploadMedia(mediaType, filename, reader)
}

// Download media with writer
func (w responseWriter) DownloadMedia(mediaId string, writer io.Writer) error {
	return w.wx.DownloadMedia(mediaId, writer)
}
