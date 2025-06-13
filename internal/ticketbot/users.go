package ticketbot

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"tctg-automation/pkg/util"
)

const (
	usersTableName = "users"
)

var (
	usersTableSchema = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			cw_id TEXT NOT NULL,
			email TEXT NOT NULL,
			mute BOOLEAN NOT NULL DEFAULT 0,
			ignore_update BOOLEAN NOT NULL DEFAULT 0,
			UNIQUE (cw_id, email)
		);`, usersTableName)
)

type user struct {
	ID           int    `db:"id" json:"-"`
	CWId         string `db:"cw_id" json:"cw_id"`
	Email        string `db:"email" json:"email"`
	Mute         bool   `db:"mute" json:"mute"`
	IgnoreUpdate bool   `db:"ignore_update" json:"ignore_update"`
}

type usersHandler struct {
	db *sqlx.DB
}

func newUsersHandler(db *sqlx.DB) *usersHandler {
	return &usersHandler{
		db: db,
	}
}

func addUserRoutes(r *gin.Engine, h *usersHandler) {
	g := r.Group("/users")

	g.POST("", h.createOrUpdateUser)
	g.GET("", h.listUsers)
	g.GET("/:user_id", h.getUser)
	// TODO: add PUT
	g.DELETE("/:user_id", h.deleteUser)
}

func (h *usersHandler) createOrUpdateUser(c *gin.Context) {
	u := &user{}
	if err := c.ShouldBindJSON(u); err != nil {
		slog.Error("failed to bind JSON to user", "error", err)
		c.JSON(http.StatusBadRequest, util.ErrorJSON("invalid request body"))
		return
	}

	if err := createOrUpdateUser(h.db, u); err != nil {
		slog.Error("failed to add or update user", "cwId", u.CWId, "email", u.Email, "mute", u.Mute, "ignoreUpdate", u.IgnoreUpdate, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to add or update user"))
		return
	}

	updatedUser, err := getUser(h.db, u.CWId)
	if err != nil {
		slog.Error("failed to retrieve updated user", "userID", u.CWId, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to retrieve created or updated user after creation"))
		return
	}

	slog.Info("user added or updated", "cwId", updatedUser.CWId, "email", updatedUser.Email, "mute", updatedUser.Mute, "ignoreUpdate", updatedUser.IgnoreUpdate)
	c.JSON(http.StatusOK, updatedUser)
}

func (h *usersHandler) listUsers(c *gin.Context) {
	users, err := listUsers(h.db)
	if err != nil {
		slog.Error("failed to get users from database", "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("failed to get users"))
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *usersHandler) getUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		slog.Error("user ID is required")
		c.JSON(http.StatusBadRequest, util.ErrorJSON("user ID is required"))
		return
	}

	u, err := getUser(h.db, userID)
	if err != nil {
		slog.Error("failed to get user from database", "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("failed to get user: %v", err)))
		return
	}

	if u == nil {
		slog.Info("user not found", "userID", userID)
		c.JSON(http.StatusNotFound, util.ErrorJSON("user not found"))
		return
	}

	slog.Info("user retrieved", "userID", u.CWId, "email", u.Email, "mute", u.Mute, "ignoreUpdate", u.IgnoreUpdate)
	c.JSON(http.StatusOK, u)
}

func (h *usersHandler) deleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		slog.Error("user ID is required")
		c.JSON(http.StatusBadRequest, util.ErrorJSON("user ID is required"))
		return
	}

	// TODO: check presence, return 404 if not found
	if err := deleteUser(h.db, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("user not found for deletion", "userID", userID)
			c.JSON(http.StatusNotFound, util.ErrorJSON("user not found"))
			return
		}
		slog.Error("failed to delete user from database", "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("failed to delete user: %v", err)))
		return
	}

	slog.Info("user deleted", "userID", userID)
	c.Status(http.StatusNoContent)
}

func createOrUpdateUser(db *sqlx.DB, user *user) error {
	_, err := db.NamedExec(`
			INSERT INTO users (cw_id, email, mute, ignore_update)
			VALUES(:cw_id, :email, :mute, :ignore_update)
			ON CONFLICT(cw_id, email) DO UPDATE SET
			    cw_id = excluded.cw_id,
			    email = excluded.email,
			    mute = excluded.mute,
			    ignore_update = excluded.ignore_update
			`, user,
	)
	if err != nil {
		return fmt.Errorf("adding or updating user: %w", err)
	}

	return nil
}

func listUsers(db *sqlx.DB) ([]user, error) {
	var users []user
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}

func getUser(db *sqlx.DB, cwId string) (*user, error) {
	var user user
	err := db.Get(&user, "SELECT * FROM users WHERE cw_id = ?", cwId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No user found with the given cw_id
		}
		return nil, fmt.Errorf("getting user by cw_id: %w", err)
	}
	return &user, nil
}

func deleteUser(db *sqlx.DB, userCwId string) error {
	_, err := db.Exec("DELETE FROM users WHERE cw_id = ?", userCwId)
	return err
}
