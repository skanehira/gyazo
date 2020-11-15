package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var uploadEndpoint = "https://upload.gyazo.com/api/upload"

type Meta struct {
	Title    string
	Desc     string
	Filename string
}

type Image struct {
	ID           string `json:"image_id"`
	PermalinkURL string `json:"permalink_url"`
	ThumbURL     string `json:"thumb_url"`
	URL          string `json:"url"`
	Type         string `json:"type"`
	Star         bool   `json:"star"`
	CreatedAt    string `json:"created_at"`
}

type Client struct {
	*http.Client
}

func NewClient(token string) (*Client, error) {
	oauthClient := oauth2.NewClient(
		oauth2.NoContext,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	)

	c := &Client{
		Client: oauthClient,
	}

	return c, nil
}

func (c *Client) Upload(meta Meta, r io.Reader) (*Image, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filename := time.Now().Format("20060102150405")
	if meta.Filename != "" {
		filename = meta.Filename
	}

	part, err := writer.CreateFormFile("imagedata", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create a form data: %w", err)
	}
	if _, err := io.Copy(part, r); err != nil {
		return nil, err
	}

	if meta.Title != "" {
		w, err := writer.CreateFormField("title")
		if err != nil {
			return nil, err
		}

		if _, err := w.Write([]byte(meta.Title)); err != nil {
			return nil, err
		}
	}

	if meta.Desc != "" {
		w, err := writer.CreateFormField("desc")
		if err != nil {
			return nil, err
		}

		if _, err := w.Write([]byte(meta.Desc)); err != nil {
			return nil, err
		}

	}

	contentType := writer.FormDataContentType()

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close a multipart writer: %w", err)
	}

	// Be aware that the URL is different from the other API.
	req, err := http.NewRequest("POST", uploadEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %w", err)
	}
	req.Header.Add("Content-Type", contentType)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newResponseError(res)
	}

	img := &Image{}
	if err = json.NewDecoder(res.Body).Decode(img); err != nil {
		return nil, fmt.Errorf("failed to decode a responsed JSON: %w", err)
	}

	return img, nil
}

func newResponseError(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	return errors.New(string(b))
}
