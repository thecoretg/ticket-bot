package ticketbot

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"strconv"
	"tctg-automation/pkg/util"
)

const (
	boardSettingsTableName = "boards"
)

var (
	boardSettingsSchema = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			board_id INTEGER NOT NULL,
			board_name TEXT NOT NULL,
			webex_room_id TEXT NOT NULL,
			enabled BOOLEAN NOT NULL DEFAULT 1,
			UNIQUE (board_id)
		);`, boardSettingsTableName)
)

type boardSetting struct {
	ID          int    `db:"id" json:"-"`
	BoardID     int    `db:"board_id" json:"board_id"`
	BoardName   string `db:"board_name" json:"board_name"`
	WebexRoomID string `db:"webex_room_id" json:"webex_room_id"`
	Enabled     bool   `db:"enabled" json:"enabled"`
}

type boardsHandler struct {
	db *sqlx.DB
}

func newBoardsHandler(db *sqlx.DB) *boardsHandler {
	return &boardsHandler{
		db: db,
	}
}

func addBoardRoutes(r *gin.Engine, h *boardsHandler) {
	g := r.Group("/boards")

	g.POST("", h.createOrUpdateBoard)
	g.GET("", h.listBoards)
	g.GET("/:board_id", h.getBoard)
	// TODO: add PUT
	g.DELETE("/:board_id", h.deleteBoard)
}

func (h *boardsHandler) createOrUpdateBoard(c *gin.Context) {
	b := &boardSetting{}
	if err := c.ShouldBindJSON(b); err != nil {
		slog.Error("failed to bind JSON to boardSetting", "error", err)
		c.JSON(http.StatusBadRequest, util.ErrorJSON("invalid request body"))
		return
	}

	if err := createOrUpdateBoard(h.db, b); err != nil {
		slog.Error("failed to add or update board setting", "boardId", b.BoardID, "boardName", b.BoardName, "webexRoomID", b.WebexRoomID, "enabled", b.Enabled, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to add or update board setting"))
		return
	}

	updatedBoard, err := getBoard(h.db, b.BoardID)
	if err != nil {
		slog.Error("failed to retrieve updated board setting", "boardId", b.BoardID, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to retrieve updated board setting"))
		return
	}

	slog.Info("board setting added or updated", "boardId", b.BoardID, "boardName", b.BoardName, "webexRoomID", b.WebexRoomID, "enabled", b.Enabled)
	c.JSON(http.StatusOK, updatedBoard)
}

func (h *boardsHandler) listBoards(c *gin.Context) {
	boards, err := listBoards(h.db)
	if err != nil {
		slog.Error("failed to get boards from database", "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to get boards"))
		return
	}

	c.JSON(http.StatusOK, boards)
}

func (h *boardsHandler) getBoard(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		slog.Error("invalid board_id parameter", "board_id", boardIDStr, "error", err)
		c.JSON(http.StatusBadRequest, util.ErrorJSON("board_id must be a valid integer"))
		return
	}

	b, err := getBoard(h.db, boardID)
	if err != nil {
		slog.Error("failed to get board from database", "boardId", boardID, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to get board"))
		return
	}

	if b == nil {
		slog.Info("board not found", "boardId", boardID)
		c.JSON(http.StatusNotFound, util.ErrorJSON("board not found"))
		return
	}

	c.JSON(http.StatusOK, b)
}

func (h *boardsHandler) deleteBoard(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		slog.Error("invalid board_id parameter", "board_id", boardIDStr, "error", err)
		c.JSON(http.StatusBadRequest, util.ErrorJSON("board_id must be a valid integer"))
		return
	}

	// TODO: check presence, return 404 if not found
	if err := deleteBoard(h.db, boardID); err != nil {
		slog.Error("failed to delete board setting", "boardId", boardID, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("internal server error"))
		return
	}

	slog.Info("board setting deleted", "boardId", boardID)
	c.Status(http.StatusNoContent)
}

func createOrUpdateBoard(db *sqlx.DB, board *boardSetting) error {
	_, err := db.NamedExec(`
			INSERT INTO boards (board_id, board_name, webex_room_id, enabled)
		   	VALUES(:board_id, :board_name, :webex_room_id, :enabled)
		   	ON CONFLICT(board_id) DO UPDATE SET
				board_name = excluded.board_name,
				webex_room_id = excluded.webex_room_id,
				enabled = excluded.enabled
			`, board,
	)

	if err != nil {
		return err
	}

	return nil
}

func listBoards(db *sqlx.DB) ([]boardSetting, error) {
	var boards []boardSetting
	err := db.Select(&boards, "SELECT * FROM boards")
	return boards, err
}

func getBoard(db *sqlx.DB, boardID int) (*boardSetting, error) {
	var board boardSetting
	err := db.Get(&board, "SELECT * FROM boards WHERE board_id = ?", boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting board by id: %w", err)
	}
	return &board, nil
}

func deleteBoard(db *sqlx.DB, boardID int) error {
	_, err := db.Exec("DELETE FROM boards WHERE board_id = ?", boardID)
	return err
}
