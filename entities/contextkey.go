package entities

type contextKey string

const (
	LoggerKey = contextKey("logger")
	TicketKey = contextKey("ticket_id")
)
