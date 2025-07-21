package connectwise

import (
	"fmt"
)

func memberIdEndpoint(memberId int) string {
	return fmt.Sprintf("system/members/%d", memberId)
}

func (c *Client) PostMember(member *Member) (*Member, error) {
	return Post[Member](c, "system/members", member)
}

func (c *Client) ListMembers(params map[string]string) ([]Member, error) {
	return GetMany[Member](c, "system/members", params)
}

func (c *Client) GetMember(memberID int, params map[string]string) (*Member, error) {
	return GetOne[Member](c, memberIdEndpoint(memberID), params)
}

func (c *Client) PutMember(memberID int, member *Member) (*Member, error) {
	return Put[Member](c, memberIdEndpoint(memberID), member)
}

func (c *Client) PatchMember(memberID int, patchOps []PatchOp) (*Member, error) {
	return Patch[Member](c, memberIdEndpoint(memberID), patchOps)
}

func (c *Client) DeleteMember(memberID int) error {
	return Delete(c, memberIdEndpoint(memberID))
}
