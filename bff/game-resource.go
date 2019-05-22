package bff

import (
	"context"
	pbgameengine "github.com/emailtovamos/ms-apis/ms-game-engine/v1"
	pbhighscore "github.com/emailtovamos/ms-apis/ms-highscore/v1"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"strconv"
)

func NewGameResource(gameClient pbhighscore.GameClient, gameEngineClient pbgameengine.GameEngineClient) *gameResource {
	return &gameResource{
		gameClient:       gameClient,
		gameEngineClient: gameEngineClient,
	}
}

type gameResource struct {
	gameClient       pbhighscore.GameClient
	gameEngineClient pbgameengine.GameEngineClient
}

func (gr *gameResource) SetHighScore(c *gin.Context) {

	highScoreString := c.Param("hs")
	highScoreFloat64, _ := strconv.ParseFloat(highScoreString, 64)
	_, _ = gr.gameClient.SetHighScore(context.Background(), &pbhighscore.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})
}

func (gr *gameResource) GetHighScore(c *gin.Context) {
	if gr.gameClient == nil {
		log.Error().Msg("nil game client")
	}

	highScoreResponse, err := gr.gameClient.GetHighScore(context.Background(), &pbhighscore.GetHighScoreRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting high score")
		log.Panic()
	}

	bsString := strconv.FormatFloat(highScoreResponse.HighScore, 'e', -1, 64)
	c.JSONP(200, gin.H{
		"hs": bsString,
	})

}

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

func (gr *gameResource) SetScore(c *gin.Context) {
	if gr.gameEngineClient == nil {
		log.Error().Msg("nil gameEngine client")
	}

	scoreString := c.Param("score")
	scoreFloat64, _ := strconv.ParseFloat(scoreString, 64)

	_, err := gr.gameEngineClient.SetScore(context.Background(), &pbgameengine.SetScoreRequest{
		Score: scoreFloat64,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting high score")
		log.Panic()
	}

}

// Make NewGrpcGameEngineServiceClient and NewGrpcGameServiceClient together

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
