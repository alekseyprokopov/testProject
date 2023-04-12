package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"testProject/app/response"
	"testProject/app/storage"
	"testProject/configs"
)

type Server struct {
	router  *gin.Engine
	storage *storage.Storage

	httpClient *http.Client
}

func New() (*Server, error) {
	store, err := storage.New(configs.SqlPath)
	if err != nil {
		return nil, fmt.Errorf("can't create Storage: %w", err)
	}

	if err := store.Init(); err != nil {
		return nil, fmt.Errorf("can't init storage: %w", err)
	}

	return &Server{
		router:     gin.Default(),
		storage:    store,
		httpClient: &http.Client{},
	}, nil

}

func (s *Server) Start() error {
	s.storage.Init()
	s.router.POST("/saveData", s.saveDataHandler)

	return s.router.Run()
}

func (s *Server) saveDataHandler(c *gin.Context) {
	data, err := s.doRequest(configs.ApiUrl)
	if err != nil {
		log.Printf("can't do request: %v", err)
	}

	resp := response.Response{}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("can't unmarshall data: %v", err)
	}

	for _, item := range resp.Results {
		s.storage.Save(&item)
	}
	c.IndentedJSON(http.StatusOK, "items saved")
}

func (s *Server) doRequest(url string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read info from response %w", err)

	}

	return body, nil
}
