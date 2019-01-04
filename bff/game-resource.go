package bff

import (
	"context"
	// pb "github.com/Graphmasters/go-genproto/workshop/game_engine/v1"
	pb "github.com/teach/ms-apis/ms-highscore/v1"
	// "io/ioutil"
	"github.com/rs/zerolog/log"
	"net/http"
	// "fmt"
	"strconv"
)

const (
	playerName = "playerName"
	voteFor    = "votefor"
	highScore  = "highScore"
)

func NewGameResource(gameClient pb.GameClient) *gameResource {
	return &gameResource{
		gameClient: gameClient,
	}
}

type gameResource struct {
	gameClient pb.GameClient
}

func (gr *gameResource) SetHighScore(writer http.ResponseWriter, request *http.Request) {
	highScoreString := request.URL.Query().Get(highScore)
	highScoreFloat64, _ := strconv.ParseFloat(highScoreString, 64)
	_, err := gr.gameClient.SetHighScore(context.Background(), &pb.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})

	if err != nil {
		writer.WriteHeader(500)
		respondError(writer, err)
	} else {
		respondSuccess(writer)
	}
}

func getError(err error) []byte {
	return []byte("Oh no...Matthias screwed up or did a bad review:/n" + err.Error())
}

func respondError(writer http.ResponseWriter, err error) {
	_, err = writer.Write(getError(err))

	if err != nil {
		log.Error().Err(err).Msg("There was an error writer to client")
	}
}

func respondSuccess(writer http.ResponseWriter) {
	_, err := writer.Write([]byte("Success"))
	if err != nil {
		log.Error().Err(err).Msg("There was an error writer to client")
	}
}
