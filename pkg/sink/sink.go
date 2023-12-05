package sink

import (
	"connectors/pkg/entities"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// why no tests?
type Sink struct {
	clientChan    chan entities.Entity
	readingIsDone chan struct{}
	allEntities   []entities.Entity
	ownerId       string
}

func New(bufferSize uint64, ownerId string) *Sink {
	s := &Sink{
		clientChan:    make(chan entities.Entity, bufferSize),
		readingIsDone: make(chan struct{}),
		allEntities:   []entities.Entity{},
		ownerId:       ownerId,
	}

	go func() {
		defer close(s.readingIsDone)
		for e := range s.clientChan {
			e := e
			s.add(e)
		}
	}()

	return s
}

func (s *Sink) Close() {
	close(s.clientChan)
	<-s.readingIsDone
}

func (s *Sink) Push(e entities.Entity) {
	s.clientChan <- e
}

// consider using io.ReadCloser or function returning io.ReadCloser, error
// add proper test
// similar concept to https://github.com/gciezkowskiobjectivity/connectors/blob/main/pkg/idstorage/idstorage.go
func (s *Sink) Dump() ([]entities.Entity, error) {
	<-s.readingIsDone

	for _, e := range s.allEntities {
		// revers logic, no nested if
		if e.IsFile {
			err := s.downloadFile(e.ExternalId+"."+filepath.Ext(e.ContentUrl), e.ContentUrl)
			if err != nil {
				log.Printf("failed to download file %s: %s", e.ContentUrl, err.Error())
			}
		}
	}

	fileBytes, err := json.Marshal(s.allEntities)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entities: %w", err)
	}

	// close file/io.ReadCloser in defer
	file, err := os.Create(fmt.Sprintf("%s.json", s.ownerId))
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return s.allEntities, nil
}

func (s *Sink) add(e entities.Entity) {
	s.allEntities = append(s.allEntities, e)
}

func (s *Sink) downloadFile(file string, url string) error {
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Printf("Downloading file %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
