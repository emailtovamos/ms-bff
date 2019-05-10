package main

import (
	"flag"
	"github.com/emailtovamos/ms-bff/bff"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var router *gin.Engine
var bestScore = 999999.0

func main() {
	grpcAddrHighScore := flag.String("address-ms-highscore", "localhost:50051", "The GRPC server address for highscore service")
	grpcAddrGameEngine := flag.String("address-ms-game-engine", "localhost:60051", "The GRPC server address for game engine service")

	serverAddr := flag.String("address-http", ":8081", "The HTTP server address")

	flag.Parse()

	gameClient, err := bff.NewGrpcGameServiceClient(*grpcAddrHighScore)
	if err != nil {
		log.Error().Err(err).Msg("Error creating a client for ms-highscore")
	}
	gameEngineClient, err := bff.NewGrpcGameEngineServiceClient(*grpcAddrGameEngine)
	if err != nil {
		log.Error().Err(err).Msg("Error creating a client for ms-highscore")
	}
	gr := bff.NewGameResource(gameClient, gameEngineClient)

	router = gin.Default()
	router.GET("/getbs", gr.GetHighScore)
	router.GET("/setbs/:hs", gr.SetHighScore)
	router.GET("/getsize", gr.GetSize)
	router.GET("/setscore/:score", gr.SetScore)
	err = router.Run(*serverAddr)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not init server")
	}

	log.Info().Msgf("Started HTTP-Service at [%s]", *serverAddr)
}
