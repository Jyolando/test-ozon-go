package entities

import (
	"github.com/gorilla/mux"
	pb "github.com/jyolando/test-ozon-go/pkg/api/proto"
)

type Server struct {
	pb.UnimplementedURLShortenerServer
	router *mux.Router
}
