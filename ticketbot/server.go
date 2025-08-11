package ticketbot

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/connectwise"
	"github.com/thecoretg/ticketbot/db"
	"github.com/thecoretg/ticketbot/webex"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config      *Cfg
	queries     *db.Queries
	cwClient    *connectwise.Client
	cwCompanyID string
	webexClient *webex.Client
	ticketLocks sync.Map
	GinEngine   *gin.Engine
}

func (s *Server) Run() error {
	if !s.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return s.GinEngine.Run()
}

func (s *Server) addAllRoutes() {
	s.addHooksGroup()
	s.addBoardsGroup()
}

func NewServer(ctx context.Context, cfg *Cfg, initHooks bool) (*Server, error) {
	dbConn, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	cwCreds := &connectwise.Creds{
		PublicKey:  cfg.CwPubKey,
		PrivateKey: cfg.CwPrivKey,
		ClientId:   cfg.CwClientID,
		CompanyId:  cfg.CwCompanyID,
	}

	s := &Server{
		config:      cfg,
		cwClient:    connectwise.NewClient(cwCreds),
		webexClient: webex.NewClient(cfg.WebexSecret),
		queries:     db.New(dbConn),
		GinEngine:   gin.Default(),
	}

	if initHooks {
		if err := s.InitAllHooks(); err != nil {
			return nil, fmt.Errorf("initiating webhooks: %w", err)
		}
	}

	s.addAllRoutes()

	return s, nil
}
