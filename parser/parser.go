package parser

import (
	"bytes"
	"fmt"

	"roland/config"
	"roland/entity/request"
	"roland/logger"

	"github.com/hybridgroup/yzma/pkg/llama"
	"go.uber.org/zap"
)

type Parser struct {
	model llama.Model
	lctx  llama.Context
	vocab llama.Vocab

	logger *logger.Logger
	cfg *config.Config
}

func NewParser(cfg *config.Config, logger *logger.Logger) (*Parser, error) {
	logger.Info("create parser",
		zap.Any("cfg", cfg))

	logger.Info("load llama library",
		zap.String("path", cfg.LLMLibPath))

	if err := llama.Load(cfg.LLMLibPath); err != nil {
		logger.Error("failed load llama library",
			zap.String("path", cfg.LLMLibPath),
			zap.Error(err))

		return nil, fmt.Errorf("failed load llama library: %w", err)
	}

	logger.Info("init llama")

	llama.Init()

	logger.Info("load model",
		zap.String("path", cfg.ModelFile))

	model, err := llama.ModelLoadFromFile(cfg.ModelFile, llama.ModelDefaultParams())
	if err != nil {
		logger.Error("failed load model from file",
			zap.String("path", cfg.ModelFile),
			zap.Error(err))

		return nil, fmt.Errorf("failed load model from file: %w", err)
	}

	logger.Info("get context from model")

	lctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		logger.Error("failed get context from model",
			zap.Error(err))

		return nil, fmt.Errorf("failed get context from model: %w", err)
	}

	return &Parser{
		model: model,
		lctx: lctx,
		logger: logger,
		vocab: llama.ModelGetVocab(model),
		cfg: cfg,
	}, nil
}
 
func (p *Parser) Parse(phrase string) (*request.Request, error) {
	tokens := llama.Tokenize(p.vocab, phrase, true, false)

	batch := llama.BatchGetOne(tokens)
	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())

	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

	var response bytes.Buffer

	for pos := int32(0); pos < p.cfg.ResponseLength; pos += batch.NTokens {
		llama.Decode(p.lctx, batch)
		token := llama.SamplerSample(sampler, p.lctx, -1)

		if llama.VocabIsEOG(p.vocab, token) {
			fmt.Println()
			break
		}

		buf := make([]byte, 36)
		len := llama.TokenToPiece(p.vocab, token, buf, 0, true)

		response.Write(buf[:len])

		batch = llama.BatchGetOne([]llama.Token{token})
	}
} 