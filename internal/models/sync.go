package models

type SyncPayload struct {
	WebexRecipients    bool  `json:"webex_rooms"`
	CWBoards           bool  `json:"cw_boards"`
	CWTickets          bool  `json:"cw_tickets"`
	BoardIDs           []int `json:"board_ids"`
	MaxConcurrentSyncs int   `json:"max_concurrent_syncs"`
}
