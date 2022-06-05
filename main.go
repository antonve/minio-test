package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Server struct {
	storage *minio.Client
}

func main() {
	e := echo.New()

	storageClient, err := minio.New("minio:9000", &minio.Options{
		Creds: credentials.NewStaticV4("root", "password", ""),
	})
	if err != nil {
		panic(err)
	}

	s := Server{
		storage: storageClient,
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/files", s.uploadHandler)

	e.Logger.Fatal(e.Start(":8000"))
}

func (s *Server) uploadHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	fileReader, err := file.Open()
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = s.storage.PutObject(
		c.Request().Context(),
		"videos",
		file.Filename,
		fileReader,
		file.Size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}
