package bff

import (
	"context"
	pbgameengine "github.com/emailtovamos/ms-apis/ms-game-engine/v1"
	pbhighscore "github.com/emailtovamos/ms-apis/ms-highscore/v1"

	// "io/ioutil"
	"github.com/rs/zerolog/log"
	"net/http"
	// "fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"strconv"
)

const (
	playerName = "playerName"
	voteFor    = "votefor"
	highScore  = "highScore"
)

var bestScore = 999999.0

func NewGameResource(gameClient pbhighscore.GameClient, gameEngineClient pbgameengine.GameEngineClient) *gameResource {
	return &gameResource{
		gameClient:       gameClient,
		gameEngineClient: gameEngineClient,
	}
}

func NewGameResourceTemp() *gameResource {
	return &gameResource{}
}

type gameResource struct {
	gameClient       pbhighscore.GameClient
	gameEngineClient pbgameengine.GameEngineClient
}

type gameEngineResource struct {
	gameEngineClient pbgameengine.GameEngineClient
}

// func (gr *gameResource) SetHighScore(writer http.ResponseWriter, request *http.Request) {
func (gr *gameResource) SetHighScore(c *gin.Context) {

	highScoreString := c.Param("hs")
	highScoreFloat64, _ := strconv.ParseFloat(highScoreString, 64)
	_, _ = gr.gameClient.SetHighScore(context.Background(), &pbhighscore.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})
}

// func (gr *gameResource) GetHighScore(writer http.ResponseWriter, request *http.Request) {
func (gr *gameResource) GetHighScore(c *gin.Context) {
	if gr.gameClient == nil {
		log.Error().Msg("nil game client")
	}

	highScoreResponse, err := gr.gameClient.GetHighScore(context.Background(), &pbhighscore.GetHighScoreRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting high score")
		log.Panic()
	}
	// TODO: Do error check
	// _, err := gr.gameClient.GetHighScore(context.Background(), &pbhighscore.GetHighScoreRequest{})
	bsString := strconv.FormatFloat(highScoreResponse.HighScore, 'e', -1, 64)
	c.JSONP(200, gin.H{
		"hs": bsString,
	})

}

// func (gr *gameResource) GetHighScore(writer http.ResponseWriter, request *http.Request) {
func (gr *gameResource) GetSize(c *gin.Context) {
	if gr.gameEngineClient == nil {
		log.Error().Msg("nil gameEngine client")
	}

	sizeResponse, err := gr.gameEngineClient.GetSize(context.Background(), &pbgameengine.GetSizeRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting high score")
		log.Panic()
	}

	c.JSONP(200, gin.H{
		"size": sizeResponse.GetSize(),
	})

}

func (gr *gameResource) HandleGet(c *gin.Context) {
	bsString := strconv.FormatFloat(bestScore, 'e', -1, 64)
	c.JSONP(200, gin.H{
		"hs": bsString,
	})
}

// NewGrpcGameServiceClient dials grpc connection to highscore service and returns client and error
func NewGrpcGameServiceClient(serverAddr string) (pbhighscore.GameClient, error) {

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}
	if conn == nil {
		log.Info().Msg("ms-highscore conn is nil in ms-bff")
	}

	client := pbhighscore.NewGameClient(conn)

	if client == nil {
		log.Info().Msg("ms-highscore client is nil in ms-bff")
	}

	return client, nil
}

// NewGrpcGameEngineServiceClient dials grpc connection to game engine service and returns client and error
func NewGrpcGameEngineServiceClient(serverAddr string) (pbgameengine.GameEngineClient, error) {

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}
	if conn == nil {
		log.Info().Msg("ms-game-engine conn is nil in ms-bff")
	}

	client := pbgameengine.NewGameEngineClient(conn)

	if client == nil {
		log.Info().Msg("ms-game-engine client is nil in ms-bff")
	}

	return client, nil
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
