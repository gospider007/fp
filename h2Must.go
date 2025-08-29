package fp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"golang.org/x/net/http2/hpack"
)

type errHttp2Code uint32

const (
	errHttp2CodeProtocol    errHttp2Code = 0x1
	errHttp2CodeFlowControl errHttp2Code = 0x3
	errHttp2CodeFrameSize   errHttp2Code = 0x6
	errHttp2CodeCompression errHttp2Code = 0x9

	// errHttp2CodeNo                 errHttp2Code = 0x0
	// errHttp2CodeInternal           errHttp2Code = 0x2
	// errHttp2CodeSettingsTimeout    errHttp2Code = 0x4
	// errHttp2CodeStreamClosed       errHttp2Code = 0x5
	// errHttp2CodeRefusedStream      errHttp2Code = 0x7
	// errHttp2CodeCancel             errHttp2Code = 0x8
	// errHttp2CodeConnect            errHttp2Code = 0xa
	// errHttp2CodeEnhanceYourCalm    errHttp2Code = 0xb
	// errHttp2CodeInadequateSecurity errHttp2Code = 0xc
	// errHttp2CodeHTTP11Required     errHttp2Code = 0xd
)

type Http2ReaderFramer struct {
	r               io.Reader
	getReadBuf      func(size uint32) []byte
	ReadMetaHeaders *hpack.Decoder
	frameCache      *http2frameCache
	readBuf         []byte
	headerBuf       [http2frameHeaderLen]byte
}

func Http2NewReaderFramer(r io.Reader) *Http2ReaderFramer {
	fr := &Http2ReaderFramer{
		r:               r,
		ReadMetaHeaders: hpack.NewDecoder(65536, nil),
	}
	fr.getReadBuf = func(size uint32) []byte {
		if cap(fr.readBuf) >= int(size) {
			return fr.readBuf[:size]
		}
		fr.readBuf = make([]byte, size)
		return fr.readBuf
	}
	return fr
}

type http2ConnectionError errHttp2Code

func (e http2ConnectionError) Error() string {
	return fmt.Sprintf("connection error: %d", errHttp2Code(e))
}

const http2frameHeaderLen = 9

type http2FrameType uint8

const (
	http2FrameData         http2FrameType = 0x0
	http2FrameHeaders      http2FrameType = 0x1
	http2FramePriority     http2FrameType = 0x2
	http2FrameRSTStream    http2FrameType = 0x3
	http2FrameSettings     http2FrameType = 0x4
	http2FramePushPromise  http2FrameType = 0x5
	http2FramePing         http2FrameType = 0x6
	http2FrameGoAway       http2FrameType = 0x7
	http2FrameWindowUpdate http2FrameType = 0x8
	http2FrameContinuation http2FrameType = 0x9
)

var http2frameName = map[http2FrameType]string{
	http2FrameData:         "DATA",
	http2FrameHeaders:      "HEADERS",
	http2FramePriority:     "PRIORITY",
	http2FrameRSTStream:    "RST_STREAM",
	http2FrameSettings:     "SETTINGS",
	http2FramePushPromise:  "PUSH_PROMISE",
	http2FramePing:         "PING",
	http2FrameGoAway:       "GOAWAY",
	http2FrameWindowUpdate: "WINDOW_UPDATE",
	http2FrameContinuation: "CONTINUATION",
}

func (t http2FrameType) String() string {
	if s, ok := http2frameName[t]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN_FRAME_TYPE_%d", uint8(t))
}

type http2Flags uint8

func (f http2Flags) Has(v http2Flags) bool {
	return (f & v) == v
}

const (
	http2FlagDataEndStream http2Flags = 0x1
	http2FlagDataPadded    http2Flags = 0x8

	http2FlagHeadersEndStream  http2Flags = 0x1
	http2FlagHeadersEndHeaders http2Flags = 0x4
	http2FlagHeadersPadded     http2Flags = 0x8
	http2FlagHeadersPriority   http2Flags = 0x20

	http2FlagSettingsAck http2Flags = 0x1

	http2FlagPingAck http2Flags = 0x1

	http2FlagContinuationEndHeaders http2Flags = 0x4

	http2FlagPushPromiseEndHeaders http2Flags = 0x4
	http2FlagPushPromisePadded     http2Flags = 0x8
)

