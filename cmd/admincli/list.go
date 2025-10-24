package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/thecoretg/ticketbot/internal/db"
)

var (
	listCmd = &cobra.Command{
		Use: "list",
	}

	listWebexCmd = &cobra.Command{
		Use: "webex",
	}

	roomType                 string
	showRoomType, showRoomID bool
	listWebexRoomsCmd        = &cobra.Command{
		Use: "rooms",
		RunE: func(cmd *cobra.Command, args []string) error {
			rooms, err := client.ListRooms(nil)
			if err != nil {
				return fmt.Errorf("listing rooms: %w", err)
			}

			rooms, err = filterWebexRooms(rooms, roomType)
			if err != nil {
				return fmt.Errorf("filtering webex rooms: %w", err)
			}

			fmt.Println(webexRoomsTable(rooms))
			return nil
		},
	}

	listCWCmd = &cobra.Command{
		Use: "cw",
	}

	listCWBoardsCmd = &cobra.Command{
		Use: "boards",
		RunE: func(cmd *cobra.Command, args []string) error {
			boards, err := client.ListBoards(nil)
			if err != nil {
				return fmt.Errorf("listing boards: %w", err)
			}

			fmt.Println(cwBoardsTable(boards))
			return nil
		},
	}
)

func addListCmds() {
	rootCmd.AddCommand(listCmd)

	listCmd.AddCommand(listWebexCmd)
	listWebexCmd.AddCommand(listWebexRoomsCmd)

	listCmd.AddCommand(listCWCmd)
	listCWCmd.AddCommand(listCWBoardsCmd)

	listWebexRoomsCmd.Flags().StringVar(&roomType, "type", "all", "type of room to filter by: 'group', 'direct', or 'all'")
	listWebexRoomsCmd.Flags().BoolVarP(&showRoomID, "show-id", "i", true, "show room id in table")
	listWebexRoomsCmd.Flags().BoolVarP(&showRoomType, "show-type", "t", false, "show room type in table")
}

func filterWebexRooms(rooms []db.WebexRoom, roomType string) ([]db.WebexRoom, error) {
	if !validRoomType(roomType) {
		return nil, fmt.Errorf("room type '%s' not valid, expected 'group' or 'direct'", roomType)
	}

	if roomType == "all" {
		return rooms, nil
	}

	var filtered []db.WebexRoom
	for _, r := range rooms {
		if r.Type == roomType {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

func validRoomType(t string) bool {
	return t == "group" || t == "direct" || t == "all"
}

func cwBoardsTable(boards []db.CwBoard) string {
	if boards == nil || len(boards) == 0 {
		return "No boards found"
	}

	headers := []string{"ID", "NAME", "NOTIFY"}
	sort.Slice(boards, func(i, j int) bool {
		return boards[i].ID < boards[j].ID
	})

	t := table.New().
		Headers(headers...).
		StyleFunc(spacingStyleFunc())

	for _, b := range boards {
		if !b.Deleted {
			addBoardRow(t, b)
		}
	}

	return t.String()
}

func addBoardRow(t *table.Table, b db.CwBoard) {
	t.Row(strconv.Itoa(b.ID), b.Name)
}

func webexRoomsTable(rooms []db.WebexRoom) string {
	if rooms == nil || len(rooms) == 0 {
		return "No rooms found"
	}

	headers := []string{"TITLE"}

	if showRoomType {
		headers = append(headers, "TYPE")
	}

	if showRoomID {
		headers = append(headers, "ID")
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Name < rooms[j].Name
	})

	t := table.New().
		Headers(headers...).
		StyleFunc(spacingStyleFunc())

	for _, r := range rooms {
		// terminated users still show but with empty title, don't show them
		if r.Name != "Empty Name" {
			addRoomRow(t, r, showRoomType, showRoomID)
		}
	}

	return t.String()
}

// lowkey this sounds like something scooby doo would say
func addRoomRow(t *table.Table, room db.WebexRoom, showType, showID bool) {
	row := []string{room.Name}
	if showType {
		row = append(row, room.Type)
	}

	if showID {
		row = append(row, strconv.Itoa(room.ID))
	}

	t.Row(row...)
}
