package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

// PrintDump prints dump of request, optionally writing it in the response
func (c *Context) PrintDump(w http.ResponseWriter, r *http.Request, write bool) {
	dump, _ := httputil.DumpRequest(r, true)
	c.Logger.Debugf("%v", string(dump))
	if write == true {
		w.Write(dump)
	}
}

func (c *Context) ConfugureLogger() {
	level, err := logrus.ParseLevel(c.Config.LogLevel)
	if err == nil {
		c.Logger.SetLevel(level)
	} else {
		c.Logger.SetLevel(logrus.DebugLevel)
		c.Logger.Errorln("Can't configure logger. Using default log level: debug.")
	}
}
