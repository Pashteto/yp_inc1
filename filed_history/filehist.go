package filedb

import (
<<<<<<< HEAD
	"context"
	"encoding/gob"
=======
	//	"bytes"

	"context"
	"encoding/gob"
	"strings"

	//	"encoding/json"
>>>>>>> main
	"errors"
	"log"
	"net/url"
	"os"
<<<<<<< HEAD
	"sort"
	"strings"
=======
>>>>>>> main
	"time"

	"github.com/Pashteto/yp_inc1/config"
	"github.com/Pashteto/yp_inc1/repos"
)

var ctx, _ = context.WithCancel(context.Background())

const urlTTL1 = time.Second * 1000

<<<<<<< HEAD
type iDLongURLPair struct {
=======
type iDShortURL struct {
>>>>>>> main
	ID      string
	LongURL string
}

<<<<<<< HEAD
func URL(m iDLongURLPair) string {
	return m.LongURL
}

type userAndPairs struct {
	User       string
	PairsSlice []iDLongURLPair
}

func User(m userAndPairs) string {
	return m.User
}

type FWriter interface {
	WriteUserAndPairs(userPairs []userAndPairs) error
	Close() error
}
type FReader interface {
	ReadUserAndPairs() ([]userAndPairs, error)
=======
func ID(m iDShortURL) string {
	return m.ID
}
func URL(m iDShortURL) string {
	return m.LongURL
}

type FWriter interface {
	WriteIDShortURL(idShURL []iDShortURL) error
	Close() error
}
type FReader interface {
	ReadIDShortURL() ([]iDShortURL, error)
>>>>>>> main
	Close() error
}

type fWriter struct {
	file    *os.File
	encoder *gob.Encoder
}

type fReader struct {
	file    *os.File
	decoder *gob.Decoder
}

func NewFWriter(fileName string) (*fWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
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
<<<<<<< HEAD
=======

>>>>>>> main
	return &fReader{
		file:    file,
		decoder: gob.NewDecoder(file),
	}, nil
}

<<<<<<< HEAD
func (p *fWriter) WriteUserAndPairs(userPairs *userAndPairs) error {
	return p.encoder.Encode(userPairs)
}

func (c *fReader) ReadUserAndPairs() ([]userAndPairs, error) {
	userPairs := []userAndPairs{}
	if err := c.decoder.Decode(&userPairs); err != nil {
=======
func (p *fWriter) WriteIDShortURL(idShURL *iDShortURL) error {
	return p.encoder.Encode(idShURL)
}

func (c *fReader) ReadIDShortURL() ([]iDShortURL, error) {
	idShURL := []iDShortURL{}
	if err := c.decoder.Decode(&idShURL); err != nil {
>>>>>>> main
		if err.Error() != "EOF" {
			return nil, err
		}
	}
<<<<<<< HEAD
	return userPairs, nil
=======
	return idShURL, nil
>>>>>>> main
}

func (p *fWriter) Close() error {
	return p.file.Close()
}

func (c *fReader) Close() error {
	return c.file.Close()
}
func CreateDirFileDBExists(cfg config.Config) error {
	fjnv := strings.SplitAfter(cfg.FStorPath, "/")
	if len(fjnv) > 0 {
		fjnv = fjnv[:len(fjnv)-1]
		fjnv1 := strings.Join(fjnv, "")
		return os.MkdirAll(fjnv1, 0777)
	}
	return nil
}

<<<<<<< HEAD
////
/*
// Set attaches the redis repository and set the data
func (r *repository) Set(ctx context.Context, key string, value UserAndString, exp time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, p, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *repository) Get(ctx context.Context, key string) (UserAndString, error) {
	var res UserAndString
	p, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return UserAndString{"", ""}, err
	}
	err = json.Unmarshal([]byte(p), &res)
	return res, err
}
*/
///

func UpdateDBSlice(rdb repos.SetterGetter, cfg config.Config) ([]string, error) {
=======
func UpdateDBSlice(rdb repos.SetterGetter, cfg config.Config) error {
>>>>>>> main
	fileName := cfg.FStorPath
	reader, err := NewFReader(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()
<<<<<<< HEAD
	readUserPairsSlice, err := reader.ReadUserAndPairs()
	if err != nil {
		return nil, err
	}
	errReadDB := pingRedisDB(rdb)
	if errReadDB != nil {
		return nil, errReadDB
	}
	err = rdb.FlushAllKeys(ctx)
	if err != nil {
		return nil, err
	}
	newUserList := make([]string, len(readUserPairsSlice))
	i := 0
	for _, userPairsSlice := range readUserPairsSlice {
		pairs := userPairsSlice.PairsSlice
		user := userPairsSlice.User
		newUserList[i] = user
		i++
		for _, pair := range pairs {
			err = testFiledURLAndConvert(&pair)
			if err != nil {
				res := deleteEmpty(newUserList)
				sort.Strings(res)
				return res, err
			}
			key := pair.ID
			value := pair.LongURL
			err = rdb.SetValueByKeyAndUser(ctx, key, user, value, urlTTL1)
			if err != nil {
				res := deleteEmpty(newUserList)
				sort.Strings(res)
				return res, err
			}
		}
	}
	sort.Strings(newUserList)
	return newUserList, nil
}

func testFiledURLAndConvert(in *iDLongURLPair) error {
=======
	readIDShortURLSlice, err := reader.ReadIDShortURL()
	if err != nil {
		return err
	}
	errReadDB := pingRedisDB(rdb)
	if errReadDB != nil {
		return errReadDB
	}
	err = rdb.FlushAllKeys(ctx)
	if err != nil {
		return err
	}
	for i := range readIDShortURLSlice {
		strIDShortURL := readIDShortURLSlice[i]
		err = testFiledURLAndConvert(&strIDShortURL)
		if err != nil {
			return err
		}
		key := strIDShortURL.ID
		value := strIDShortURL.LongURL
		err = rdb.Set(ctx, key, value, urlTTL1)
		if err != nil {
			return err
		}
	}
	return nil
}

func testFiledURLAndConvert(in *iDShortURL) error {
>>>>>>> main
	if in == nil {
		return errors.New("nil filed id")
	}
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
<<<<<<< HEAD
	fileName := cfg.FStorPath
=======
	fileName := cfg.FStorPath // + "/URLs.txt"
>>>>>>> main
	writer, err := NewFWriter(fileName)
	if err != nil {
		return err
	}
	defer writer.Close()
<<<<<<< HEAD
	idShURL := &iDLongURLPair{ID: id, LongURL: longURL.String()}
=======
	idShURL := &iDShortURL{ID: id, LongURL: longURL.String()}
>>>>>>> main
	if err := writer.encoder.Encode(&idShURL); err != nil {
		return err
	}
	return nil
}

<<<<<<< HEAD
func WriteAll(rdb repos.SetterGetter, cfg config.Config, UsersList *[]string) error {
=======
func WriteAll(rdb repos.SetterGetter, cfg config.Config) error {
>>>>>>> main
	fileName := cfg.FStorPath

	writer, err := NewFWriter(fileName)
	if err != nil {
		return err
	}
	defer writer.Close()

<<<<<<< HEAD
	var DBWrite []userAndPairs

	for _, user := range *UsersList {
		pairsFromDB, err := rdb.ListAllKeysByUser(ctx, user)
		if err != nil {
			return err
		}
		pairs := make([]iDLongURLPair, len(pairsFromDB))
		i := 0
		for key, value := range pairsFromDB {
			pairs[i] = iDLongURLPair{ID: key, LongURL: value}
			i++
		}

		DBWrite = append(DBWrite, userAndPairs{User: user, PairsSlice: pairs})
=======
	var DBWrite []iDShortURL
	keys, err := rdb.ListAllKeys(ctx)
	if err != nil {
		return err
	}
	for _, key := range keys {
		value, err := rdb.Get(ctx, key)
		if err != nil {
			return err
		}
		DBWrite = append(DBWrite, iDShortURL{ID: key, LongURL: value})
>>>>>>> main
	}
	if err := writer.encoder.Encode(&DBWrite); err != nil {
		return err
	}

	return nil
}

func pingRedisDB(client repos.SetterGetter) error {
	if client == nil {
		return errors.New("no redis db")
	}
	err := client.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
<<<<<<< HEAD

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
=======
>>>>>>> main
