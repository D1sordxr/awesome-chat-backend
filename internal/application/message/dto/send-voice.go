package dto

type (
	SendVoiceRequest struct {
		UserID   string `json:"user_id"`
		ChatID   string `json:"chat_id"`
		Blob     string `json:"blob"`
		Duration int    `json:"duration"`
	}
	SendVoiceResponse struct {
		Waveform []byte `json:"waveform"`
		URL      string `json:"url,omitempty"`
	}
)