type http2frameParser func(fc *http2frameCache, fh http2FrameHeader, payload []byte) (any, error)

func http2parseHeadersFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (_ any, err error) {
	hf := &Http2HeadersFrame{
		http2FrameHeader: fh,
	}
	var padLength uint8
	if fh.Flags.Has(http2FlagHeadersPadded) {
		if p, padLength, err = http2readByte(p); err != nil {
			return
		}
	}
	if fh.Flags.Has(http2FlagHeadersPriority) {
		var v uint32
		p, v, err = http2readUint32(p)
		if err != nil {
			return nil, err
		}
		p, hf.Priority.Weight, err = http2readByte(p)
		if err != nil {
			return nil, err
		}
		hf.Priority.StreamDep = v & 0x7fffffff
		hf.Priority.Exclusive = (v != hf.Priority.StreamDep)
	}
	if len(p)-int(padLength) < 0 {
		return nil, errors.New("frame_headers_pad_too_big")
	}
	hf.headerFragBuf = p[:len(p)-int(padLength)]
	return hf, nil
}

var http2frameParsers = map[http2FrameType]http2frameParser{
	http2FrameData:         http2parseDataFrame,
	http2FrameHeaders:      http2parseHeadersFrame,
	http2FramePriority:     http2parsePriorityFrame,
	http2FrameRSTStream:    http2parseRSTStreamFrame,
	http2FrameSettings:     http2parseSettingsFrame,
	http2FramePushPromise:  http2parsePushPromise,
	http2FramePing:         http2parsePingFrame,
	http2FrameGoAway:       http2parseGoAwayFrame,
	http2FrameWindowUpdate: http2parseWindowUpdateFrame,
	http2FrameContinuation: http2parseContinuationFrame,
}

func http2typeFrameParser(t http2FrameType) http2frameParser {
	if f := http2frameParsers[t]; f != nil {
		return f
	}
	return http2parseUnknownFrame
}

type http2FrameHeader struct {
	Type     http2FrameType
	Flags    http2Flags
	Length   uint32
	StreamID uint32
}

func http2readFrameHeader(buf []byte, r io.Reader) (http2FrameHeader, error) {
	_, err := io.ReadFull(r, buf[:http2frameHeaderLen])
	if err != nil {
		return http2FrameHeader{}, err
	}
	return http2FrameHeader{
		Length:   (uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])),
		Type:     http2FrameType(buf[3]),
		Flags:    http2Flags(buf[4]),
		StreamID: binary.BigEndian.Uint32(buf[5:]) & (1<<31 - 1),
	}, nil
}

type http2frameCache struct {
	dataFrame Http2DataFrame
}

func (fc *http2frameCache) getDataFrame() *Http2DataFrame {
	if fc == nil {
		return &Http2DataFrame{}
	}
	return &fc.dataFrame
}

func (fr *Http2ReaderFramer) ReadFrame() (any, []byte, error) {
	fh, err := http2readFrameHeader(fr.headerBuf[:], fr.r)
	if err != nil {
		return nil, nil, err
	}
	payload := fr.getReadBuf(fh.Length)
	if _, err := io.ReadFull(fr.r, payload); err != nil {
		return nil, nil, err
	}
	data := fr.headerBuf[:http2frameHeaderLen]
	data = append(data, payload...)
	f, err := http2typeFrameParser(fh.Type)(fr.frameCache, fh, payload)
	if err != nil {
		return nil, nil, err
	}
	if setting, ok := f.(*Http2SettingsFrame); ok {
		setting.ForeachSetting(
			func(hs Http2Setting) error {
				if hs.ID == Http2SettingHeaderTableSize {
					fr.ReadMetaHeaders.SetMaxDynamicTableSize(hs.Val)
				}
				return nil
			},
		)
	}
	if fh.Type == http2FrameHeaders {
		f2, data2, err2 := fr.readMetaFrame(f.(*Http2HeadersFrame))
		if err2 != nil {
			return nil, nil, err2
		}
		data = append(data, data2...)
		return f2, data, nil
	}
	return f, data, nil
}
func (fr *Http2ReaderFramer) readMetaFrame(hf *Http2HeadersFrame) (any, []byte, error) {
	mh := &Http2MetaHeadersFrame{
		Http2HeadersFrame: hf,
	}
	fr.ReadMetaHeaders.SetEmitEnabled(true)
	fr.ReadMetaHeaders.SetEmitFunc(func(hf hpack.HeaderField) {
		mh.Fields = append(mh.Fields, hf)
	})
	defer func() {
		fr.ReadMetaHeaders.SetEmitEnabled(false)
		fr.ReadMetaHeaders.SetEmitFunc(nil)
	}()
	var hc http2headersOrContinuation = hf
	allData := []byte{}
	for {
		frag := hc.HeaderBlockFragment()
		if _, err := fr.ReadMetaHeaders.Write(frag); err != nil {
			return mh, nil, http2ConnectionError(errHttp2CodeCompression)
		}
		if hc.HeadersEnded() {
			break
		}
		if f, data, err := fr.ReadFrame(); err != nil {
			return nil, nil, err
		} else {
			hc = f.(*http2ContinuationFrame)
			allData = append(allData, data...)
		}
	}
	mh.Http2HeadersFrame.headerFragBuf = nil
	if err := fr.ReadMetaHeaders.Close(); err != nil {
		return mh, nil, http2ConnectionError(errHttp2CodeCompression)
	}
	return mh, allData, nil
}

