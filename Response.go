package wxhelper

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

// Thin wrapper of http.Response(can also be used as http.Response).
type Response struct {
	*http.Response
}

// Read response body into a byte slice.
func (this *Response) ReadAll() ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch this.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(this.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = this.Body
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// Read response body into string.
func (this *Response) ToString() (string, error) {
	bytes, err := this.ReadAll()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
