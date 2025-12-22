// Package atem provides a Go client for communicating with Blackmagic Design ATEM switchers.
package atem

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/krombel/atem-go/models"
	"github.com/krombel/atem-go/packet"
	"github.com/krombel/atem-go/packet/cmds"
	"github.com/netlify/netlify-commons/util"
	"github.com/pkg/errors"
)

// Client represents a connection to an ATEM switcher.
type Client struct {
	host string
	port int

	conn *net.UDPConn
	log  *slog.Logger

	started   util.AtomicBool
	connected util.AtomicBool

	sessionID uint32 // atomic value (actually uint16)
	seqNum    uint32 // atomic counter (actually uint16)

	state      SwitcherState
	stateMutex sync.RWMutex

	timeRequests chan chan models.Timecode
}

// NewClient creates a new ATEM client with the default port (9910).
func NewClient(log *slog.Logger, host string) *Client {
	return NewClientWithPort(log, host, 9910)
}

// NewClientWithPort creates a new ATEM client with a custom port.
func NewClientWithPort(log *slog.Logger, host string, port int) *Client {
	return &Client{
		host:      host,
		port:      port,
		log:       log,
		started:   util.NewAtomicBool(false),
		connected: util.NewAtomicBool(false),

		timeRequests: make(chan chan models.Timecode, 100),
	}
}

// State returns the current switcher state.
func (c *Client) State() SwitcherState {
	c.stateMutex.RLock()
	state := c.state
	c.stateMutex.RUnlock()
	return state
}

// ErrAlreadyStarted is returned when attempting to start an already started client.
var ErrAlreadyStarted = errors.New("Already started")

