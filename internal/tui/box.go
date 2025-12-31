package tui

import (
	"github.com/dsrosen6/teabox/flexbox"
	"github.com/dsrosen6/teabox/titlebox"
)

func (m *Model) mainFlexbox() *flexbox.Box {
	return flexbox.New(flexbox.Vertical, 1).
		AddFlexBox(m.topBox(), "topBox", 1, nil, nil, nil)
}

func (m *Model) topBox() *flexbox.Box {
	return flexbox.New(flexbox.Horizontal, 1).
		AddTitleBox(m.rulesBox(), "rules", 1, nil, nil, nil)
}

func (m *Model) rulesBox() titlebox.Box {
	return titlebox.New().
		SetTitle("rules")
}
