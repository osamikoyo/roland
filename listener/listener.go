package listener

import (
	"encoding/json"
	"fmt"
	"roland/config"
	"roland/logger"

	"github.com/alphacep/vosk-api/go"
	"github.com/gordonklaus/portaudio"
	"go.uber.org/zap"
)

type Listener struct {
	rec      *vosk.VoskRecognizer
	logger   *logger.Logger
	model    *vosk.VoskModel
	outBound chan string
}

func NewListener(cfg *config.Config, logger *logger.Logger, outbound chan string) (*Listener, error) {
	logger.Info("setup model")

	model, err := vosk.NewModel(cfg.LLMs.Listener)
	if err != nil {
		logger.Error("failed create model",
			zap.String("path", cfg.LLMs.Listener),
			zap.Error(err))

		return nil, fmt.Errorf("failed create model: %w", err)
	}

	logger.Info("setup recognizer")

	rec, err := vosk.NewRecognizer(model, 1600)
	if err != nil {
		logger.Error("failed create recognizer",
			zap.Error(err))

		return nil, fmt.Errorf("failed create recognizer: %w", err)
	}

	return &Listener{
		rec:    rec,
		logger: logger,
	}, nil
}

func (lis *Listener) Start() error {
	lis.logger.Info("init portaudio")

	if err := portaudio.Initialize(); err != nil {
		lis.logger.Error("failed init portaudio",
			zap.Error(err))

		return fmt.Errorf("failed init portaudio: %w", err)
	}

	lis.logger.Info("create stream")

	stream, err := portaudio.OpenDefaultStream(1, 0, 1600, 512, func(in []byte) {
		lis.rec.AcceptWaveform(in)

		result := lis.rec.Result()
		var res struct{ Text string }

		if err := json.Unmarshal([]byte(result), &res); err != nil {
			lis.logger.Error("failed unmarshal result",
				zap.Error(err))
		}

		if res.Text != "" {
			lis.outBound <- res.Text
		}
	})
	if err != nil {
		lis.logger.Error("failed open stream",
			zap.Error(err))

		return fmt.Errorf("failed open stream: %w", err)
	}

	lis.logger.Info("start stream")

	if err = stream.Start(); err != nil {
		lis.logger.Error("failed start stream",
			zap.Error(err))

		return fmt.Errorf("failed start stream: %w", err)
	}

	return nil
}

func (lis *Listener) Close() {
	lis.rec.Free()
	lis.model.Free()
	portaudio.Terminate()
}
