package sigv4http

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

func NewSigner(cfg aws.Config, core http.RoundTripper) Signer {
	return Signer{Cfg: cfg, Core: core}
}

type Signer struct {
	Cfg  aws.Config
	Core http.RoundTripper
}

type closeSegue struct {
	reader io.Reader
	closer io.Closer
}

func (r *closeSegue) Read(p []byte) (n int, err error) { return r.reader.Read(p) }
func (r *closeSegue) Close() error                     { return r.closer.Close() }

func (i Signer) RoundTrip(req *http.Request) (*http.Response, error) {
	credentials, err := i.Cfg.Credentials.Retrieve(req.Context())
	if err != nil {
		return nil, err
	}

	var bodyBuf bytes.Buffer
	tee := io.TeeReader(req.Body, &bodyBuf)

	body, err := io.ReadAll(tee)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(body)
	hash := hex.EncodeToString(h.Sum(nil))

	req.Body = &closeSegue{reader: &bodyBuf, closer: req.Body}

	signer := v4.NewSigner()
	err = signer.SignHTTP(req.Context(), credentials, req, hash, "execute-api", i.Cfg.Region, time.Now())
	if err != nil {
		return nil, err
	}

	return i.Core.RoundTrip(req)
}
