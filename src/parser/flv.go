package parser

import (
	"github.com/go-olive/flv"
	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
)

func init() {
	SharedManager.Register(
		new(customFlv),
	)
}

type customFlv struct {
	*flv.Parser
}

func (this *customFlv) New() Parser {
	return &customFlv{
		Parser: flv.NewParser(),
	}
}

func (this *customFlv) Parse(streamURL string, out string) (err error) {
	l.Logger.WithFields(logrus.Fields{
		// "streamURL": streamURL,
		"out": out,
	}).Debug("flv working")

	return this.Parser.Parse(streamURL, out)
}
