package notify

func (n *Notify) useService(service Notifier) {
	if service != nil {
		n.notifiers = append(n.notifiers, service)
	}
}

func (n *Notify) useServices(services ...Notifier) {
	for _, s := range services {
		n.useService(s)
	}
}

func (n *Notify) UseServices(services ...Notifier) {
	n.useServices(services...)
}

func UseServices(services ...Notifier) {
	std.UseServices(services...)
}