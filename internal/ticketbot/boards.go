package ticketbot

import "fmt"

func (s *server) ensureBoardInStore(cwData *cwData) (*Board, error) {
	board, err := s.dataStore.GetBoard(cwData.ticket.Board.ID)
	if err != nil {
		return nil, fmt.Errorf("getting board from storage: %w", err)
	}

	if board == nil {
		board, err = s.addBoard(cwData.ticket.Board.ID)
		if err != nil {
			return nil, fmt.Errorf("inserting board into store: %w", err)
		}
	}

	return board, nil
}

// addBoard adds Connecwise boards to the data store, with a default of
// notifications not enabled.
func (s *server) addBoard(boardID int) (*Board, error) {
	cwBoard, err := s.cwClient.GetBoard(boardID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting board from connectwise: %w", err)
	}

	storeBoard := &Board{
		ID:            cwBoard.ID,
		Name:          cwBoard.Name,
		NotifyEnabled: false,
	}

	if err := s.dataStore.UpsertBoard(storeBoard); err != nil {
		return nil, fmt.Errorf("adding board to store: %w", err)
	}

	return storeBoard, nil
}
