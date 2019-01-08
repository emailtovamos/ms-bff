package bff

import (
	"context"
	// pb "github.com/Graphmasters/go-genproto/workshop/game_engine/v1"
	pb "github.com/teach/ms-apis/ms-highscore/v1"
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

func NewGameResource(gameClient pb.GameClient) *gameResource {
	return &gameResource{
		gameClient: gameClient,
	}
}

func NewGameResourceTemp() *gameResource {
	return &gameResource{}
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

// func (gr *gameResource) GetHighScore(writer http.ResponseWriter, request *http.Request) {
func (gr *gameResource) GetHighScore(c *gin.Context) {

	highScoreResponse, _ := gr.gameClient.GetHighScore(context.Background(), &pb.GetHighScoreRequest{})
	// TODO: Do error check
	// _, err := gr.gameClient.GetHighScore(context.Background(), &pb.GetHighScoreRequest{})
	bsString := strconv.FormatFloat(highScoreResponse.HighScore, 'e', -1, 64)
	c.JSONP(200, gin.H{
		"hs": bsString,
	})
	// if err != nil {
	// 	writer.WriteHeader(500)
	// 	respondError(writer, err)
	// } else {
	// 	respondSuccess(writer)
	// }
}

func (gr *gameResource) HandleGet(c *gin.Context) {
	bsString := strconv.FormatFloat(bestScore, 'e', -1, 64)
	c.JSONP(200, gin.H{
		"hs": bsString,
	})
}

// NewGrpcGameServiceClient dials grpc connection and returns client and error
func NewGrpcGameServiceClient(serverAddr string) (pb.GameClient, error) {
	// tracers will default to a NOOP tracer if nothing was configured
	// streamTracingInterceptor := grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))
	// unaryTracingInterceptor := grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))

	// serverOpts := []grpc.DialOption{
	// 	grpc.WithStreamInterceptor(streamTracingInterceptor),
	// 	grpc.WithUnaryInterceptor(unaryTracingInterceptor),
	// 	grpc.WithInsecure()}

	// conn, err := grpc.Dial(serverAddr, serverOpts...)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	client := pb.NewGameClient(conn)

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
