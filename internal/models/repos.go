package models

type AllRepos struct {
	APIKey              APIKeyRepository
	APIUser             APIUserRepository
	Config              ConfigRepository
	TicketNotifications TicketNotificationRepository
	NotifierForwards    NotifierForwardRepository
	NotifierRules       NotifierRuleRepository
	WebexRecipients     WebexRecipientRepository
	CW                  CWRepos
	Addigy              AddigyRepos
}

type CWRepos struct {
	Board        BoardRepository
	Company      CompanyRepository
	Contact      ContactRepository
	Member       MemberRepository
	Note         TicketNoteRepository
	Ticket       TicketRepository
	TicketStatus TicketStatusRepository
}

type AddigyRepos struct {
	Alert       AddigyAlertRepository
	AlertConfig AddigyAlertConfigRepository
}
