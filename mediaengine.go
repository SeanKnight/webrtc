// +build !js

package webrtc

import (
	"fmt"
	"strings"
	"time"

	"github.com/pion/sdp/v3"
)

const (
	mimeTypeH264 = "video/h264"
	mimeTypeOpus = "audio/opus"
	mimeTypeVP8  = "video/vp8"
	mimeTypeVP9  = "video/vp9"
	mimeTypeG722 = "audio/G722"
	mimeTypePCMU = "audio/PCMU"
	mimeTypePCMA = "audio/PCMA"
)

type mediaEngineHeaderExtension struct {
	id                               int // Negotiated ID
	uri                              string
	isAudio, isVideo                 bool
	negotiatedVideo, negotiatedAudio bool // Has this been negotiated with the remote
}

// A MediaEngine defines the codecs supported by a PeerConnection.
// MediaEngines populated using RegisterCodec (and RegisterDefaultCodecs)
// may be set up once and reused, including concurrently,
// as long as no other codecs are added subsequently.
type MediaEngine struct {
	negotiatedVideo, negotiatedAudio bool
	videoCodecs                      []RTPCodecParameters
	audioCodecs                      []RTPCodecParameters

	headerExtensions []mediaEngineHeaderExtension
}

// RegisterDefaultCodecs registers the default codecs supported by Pion WebRTC.
// RegisterDefaultCodecs is not safe for concurrent use.
func (m *MediaEngine) RegisterDefaultCodecs() error {
	// Default Pion Audio Codecs
	for _, codec := range []RTPCodecParameters{
		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeOpus, 48000, 2, "minptime=10;useinbandfec=1", []RTCPFeedback{{"transport-cc", ""}}},
			PayloadType:        111,
		},
		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeG722, 8000, 0, "", nil},
			PayloadType:        9,
		},
		{
			RTPCodecCapability: RTPCodecCapability{mimeTypePCMU, 8000, 0, "", nil},
			PayloadType:        0,
		},
		{
			RTPCodecCapability: RTPCodecCapability{mimeTypePCMA, 8000, 0, "", nil},
			PayloadType:        8,
		},
	} {
		if err := m.RegisterCodec(codec, RTPCodecTypeAudio); err != nil {
			return err
		}
	}

	// Default Pion Audio Header Extensions
	for _, extension := range []string{
		"urn:ietf:params:rtp-hdrext:sdes:mid",
		"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id",
		"urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id",
	} {
		if err := m.RegisterHeaderExtension(RTPHeaderExtensionCapability{extension}, RTPCodecTypeAudio); err != nil {
			return err
		}
	}

	videoRTCPFeedback := []RTCPFeedback{{"goog-remb", ""}, {"transport-cc", ""}, {"ccm", "fir"}, {"nack", ""}, {"nack", "pli"}}
	for _, codec := range []RTPCodecParameters{
		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeVP8, 90000, 0, "", videoRTCPFeedback},
			PayloadType:        96,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=96", nil},
			PayloadType:        97,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeVP9, 90000, 0, "profile-id=0", videoRTCPFeedback},
			PayloadType:        98,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=98", nil},
			PayloadType:        99,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeVP9, 90000, 0, "profile-id=1", videoRTCPFeedback},
			PayloadType:        100,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=100", nil},
			PayloadType:        101,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeH264, 90000, 0, "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f", videoRTCPFeedback},
			PayloadType:        102,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=102", nil},
			PayloadType:        121,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeH264, 90000, 0, "level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f", videoRTCPFeedback},
			PayloadType:        127,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=127", nil},
			PayloadType:        120,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeH264, 90000, 0, "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f", videoRTCPFeedback},
			PayloadType:        125,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=125", nil},
			PayloadType:        107,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeH264, 90000, 0, "level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42e01f", videoRTCPFeedback},
			PayloadType:        108,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=108", nil},
			PayloadType:        109,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeH264, 90000, 0, "level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f", videoRTCPFeedback},
			PayloadType:        127,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=127", nil},
			PayloadType:        120,
		},

		{
			RTPCodecCapability: RTPCodecCapability{mimeTypeH264, 90000, 0, "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=640032", videoRTCPFeedback},
			PayloadType:        123,
		},
		{
			RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=123", nil},
			PayloadType:        118,
		},

		{
			RTPCodecCapability: RTPCodecCapability{"video/ulpfec", 90000, 0, "", nil},
			PayloadType:        116,
		},
	} {
		if err := m.RegisterCodec(codec, RTPCodecTypeVideo); err != nil {
			return err
		}
	}

	// Default Pion Video Header Extensions
	for _, extension := range []string{
		"urn:ietf:params:rtp-hdrext:sdes:mid",
		"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id",
		"urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id",
	} {
		if err := m.RegisterHeaderExtension(RTPHeaderExtensionCapability{extension}, RTPCodecTypeVideo); err != nil {
			return err
		}
	}

	return nil
}

