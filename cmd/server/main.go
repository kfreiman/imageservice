package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/kfreiman/imageservice/internal/app"
	"github.com/kfreiman/imageservice/pkg/handler"
	"github.com/kfreiman/imageservice/pkg/logging"
	"github.com/kfreiman/imageservice/pkg/processor/imaging"
	"github.com/kfreiman/imageservice/pkg/repo"
	"github.com/kfreiman/imageservice/pkg/service"
)

type config struct {
	Port    int    `envconfig:"port" default:"8080"`
	IsDev   bool   `envconfig:"is_dev" default:"false"`
	RepoDir string `envconfig:"repo_dir" default:"/tmp"`
}

func main() {
	conf := parseConfig()
	logger := logging.NewLogger(conf.IsDev)
	repo, err := repo.NewRepo(conf.RepoDir)
	if err != nil {
		logger.Fatal(err)
	}

	processor := imaging.NewProcessor()
	service := service.NewService(repo, processor)
	extractor := handler.NewExtractor()
	handler := handler.NewHTTPHandler(service, extractor)
	router := app.NewRouter(handler)
	server := app.NewHTTPServer(logger, router)

	server.Start(conf.Port)
}

func parseConfig() config {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := config{}

	err := envconfig.Process("app", &c)
	if err != nil {
		panic(err)
	}
	return c
}
