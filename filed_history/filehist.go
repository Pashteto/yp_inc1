package filedb

import (
	//	"bytes"
	"context"
	"encoding/gob"

	//	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Pashteto/yp_inc1/config"
	"github.com/Pashteto/yp_inc1/repos"
)

var ctx, _ = context.WithCancel(context.Background())

type iDShortURL struct {
	ID      string
	LongURL string
}

type FWriter interface {
	WriteIDShortURL(idShURL *iDShortURL)
	Close() error
}

type FReader interface {
	ReadIDShortURL() (*iDShortURL, error)
	Close() error
}

type fWriter struct {
	file    *os.File
	encoder *gob.Encoder
}

/*
   var buffer bytes.Buffer
   gobEncoder := gob.NewEncoder(&buffer)
   gobDecoder := gob.NewDecoder(&buffer)

*/
type fReader struct {
	file    *os.File
	decoder *gob.Decoder
}

func NewFWriter(fileName string) (*fWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &fWriter{
		file:    file,
		encoder: gob.NewEncoder(file),
	}, nil
}

func NewFReader(fileName string) (*fReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &fReader{
		file:    file,
		decoder: gob.NewDecoder(file),
	}, nil
}

func (p *fWriter) WriteIDShortURL(idShURL *iDShortURL) error {
	return p.encoder.Encode(idShURL)
}

func (c *fReader) ReadIDShortURL() (*iDShortURL, error) {
	idShURL := &iDShortURL{}
	if err := c.decoder.Decode(&idShURL); err != nil {
		return nil, err
	}
	return idShURL, nil
}

func (p *fWriter) Close() error {
	return p.file.Close()
}

func (c *fReader) Close() error {
	return c.file.Close()
}
func CreateDirFileDBExists(cfg config.Config) error {
	return os.MkdirAll(cfg.FStorPath, 0777)
}

func UpdateDB(rdb *repos.SetterGetter, cfg config.Config) error {
	fileName := cfg.FStorPath + "/URLs.bin"
	reader, err := NewFReader(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for {
		readIDShortURL, err := reader.ReadIDShortURL()
		if err != nil {
			break
		}
		if testFiledURLAndConvert(readIDShortURL) == nil {
			(*rdb).Set(ctx, readIDShortURL.ID, readIDShortURL.LongURL, 1000*time.Second)
		}
	}
	return nil
}

func testFiledURLAndConvert(in *iDShortURL) error {
	if in.ID == "" {
		return errors.New("empty filed id")
	}
	if in.LongURL == "" {
		return errors.New("empty filed url")
	}
	longURL, err := url.Parse(in.LongURL)
	if err != nil {
		return errors.New("unable to parse filed url")
	}
	if !longURL.IsAbs() {
		longURL.Scheme = "http"
	}
	in.LongURL = longURL.String()
	return nil
}

func PostInFileDB(id string, longURL *url.URL, cfg config.Config) error {
	fileName := cfg.FStorPath + "/URLs.bin"
	writer, err := NewFWriter(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()
	idShURL := &iDShortURL{ID: id, LongURL: longURL.String()}
	if err := writer.encoder.Encode(&idShURL); err != nil {
		return err
	}
	return nil
}