type Http2DataFrame struct {
	data []byte
	http2FrameHeader
}

func (f *Http2DataFrame) StreamEnded() bool {
	return f.http2FrameHeader.Flags.Has(http2FlagDataEndStream)
}

func (f *Http2DataFrame) Data() []byte {
	return f.data
}

func http2parseDataFrame(fc *http2frameCache, fh http2FrameHeader, payload []byte) (any, error) {
	if fh.StreamID == 0 {
		return nil, errors.New("DATA frame with stream ID 0")
	}
	f := fc.getDataFrame()
	f.http2FrameHeader = fh

	var padSize byte
	if fh.Flags.Has(http2FlagDataPadded) {
		var err error
		payload, padSize, err = http2readByte(payload)
		if err != nil {
			return nil, err
		}
	}
	if int(padSize) > len(payload) {
		return nil, errors.New("pad size larger than data payload")
	}
	f.data = payload[:len(payload)-int(padSize)]
	return f, nil
}

type Http2SettingsFrame struct {
	p []byte
	http2FrameHeader
}

func http2parseSettingsFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (any, error) {
	if fh.Flags.Has(http2FlagSettingsAck) && fh.Length > 0 {
		return nil, http2ConnectionError(errHttp2CodeFrameSize)
	}
	if fh.StreamID != 0 {
		return nil, http2ConnectionError(errHttp2CodeProtocol)
	}
	if len(p)%6 != 0 {
		return nil, http2ConnectionError(errHttp2CodeFrameSize)
	}
	f := &Http2SettingsFrame{http2FrameHeader: fh, p: p}
	if v, ok := f.Value(Http2SettingInitialWindowSize); ok && v > (1<<31)-1 {
		return nil, http2ConnectionError(errHttp2CodeFlowControl)
	}
	return f, nil
}

func (f *Http2SettingsFrame) IsAck() bool {
	return f.http2FrameHeader.Flags.Has(http2FlagSettingsAck)
}

func (f *Http2SettingsFrame) Value(id Http2SettingID) (v uint32, ok bool) {
	for i := 0; i < f.NumSettings(); i++ {
		if s := f.Setting(i); s.ID == id {
			return s.Val, true
		}
	}
	return 0, false
}

func (f *Http2SettingsFrame) Setting(i int) Http2Setting {
	buf := f.p
	return Http2Setting{
		ID:  Http2SettingID(binary.BigEndian.Uint16(buf[i*6 : i*6+2])),
		Val: binary.BigEndian.Uint32(buf[i*6+2 : i*6+6]),
	}
}

func (f *Http2SettingsFrame) NumSettings() int { return len(f.p) / 6 }

func (f *Http2SettingsFrame) HasDuplicates() bool {
	num := f.NumSettings()
	if num == 0 {
		return false
	}

	if num < 10 {
		for i := 0; i < num; i++ {
			idi := f.Setting(i).ID
			for j := i + 1; j < num; j++ {
				idj := f.Setting(j).ID
				if idi == idj {
					return true
				}
			}
		}
		return false
	}
	seen := map[Http2SettingID]bool{}
	for i := 0; i < num; i++ {
		id := f.Setting(i).ID
		if seen[id] {
			return true
		}
		seen[id] = true
	}
	return false
}