// RegisterCodec adds codec to the MediaEngine
// These are the list of codecs supported by this PeerConnection.
// RegisterCodec is not safe for concurrent use.
func (m *MediaEngine) RegisterCodec(codec RTPCodecParameters, typ RTPCodecType) error {
	codec.statsID = fmt.Sprintf("RTPCodec-%d", time.Now().UnixNano())
	switch typ {
	case RTPCodecTypeAudio:
		m.audioCodecs = append(m.audioCodecs, codec)
	case RTPCodecTypeVideo:
		m.videoCodecs = append(m.videoCodecs, codec)
	default:
		return ErrUnknownType
	}
	return nil
}

// RegisterHeaderExtension adds a header extension to the MediaEngine
// To determine the negotiated value use `GetHeaderExtensionID` after signaling is complete
func (m *MediaEngine) RegisterHeaderExtension(extension RTPHeaderExtensionCapability, typ RTPCodecType) error {
	extensionIndex := -1
	for i := range m.headerExtensions {
		if extension.URI == m.headerExtensions[i].uri {
			extensionIndex = i
		}
	}

	if extensionIndex == -1 {
		m.headerExtensions = append(m.headerExtensions, mediaEngineHeaderExtension{})
		extensionIndex = len(m.headerExtensions) - 1
	}

	if typ == RTPCodecTypeAudio {
		m.headerExtensions[extensionIndex].isAudio = true
	} else if typ == RTPCodecTypeVideo {
		m.headerExtensions[extensionIndex].isVideo = true
	}

	m.headerExtensions[extensionIndex].uri = extension.URI
	m.headerExtensions[extensionIndex].id = extensionIndex

	return nil
}

// GetHeaderExtensionID returns the negotiated ID for a header extension.
// If the Header Extension isn't enabled ok will be false
func (m *MediaEngine) GetHeaderExtensionID(extension RTPHeaderExtensionCapability) (val int, audioNegotiated, videoNegotiated bool) {
	for i := range m.headerExtensions {
		if extension.URI == m.headerExtensions[i].uri && (m.headerExtensions[i].negotiatedAudio || m.headerExtensions[i].negotiatedVideo) {
			return m.headerExtensions[i].id, m.headerExtensions[i].negotiatedAudio, m.headerExtensions[i].negotiatedVideo
		}
	}

	return
}

func (m *MediaEngine) getCodecByPayload(payloadType PayloadType) (RTPCodecParameters, error) {
	for _, codec := range m.videoCodecs {
		if codec.PayloadType == payloadType {
			return codec, nil
		}
	}
	for _, codec := range m.audioCodecs {
		if codec.PayloadType == payloadType {
			return codec, nil
		}
	}

	return RTPCodecParameters{}, ErrCodecNotFound
}

func (m *MediaEngine) getCodecsByKind(kind RTPCodecType) []RTPCodecParameters {
	switch kind {
	case RTPCodecTypeAudio:
		return m.audioCodecs
	case RTPCodecTypeVideo:
		return m.videoCodecs
	default:
		return []RTPCodecParameters{}
	}
}

