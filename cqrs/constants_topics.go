package cqrs

import "fmt"

const (
	Ns   = "atani"
	cqrs = "cqrs"
)

const (
	KindEvents   = "events"
	KindCommands = "commands"
)

var (
	TopicCommandsBroker = TopicCommandsAtani("broker")
	TopicEventsBroker   = TopicEventsAtani("broker")

	TopicCommandsKYC = TopicCommandsAtani("kyc")
	TopicEventsKYC   = TopicEventsAtani("kyc")

	TopicEventsAuth   = TopicEventsAtani("auth")
	TopicCommandsAuth = TopicCommandsAtani("auth")

	TopicEventsAlerts   = TopicEventsAtani("alerts")
	TopicCommandsAlerts = TopicCommandsAtani("alerts")

	TopicEventsNotifica   = TopicEventsAtani("notifica")
	TopicCommandsNotifica = TopicCommandsAtani("notifica")

	TopicEventsRebalances = TopicEventsAtani("rebalances")

	TopicEventsTracker   = TopicEventsAtani("tracker")
	TopicCommandsTracker = TopicCommandsAtani("tracker")

	TopicErrors = fmt.Sprintf("%s.%s.__errors", Ns, cqrs)
)

func TopicEventsAtani(domain string) string {
	return TopicEvents(Ns, domain)
}

func TopicCommandsAtani(domain string) string {
	return TopicCommands(Ns, domain)
}

func TopicEvents(ns string, domain string) string {
	return Topic(ns, domain, KindEvents)
}

func TopicCommands(ns string, domain string) string {
	return Topic(ns, domain, KindCommands)
}

func Topic(ns string, domain string, kind string) string {
	return fmt.Sprintf("%s.%s.%s.%s",
		ns, cqrs, domain, kind)
}
