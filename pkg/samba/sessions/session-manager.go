package sessions

import (
	"fmt"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"github.com/hirochachacha/go-smb2"
	"log/slog"
	"net"
	"sync/atomic"
)

const maxPoolSize = 20

type sessionManager struct {
	host, port     string
	user, password string
	poolSize       int32
	sessions       []*smb2.Session
	logger         *slog.Logger
}

type Session struct {
	net.Conn
	*smb2.Session
}

func NewSessionManager(logger *slog.Logger, host, port string, user, password string, poolSize int) (*sessionManager, error) {
	const op = "samba.sessions.NewSessionManager()"
	log := logger.With(slog.String("operation", op))

	if poolSize < 1 {
		poolSize = maxPoolSize
	}

	mngr := &sessionManager{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		poolSize: int32(poolSize),
		sessions: make([]*smb2.Session, poolSize),
		logger:   logger,
	}

	log.Debug("created session manager", slog.Any("pool size", poolSize))
	return mngr, nil
}

func (s *sessionManager) dial() (Session, error) {
	const op = "samba.sessions.dial()"
	logger := s.logger.With(slog.String("operation", op))

	endpoint := fmt.Sprintf("%s:%s", s.host, s.port)

	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		logger.Error("unable to connect to host via tcp", slog.String("endpoint", endpoint), sl.Err(err))
		return Session{}, errUnableConnectToHost
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{User: s.user, Password: s.password},
	}

	smbSession, err := d.Dial(conn)
	if err != nil {
		logger.Error("unable to connect to host via smb2.Dialer", slog.String("endpoint", endpoint), sl.Err(err))
		return Session{}, errUnableConnectToHost
	}

	session := Session{
		Conn:    conn,
		Session: smbSession,
	}

	return session, nil
}

func (s *sessionManager) GetSession() (Session, error) {
	const op = "samba.sessions.GetSession()"
	logger := s.logger.With(slog.String("operation", op))

	if s.poolSize <= 0 {
		return Session{}, errNoSessionAvailable
	}

	actualPoolSize := atomic.AddInt32(&s.poolSize, -1)

	logger.Debug("creating session", slog.Int("actual pool size", len(s.sessions)))

	session, err := s.dial()
	if err != nil {
		logger.Error("unable to create session", sl.Err(err))
		return Session{}, err
	}

	logger.Debug("got session", slog.Any("actual pool size", actualPoolSize))
	return session, nil
}

func (s *sessionManager) ReleaseSession(session Session) {
	const op = "samba.sessions.ReleaseSession()"
	logger := s.logger.With(slog.String("operation", op))

	logger.Debug("logging off released session")
	err := session.Logoff()
	if err != nil {
		logger.Error("unable to logoff session", sl.Err(err))
	}

	currentPoolSize := atomic.AddInt32(&s.poolSize, 1)

	logger.Debug("released session", slog.Any("actual pool size", currentPoolSize))
}