func (m *MediaEngine) collectStats(collector *statsReportCollector) {
	statsLoop := func(codecs []RTPCodecParameters) {
		for _, codec := range codecs {
			collector.Collecting()
			stats := CodecStats{
				Timestamp:   statsTimestampFrom(time.Now()),
				Type:        StatsTypeCodec,
				ID:          codec.statsID,
				PayloadType: codec.PayloadType,
				MimeType:    codec.MimeType,
				ClockRate:   codec.ClockRate,
				Channels:    uint8(codec.Channels),
				SDPFmtpLine: codec.SDPFmtpLine,
			}

			collector.Collect(stats.ID, stats)
		}
	}

	statsLoop(m.videoCodecs)
	statsLoop(m.audioCodecs)
}

// Look up a codec and enable if it exists
func (m *MediaEngine) updateCodecParameters(codec RTPCodecParameters, typ RTPCodecType) error {
	codecs := m.videoCodecs
	if typ == RTPCodecTypeAudio {
		codecs = m.audioCodecs
	}

	// First attempt to match on MimeType + SDPFmtpLine
	for i := range codecs {
		if strings.EqualFold(codecs[i].RTPCodecCapability.MimeType, codec.RTPCodecCapability.MimeType) &&
			codecs[i].RTPCodecCapability.SDPFmtpLine == codec.RTPCodecCapability.SDPFmtpLine {
			codecs[i].negotiated = true
			codecs[i].PayloadType = codec.PayloadType

			return nil
		}
	}

	// Fallback to just MimeType
	for i := range codecs {
		if strings.EqualFold(codecs[i].RTPCodecCapability.MimeType, codec.RTPCodecCapability.MimeType) {
			codecs[i].negotiated = true
			codecs[i].PayloadType = codec.PayloadType
		}
	}

	return nil
}

// Look up a header extension and enable if it exists
func (m *MediaEngine) updateHeaderExtension(id int, extension string, typ RTPCodecType) error {
	for i := range m.headerExtensions {
		if m.headerExtensions[i].uri == extension {
			m.headerExtensions[i].id = id

			if m.headerExtensions[i].isAudio && typ == RTPCodecTypeAudio {
				m.headerExtensions[i].negotiatedAudio = true
			}

			if m.headerExtensions[i].isVideo && typ == RTPCodecTypeVideo {
				m.headerExtensions[i].negotiatedVideo = true
			}
		}
	}
	return nil
}

// Update the MediaEngine from a remote description
func (m *MediaEngine) updateFromRemoteDescription(desc sdp.SessionDescription) error {
	for _, media := range desc.MediaDescriptions {
		var typ RTPCodecType
		switch {
		case !m.negotiatedAudio && strings.EqualFold(media.MediaName.Media, "audio"):
			m.negotiatedAudio = true
			typ = RTPCodecTypeAudio
		case !m.negotiatedVideo && strings.EqualFold(media.MediaName.Media, "video"):
			m.negotiatedVideo = true
			typ = RTPCodecTypeVideo
		default:
			continue
		}

		codecs, err := codecsFromMediaDescription(media)
		if err != nil {
			return err
		}

		for _, codec := range codecs {
			if err = m.updateCodecParameters(codec, typ); err != nil {
				return err
			}
		}

		extensions, err := rtpExtensionsFromMediaDescription(media)
		if err != nil {
			return err
		}

		for id, extension := range extensions {
			if err = m.updateHeaderExtension(extension, id, typ); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MediaEngine) negotiatedCodecsForType(typ RTPCodecType) []RTPCodecParameters {
	if typ == RTPCodecTypeAudio {
		return m.audioCodecs
	} else if typ == RTPCodecTypeVideo {
		return m.videoCodecs
	}

	return nil
}

// func (m *MediaEngine) negotiatedHeaderExtensionsForType(RTPCodecType) []mediaEngineHeaderExtension {
// 	return m.headerExtensions
// }
