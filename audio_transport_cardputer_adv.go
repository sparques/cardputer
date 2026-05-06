//go:build esp32 && cardputer_adv

package cardputer

type advAudioTransport struct {
	cfg AudioTransportConfig
}

var sharedAudioTransport audioTransport = &advAudioTransport{}

func openAudioTransport() (audioTransport, error) {
	return sharedAudioTransport, nil
}

func configureSharedAudioTransport(cfg ES8311Config) error {
	transport, err := openAudioTransport()
	if err != nil {
		return err
	}
	return transport.Configure(audioTransportConfigFromES8311(cfg))
}

func (t *advAudioTransport) Configure(cfg AudioTransportConfig) error {
	t.cfg = cfg
	return nil
}

func (*advAudioTransport) Write([]int16) (int, error) {
	return 0, errAudioStreamUnavailable
}

func (*advAudioTransport) Read([]int16) (int, error) {
	return 0, errAudioStreamUnavailable
}
