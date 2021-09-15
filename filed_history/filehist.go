package filedb

import (
	//	"bytes"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math"

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

func Id(m iDShortURL) string {
	return m.ID
}
func URL(m iDShortURL) string {
	return m.LongURL
}

type iDShortURLInterfuck interface {
	Id() string
	URL() string
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

// interfaceDecode decodes the next interface value from the stream and returns it.
func interfaceDecodeHJNJ(dec *gob.Decoder) (iDShortURLInterfuck, error) {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p iDShortURLInterfuck
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

//s:"gob: local interface type *filedb.iDShortURLInterfuck can only be decoded from remote interface type; received concrete type iDShortURL"
func UpdateDB(rdb *repos.SetterGetter, cfg config.Config) error {
	fileName := cfg.FStorPath + "/URLs.txt"

	//var p iDShortURLInterfuck

	/*gob.Register(&p)
	//	dec := gob.NewDecoder(&network)

	reader, err := NewFReader(fileName)
	for {
		result, err := interfaceDecodeHJNJ(reader.decoder)
		if err != nil {
			break
		}
		fmt.Println(result.Id(), result.URL())
	}
	*/
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
			fmt.Println(readIDShortURL.ID, readIDShortURL.LongURL)
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
	fileName := cfg.FStorPath + "/URLs.txt"
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

type Point struct {
	X, Y int
}

func (p Point) Hypotenuse() float64 {
	return math.Hypot(float64(p.X), float64(p.Y))
}

type Pythagoras interface {
	Hypotenuse() float64
}

func scsfvsfv() {
	// This example shows how to encode an interface value. The key
	// distinction from regular types is to register the concrete type that
	// implements the interface.
	var network bytes.Buffer // Stand-in for the network.

	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Point{})

	// Create an encoder and send some values.
	enc := gob.NewEncoder(&network)
	for i := 1; i <= 3; i++ {
		interfaceEncode(enc, Point{3 * i, 4 * i})
	}

	// Create a decoder and receive some values.
	dec := gob.NewDecoder(&network)
	for i := 1; i <= 3; i++ {
		result := interfaceDecode(dec)
		fmt.Println(result.Hypotenuse())
	}

}

// interfaceEncode encodes the interface value into the encoder.
func interfaceEncode(enc *gob.Encoder, p Pythagoras) {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.

	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type. If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
	}
}

// interfaceDecode decodes the next interface value from the stream and returns it.
func interfaceDecode(dec *gob.Decoder) Pythagoras {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p Pythagoras
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return p
}
