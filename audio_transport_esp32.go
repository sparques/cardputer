//go:build esp32 && !cardputer_adv

package cardputer

type esp32AudioTransport struct {
	cfg AudioTransportConfig
}

var sharedAudioTransport audioTransport = &esp32AudioTransport{}

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

func (t *esp32AudioTransport) Configure(cfg AudioTransportConfig) error {
	t.cfg = cfg
	return nil
}

func (*esp32AudioTransport) Write([]int16) (int, error) {
	return 0, errAudioStreamUnavailable
}

func (*esp32AudioTransport) Read([]int16) (int, error) {
	return 0, errAudioStreamUnavailable
}
