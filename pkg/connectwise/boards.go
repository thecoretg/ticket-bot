package connectwise

import (
	"context"
	"fmt"
)

func boardIdEndpoint(boardId int) string {
	return fmt.Sprintf("service/boards/%d", boardId)
}

func boardIdStatusEndpoint(boardId int) string {
	return fmt.Sprintf("%s/statuses", boardIdEndpoint(boardId))
}

func boardIdStatusIdEndpoint(boardId, statusId int) string {
	return fmt.Sprintf("%s/%d", boardIdStatusEndpoint(boardId), statusId)
}

func (c *Client) ListBoards(ctx context.Context, params *QueryParams) ([]Board, error) {
	return ApiRequestPaginated[Board](ctx, c, "GET", "service/boards", params, nil)
}

func (c *Client) PostBoard(ctx context.Context, board *Board) (*Board, error) {
	return ApiRequestNonPaginated[Board](ctx, c, "POST", "service/boards", nil, board)
}

func (c *Client) GetBoard(ctx context.Context, boardId int, params *QueryParams) (*Board, error) {
	return ApiRequestNonPaginated[Board](ctx, c, "GET", boardIdEndpoint(boardId), params, nil)
}

func (c *Client) DeleteBoard(ctx context.Context, boardId int) error {
	if _, err := ApiRequestNonPaginated[struct{}](ctx, c, "DELETE", boardIdEndpoint(boardId), nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) PutBoard(ctx context.Context, boardId int, board *Board) (*Board, error) {
	return ApiRequestNonPaginated[Board](ctx, c, "PUT", boardIdEndpoint(boardId), nil, board)
}

func (c *Client) PatchBoard(ctx context.Context, boardId int, patchOps []PatchOp) (*Board, error) {
	return ApiRequestNonPaginated[Board](ctx, c, "PATCH", boardIdEndpoint(boardId), nil, patchOps)
}

func (c *Client) ListBoardStatuses(ctx context.Context, boardId int, params *QueryParams) ([]BoardStatus, error) {
	return ApiRequestPaginated[BoardStatus](ctx, c, "GET", boardIdStatusEndpoint(boardId), params, nil)
}

func (c *Client) PostBoardStatus(ctx context.Context, boardId int, status *BoardStatus) (*BoardStatus, error) {
	return ApiRequestNonPaginated[BoardStatus](ctx, c, "POST", boardIdStatusEndpoint(boardId), nil, status)
}

func (c *Client) GetBoardStatus(ctx context.Context, boardId, statusId int, params *QueryParams) (*BoardStatus, error) {
	return ApiRequestNonPaginated[BoardStatus](ctx, c, "GET", boardIdStatusIdEndpoint(boardId, statusId), params, nil)
}

func (c *Client) DeleteBoardStatus(ctx context.Context, boardId, statusId int) error {
	if _, err := ApiRequestNonPaginated[struct{}](ctx, c, "DELETE", boardIdStatusIdEndpoint(boardId, statusId), nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) PutBoardStatus(ctx context.Context, boardId, statusId int, status *BoardStatus) (*BoardStatus, error) {
	return ApiRequestNonPaginated[BoardStatus](ctx, c, "PUT", boardIdStatusIdEndpoint(boardId, statusId), nil, status)
}

func (c *Client) PatchBoardStatus(ctx context.Context, boardId, statusId int, patchOps []PatchOp) (*BoardStatus, error) {
	return ApiRequestNonPaginated[BoardStatus](ctx, c, "PATCH", boardIdStatusIdEndpoint(boardId, statusId), nil, patchOps)
}
