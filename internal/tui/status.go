package tui

type subModelStatus int

const (
	statusInit subModelStatus = iota
	statusMainView
	statusLoadingFormData
	statusEntry
	statusConfirmation
	statusRefreshing
)