func (f *Http2SettingsFrame) ForeachSetting(fn func(Http2Setting) error) error {
	for i := 0; i < f.NumSettings(); i++ {
		if err := fn(f.Setting(i)); err != nil {
			return err
		}
	}
	return nil
}

type Http2PingFrame struct {
	http2FrameHeader
	Data [8]byte
}

func (f *Http2PingFrame) IsAck() bool { return f.Flags.Has(http2FlagPingAck) }

func http2parsePingFrame(_ *http2frameCache, fh http2FrameHeader, payload []byte) (any, error) {
	if len(payload) != 8 {
		return nil, http2ConnectionError(errHttp2CodeFrameSize)
	}
	if fh.StreamID != 0 {
		return nil, http2ConnectionError(errHttp2CodeProtocol)
	}
	f := &Http2PingFrame{http2FrameHeader: fh}
	copy(f.Data[:], payload)
	return f, nil
}

type Http2GoAwayFrame struct {
	debugData []byte
	http2FrameHeader
	LastStreamID uint32
	ErrCode      errHttp2Code
}

func http2parseGoAwayFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (any, error) {
	if fh.StreamID != 0 {
		return nil, http2ConnectionError(errHttp2CodeProtocol)
	}
	if len(p) < 8 {
		return nil, http2ConnectionError(errHttp2CodeFrameSize)
	}
	return &Http2GoAwayFrame{
		http2FrameHeader: fh,
		LastStreamID:     binary.BigEndian.Uint32(p[:4]) & (1<<31 - 1),
		ErrCode:          errHttp2Code(binary.BigEndian.Uint32(p[4:8])),
		debugData:        p[8:],
	}, nil
}

type http2UnknownFrame struct {
	p []byte
	http2FrameHeader
}

func (f *http2UnknownFrame) Payload() []byte {
	return f.p
}

func http2parseUnknownFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (any, error) {
	return &http2UnknownFrame{p, fh}, nil
}

type Http2WindowUpdateFrame struct {
	http2FrameHeader
	Increment uint32
}

func http2parseWindowUpdateFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (any, error) {
	if len(p) != 4 {
		return nil, http2ConnectionError(errHttp2CodeFrameSize)
	}
	return &Http2WindowUpdateFrame{
		http2FrameHeader: fh,
		Increment:        binary.BigEndian.Uint32(p[:4]) & 0x7fffffff,
	}, nil
}

type Http2HeadersFrame struct {
	headerFragBuf []byte
	http2FrameHeader
	Priority Http2PriorityParam
}

func (f *Http2HeadersFrame) HeaderBlockFragment() []byte {
	return f.headerFragBuf
}

func (f *Http2HeadersFrame) HeadersEnded() bool {
	return f.http2FrameHeader.Flags.Has(http2FlagHeadersEndHeaders)
}

func (f *Http2HeadersFrame) StreamEnded() bool {
	return f.http2FrameHeader.Flags.Has(http2FlagHeadersEndStream)
}

type Http2HeadersFrameParam struct {
	BlockFragment []byte
	Priority      Http2PriorityParam
	StreamID      uint32
	EndStream     bool
	EndHeaders    bool
	PadLength     uint8
}

type http2PriorityFrame struct {
	http2FrameHeader
	Http2PriorityParam
}

type Http2PriorityParam struct {
	StreamDep uint32

	Exclusive bool

	Weight uint8
}

func (p Http2PriorityParam) IsZero() bool {
	return p == Http2PriorityParam{}
}

func http2parsePriorityFrame(_ *http2frameCache, fh http2FrameHeader, payload []byte) (any, error) {
	if fh.StreamID == 0 {
		return nil, errors.New("PRIORITY frame with stream ID 0")
	}
	if len(payload) != 5 {
		return nil, fmt.Errorf("PRIORITY frame payload size was %d; want 5", len(payload))
	}
	v := binary.BigEndian.Uint32(payload[:4])
	streamID := v & 0x7fffffff
	return &http2PriorityFrame{
		http2FrameHeader: fh,
		Http2PriorityParam: Http2PriorityParam{
			Weight:    payload[4],
			StreamDep: streamID,
			Exclusive: streamID != v,
		},
	}, nil
}