// Start begins the connection to the ATEM switcher.
func (c *Client) Start(ctx context.Context) error {
	if c.started.Set(true) {
		return ErrAlreadyStarted
	}

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(c.host, strconv.Itoa(c.port)))
	if err != nil {
		return errors.Wrap(err, "Failed to resolve UDP address")
	}

	c.conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return errors.Wrap(err, "Failed to open UDP socket")
	}

	c.log.Debug("Created socket. Starting message processing task...")
	connected := c.startHandleMessages(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	select {
	case err, open := <-connected:
		if !open {
			return nil
		}
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// InitPayload is the initial payload sent to establish a connection with the ATEM switcher.
var InitPayload = []byte{1, 0, 0, 0, 0, 0, 0, 0}

func (c *Client) sendInit() error {
	initMSG := c.makeHeader(packet.FlagInit)
	initMSG.Length = 20 // custom well-known length
	buf := new(bytes.Buffer)
	initMSG.Serialize(buf)
	buf.Write(InitPayload)

	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) send(msg *packet.Message) error {
	buf := new(bytes.Buffer)
	msg.Serialize(buf)

	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) resetSession() {
	var sessionIDBytes [4]byte
	_, _ = rand.Read(sessionIDBytes[:])
	sessionID := binary.BigEndian.Uint32(sessionIDBytes[:])
	atomic.StoreUint32(&c.sessionID, sessionID)
	atomic.StoreUint32(&c.seqNum, 0)
	c.connected.Set(false)

	// init state
	c.state = NewSwitcherState()
}

func (c *Client) startHandleMessages(ctx context.Context) chan error {
	errCh := make(chan error)

	go func() {
		var innerCtx context.Context
		var cancel context.CancelFunc
		for ctx.Err() == nil {
			// cancel previous session
			if cancel != nil {
				cancel()
			}
			innerCtx, cancel = context.WithCancel(ctx)

			c.resetSession()
			var connected bool

			c.log.Debug("Sending init packet")
			err := c.sendInit()
			if err != nil {
				errCh <- errors.Wrap(err, "Failed to send init message")
				break
			}

			// ack packets in background
			ackQueue := make(chan uint16, 20) // seq nums to ack
			go c.ackWorker(innerCtx, ackQueue)

			// read loop
			var rerr error
			for {
				raw := make([]byte, 1500) // safe size for expected MTU
				// todo: read deadline
				n, err := c.conn.Read(raw)
				if err != nil {
					rerr = err
					break
				}
				c.log.Debug("Got packet", "len", n)

				buf := bytes.NewBuffer(raw[:n])
				msg, err := packet.Deserialize(c.log, buf)
				if err != nil {
					rerr = err
					break
				}
				log := c.log.With("seq", msg.SeqNum)
				flags := msg.Flags.Debug()
				attrs := make([]interface{}, 0, len(flags)*2)
				for k, v := range flags {
					attrs = append(attrs, k, v)
				}
				log.Debug("Read packet", attrs...)

				oldSession := atomic.SwapUint32(&c.sessionID, uint32(msg.SessionID))
				if oldSession != uint32(msg.SessionID) {
					log.Info("Using new session id from switcher", "session", msg.SessionID)
				}

				if msg.Flags.Has(packet.FlagRetrans) {
					// todo: handle double sends
					log.Warn("Retransmission detected")
				}

				if msg.Flags.Has(packet.FlagInit) {
					// always ack init packets
					msg.Flags |= packet.FlagNeedACK
				}

				if msg.Flags.Has(packet.FlagNeedACK) {
					log.Debug("queueing packet for ack")
					ackQueue <- msg.SeqNum
				}

				gotInit, err := c.processCommands(log, msg)
				if err != nil {
					log.Warn("Error while processing commands", "error", err)
				}

				if gotInit && !connected {
					connected = true
					c.connected.Set(true)
					errCh <- nil // avoid risk of double-closing
				}
			}
			if rerr != nil {
				c.log.Warn("Failed to read", "error", rerr)
				if !connected {
					errCh <- errors.Wrap(rerr, "Failed to read message")
				}
				break
			}
		}
		if cancel != nil {
			cancel()
		}
		_ = c.conn.Close()
		c.started.Set(false)
	}()

	return errCh
}

func (c *Client) processCommands(log *slog.Logger, msg packet.Message) (init bool, err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	for _, cmd := range msg.Commands {
		log.Debug("Starting to process command", "slug", cmd.Slug())

		switch t := cmd.(type) {
		case *cmds.UnknownCommand:
			log.Debug("Got unknown command", "slug", t.Slug())
		case *cmds.IncmCmd:
			init = true
		case *cmds.VerCmd:
			c.state.Version.Major = int(t.Major)
			c.state.Version.Minor = int(t.Minor)
			log.Debug("Got version", "major", t.Major, "minor", t.Minor)
		case *cmds.PinCmd:
			c.state.Config.ProductName = t.Value()
		case *cmds.WarnCmd:
			c.state.Warning = t.Value()
		case *cmds.TopCmd:
			c.state.Config.Topology = models.Topology(*t)
		case *cmds.MecCmd:
			c.state.Config.MixEffect[int(t.ME)] = int(t.Keyers)
		case *cmds.MplCmd:
			c.state.Config.MediaPlayer = models.MediaPlayerConfig(*t)
		case *cmds.MvcCmd:
			c.state.Config.MultiViews = int(*t)
		case *cmds.SscCmd:
			c.state.Config.SuperSources = int(*t)
		case *cmds.TlcCmd:
			c.state.Config.TallyChannels = int(*t)
		case *cmds.MacCmd:
			c.state.Config.MacroBanks = int(*t)
		case *cmds.PowrCmd:
			c.state.Power = models.PowerStatus(*t)
		case *cmds.VidmCmd:
			c.state.VideoMode = models.VideoMode(*t)
		case *cmds.InprCmd:
			c.state.Inputs[t.SourceIndex] = models.InputProperties(*t)
		case *cmds.PrgiCmd:
			c.state.Program[int(t.Bus)] = t.Source
		case *cmds.PrviCmd:
			c.state.Preview[int(t.Bus)] = t.Source
		case *cmds.AuxsCmd:
			c.state.Aux[int(t.Bus)] = t.Source
		case *cmds.MpceCmd:
			mp, ok := c.state.MediaPlayer[int(t.PlayerIndex)]
			if !ok {
				mp = &models.MediaPlayer{}
				c.state.MediaPlayer[int(t.PlayerIndex)] = mp
			}
			mp.Type = t.Type
			mp.StillIndex = t.StillIndex
			mp.ClipIndex = t.ClipIndex
		case *cmds.MpfeCmd:
			c.state.MediaFiles[int(t.Index)] = models.MediaStillFrame{
				Used:     t.Used,
				Hash:     t.Hash,
				Filename: t.Filename,
			}
		case cmds.TlinCmd:
			c.state.TallyByIndex = t
		case cmds.TlsrCmd:
			c.state.TallyBySource = t
		case *cmds.TimeCmd:
			c.handleTimeCmd(t)
		case *cmds.KeOnCmd:
			c.handleKeOnCmd(t)
		case *cmds.KeBPCmd:
			c.handleKeBPCmd(t)
		case *cmds.KeDVCmd:
			c.handleKeDVCmd(t)
		default:
			log.Debug("Unhandled packet type")
		}
	}
	return
}

func (c *Client) handleTimeCmd(t *cmds.TimeCmd) {
	tc := models.Timecode(*t)
	done := false
	for !done {
		select {
		case retCh, open := <-c.timeRequests:
			if !open {
				done = true
			}
			retCh <- tc
			close(retCh)
		default:
			// no more in channel
			done = true
		}
	}
	c.state.TimeCodeLastChange = tc
}

func (c *Client) handleKeOnCmd(t *cmds.KeOnCmd) {
	me := int(t.ME)
	keyer := int(t.Keyer)
	if c.state.Keyers[me] == nil {
		c.state.Keyers[me] = make(map[int]models.KeyerState)
	}
	state := c.state.Keyers[me][keyer]
	// OnAir might be a bitmask: bit 0 = program, bit 1 = preview
	state.OnAir = (t.OnAir & 0x01) != 0
	state.PreviewOnAir = (t.OnAir & 0x02) != 0
	c.state.Keyers[me][keyer] = state
}

func (c *Client) handleKeBPCmd(t *cmds.KeBPCmd) {
	me := int(t.ME)
	keyer := int(t.Keyer)
	if c.state.Keyers[me] == nil {
		c.state.Keyers[me] = make(map[int]models.KeyerState)
	}
	state := c.state.Keyers[me][keyer]
	state.FillSource = t.FillSource
	state.KeySource = t.KeySource
	c.state.Keyers[me][keyer] = state
}

func (c *Client) handleKeDVCmd(t *cmds.KeDVCmd) {
	me := int(t.ME)
	keyer := int(t.Keyer)
	if c.state.Keyers[me] == nil {
		c.state.Keyers[me] = make(map[int]models.KeyerState)
	}
	state := c.state.Keyers[me][keyer]
	if state.DVE == nil {
		state.DVE = &models.DVEProperties{}
	}
	state.DVE.Rate = t.Rate
	state.DVE.SizeX = t.SizeX
	state.DVE.SizeY = t.SizeY
	state.DVE.PositionX = t.PositionX
	state.DVE.PositionY = t.PositionY
	state.DVE.Rotation = t.Rotation
	state.DVE.BorderEnabled = t.BorderEnabled
	state.DVE.BorderStyle = t.BorderStyle
	state.DVE.BorderWidth = t.BorderWidth
	state.DVE.BorderOpacity = t.BorderOpacity
	state.DVE.ShadowEnabled = t.ShadowEnabled
	state.DVE.LightSource = t.LightSource
	state.DVE.MaskEnabled = t.MaskEnabled
	state.DVE.MaskTop = t.MaskTop
	state.DVE.MaskBottom = t.MaskBottom
	state.DVE.MaskLeft = t.MaskLeft
	state.DVE.MaskRight = t.MaskRight
	c.state.Keyers[me][keyer] = state
}

const ackDebounceInterval = time.Millisecond * 20
const ackMaxDebounce = time.Millisecond * 100

// since the protocol seems to allow just acking the seq num
// of the packet that was received last, the actual sending
// of the ack is debounced to improve performance
func (c *Client) ackWorker(ctx context.Context, queue <-chan uint16) {
	timer := time.NewTimer(ackDebounceInterval)
	timer.Stop() // start with a stopped timer

	lastSend := time.Now()
	var latestNum uint16
	for {
		log := c.log.With("num", latestNum)
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Debug("ACK worker stopped", "error", ctx.Err())
			return
		case seqNum, open := <-queue:
			if !open {
				log.Debug("ACK work channel closed")
				return
			}
			if seqNum > latestNum {
				latestNum = seqNum
			}
			timer.Stop()
			if time.Since(lastSend) > ackMaxDebounce {
				timer.Reset(0)
			} else {
				timer.Reset(ackDebounceInterval)
			}
		case <-timer.C:
			log.Debug("Ack timer fired. Sending ack...")
			// defer timer fired, actually send ack
			msg := c.makeHeader(packet.FlagACK)
			msg.AckID = latestNum
			err := c.send(msg)
			if err != nil {
				log.Warn("Failed sending ack", "error", err)
			}
			lastSend = time.Now()
		}
	}
}

func (c *Client) makeHeader(flags ...packet.Flags) *packet.Message {
	seqNum := atomic.AddUint32(&c.seqNum, 1) - 1 // allow seq num to start at 0
	return &packet.Message{
		Flags:     packet.FlagsFrom(flags...),
		SessionID: uint16(atomic.LoadUint32(&c.sessionID)), //nolint:gosec // sessionID is always within uint16 range
		SeqNum:    uint16(seqNum),                          //nolint:gosec // seqNum is always within uint16 range
	}
}
