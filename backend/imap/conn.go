package imap

import (
	"sync"
	"errors"

	"github.com/mxk/go-imap/imap"
)

type Config struct {
	Host string
	Suffix string
}

type connBackend struct {
	config *Config
	conns map[string]*imap.Client
	locks map[string]sync.Locker
}

func (b *connBackend) insertConn(user string, conn *imap.Client) {
	b.conns[user] = conn
	b.locks[user] = &sync.Mutex{}
}

func (b *connBackend) getConn(user string) (*imap.Client, func(), error) {
	lock, ok := b.locks[user]
	if !ok {
		return nil, nil, errors.New("No such user")
	}

	lock.Lock()

	conn, ok := b.conns[user]
	if !ok {
		return nil, nil, errors.New("No such user")
	}

	return conn, lock.Unlock, nil
}

func newConnBackend() *connBackend {
	return &connBackend{
		// TODO: make this configurable
		config: &Config{
			Host: "mail.gandi.net",
			Suffix: "@emersion.fr",
		},

		conns: map[string]*imap.Client{},
		locks: map[string]sync.Locker{},
	}
}