type Http2RSTStreamFrame struct {
	http2FrameHeader
	ErrCode errHttp2Code
}

func http2parseRSTStreamFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (any, error) {
	if len(p) != 4 {
		return nil, http2ConnectionError(errHttp2CodeFrameSize)
	}
	if fh.StreamID == 0 {
		return nil, http2ConnectionError(errHttp2CodeProtocol)
	}
	return &Http2RSTStreamFrame{fh, errHttp2Code(binary.BigEndian.Uint32(p[:4]))}, nil
}

type http2ContinuationFrame struct {
	headerFragBuf []byte
	http2FrameHeader
}

func http2parseContinuationFrame(_ *http2frameCache, fh http2FrameHeader, p []byte) (any, error) {
	if fh.StreamID == 0 {
		return nil, errors.New("CONTINUATION frame with stream ID 0")
	}
	return &http2ContinuationFrame{p, fh}, nil
}

func (f *http2ContinuationFrame) HeaderBlockFragment() []byte {
	return f.headerFragBuf
}

func (f *http2ContinuationFrame) HeadersEnded() bool {
	return f.http2FrameHeader.Flags.Has(http2FlagContinuationEndHeaders)
}

type Http2PushPromiseFrame struct {
	headerFragBuf []byte
	http2FrameHeader
	PromiseID uint32
}

func (f *Http2PushPromiseFrame) HeaderBlockFragment() []byte {
	return f.headerFragBuf
}

func (f *Http2PushPromiseFrame) HeadersEnded() bool {
	return f.http2FrameHeader.Flags.Has(http2FlagPushPromiseEndHeaders)
}

func http2parsePushPromise(_ *http2frameCache, fh http2FrameHeader, p []byte) (_ any, err error) {
	pp := &Http2PushPromiseFrame{
		http2FrameHeader: fh,
	}
	if pp.StreamID == 0 {
		return nil, http2ConnectionError(errHttp2CodeProtocol)
	}
	var padLength uint8
	if fh.Flags.Has(http2FlagPushPromisePadded) {
		if p, padLength, err = http2readByte(p); err != nil {
			return
		}
	}
	p, pp.PromiseID, err = http2readUint32(p)
	if err != nil {
		return
	}
	pp.PromiseID = pp.PromiseID & (1<<31 - 1)

	if int(padLength) > len(p) {

		return nil, http2ConnectionError(errHttp2CodeProtocol)
	}
	pp.headerFragBuf = p[:len(p)-int(padLength)]
	return pp, nil
}

func http2readByte(p []byte) (remain []byte, b byte, err error) {
	if len(p) == 0 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	return p[1:], p[0], nil
}

func http2readUint32(p []byte) (remain []byte, v uint32, err error) {
	if len(p) < 4 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	return p[4:], binary.BigEndian.Uint32(p[:4]), nil
}

type http2headersEnder interface {
	HeadersEnded() bool
}

type http2headersOrContinuation interface {
	http2headersEnder
	HeaderBlockFragment() []byte
}

type Http2MetaHeadersFrame struct {
	*Http2HeadersFrame
	Fields []hpack.HeaderField
}

func (mh *Http2MetaHeadersFrame) PseudoValue(pseudo string) string {
	for _, hf := range mh.Fields {
		if !hf.IsPseudo() {
			return ""
		}
		if hf.Name[1:] == pseudo {
			return hf.Value
		}
	}
	return ""
}

func (mh *Http2MetaHeadersFrame) RegularFields() []hpack.HeaderField {
	for i, hf := range mh.Fields {
		if !hf.IsPseudo() {
			return mh.Fields[i:]
		}
	}
	return nil
}

func (mh *Http2MetaHeadersFrame) PseudoFields() []hpack.HeaderField {
	for i, hf := range mh.Fields {
		if !hf.IsPseudo() {
			return mh.Fields[:i]
		}
	}
	return mh.Fields
}

type Http2Setting struct {
	ID Http2SettingID

	Val uint32
}

type Http2SettingID uint16

const (
	Http2SettingHeaderTableSize   Http2SettingID = 0x1
	Http2SettingInitialWindowSize Http2SettingID = 0x4
	Http2SettingMaxFrameSize      Http2SettingID = 0x5
	Http2SettingMaxHeaderListSize Http2SettingID = 0x6
)
