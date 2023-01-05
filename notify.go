package notify

import "github.com/pkg/errors"

var _ Notifier = (*Notify)(nil)

var ErrSendNotification = errors.New("send notification")

type Notify struct {
	Disabled  bool
	notifiers []Notifier
}

type Option func(*Notify)

func Enable(n *Notify) {
	if n != nil {
		n.Disabled = false
	}
}

func Disable(n *Notify) {
	if n != nil {
		n.Disabled = true
	}
}

func (n *Notify) WithOptions(options ...Option) *Notify {
	if options == nil {
		return n
	}

	for _, option := range options {
		if option != nil {
			option(n)
		}
	}

	return n
}

func NewWithOptions(options ...Option) *Notify {
	n := &Notify{
		Disabled:  false,               // Enabled by default.
		notifiers: make([]Notifier, 0), // Avoid nil list.
	}

	return n.WithOptions(options...)
}

func New() *Notify {
	return NewWithOptions()
}

func NewWithServices(services ...Notifier) *Notify {
	n := New()
	n.UseServices(services...)

	return n
}

var std = New()

func Default() *Notify {
	return std
}
