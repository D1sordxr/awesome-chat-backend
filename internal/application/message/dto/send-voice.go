package dto

type (
	SendVoiceRequest struct {
		Message
		Blob     []byte `json:"blob"`
		Duration int    `json:"duration"`
	}
	SendVoiceResponse struct {
		Waveform []byte `json:"waveform"`
		URL      string `json:"url,omitempty"`
	}
)
