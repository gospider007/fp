package fp

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gospider007/gson"
	"github.com/gospider007/gtls"
)

type Option struct {
	Addr      string
	Handler   http.Handler
	TLSConfig *tls.Config
}

func Server(handler http.Handler, options ...Option) (err error) {
	var option Option
	if len(options) > 0 {
		option = options[0]
	}
	if option.Addr == "" {
		option.Addr = ":8999"
		log.Print("Starting server on https://localhost:8999")
	}
	if option.TLSConfig == nil {
		option.TLSConfig = &tls.Config{
			GetCertificate:     gtls.GetCertificate,
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2", "http/1.1"},
		}
	}
	ln, err := NewListen(option.Addr, handler, option.TLSConfig)
	if err != nil {
		return err
	}
	defer ln.Close()
	return (&http.Server{ConnContext: ConnContext, Handler: handler}).Serve(ln)
}

func Start(option ...Option) error {
	return Server(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Connection", "close")
		rawConn := GetRawConn(r.Context())
		tlsSpec := rawConn.TLSSpec()
		h1Spec := rawConn.H1Spec()
		h2Spec := rawConn.H2Spec()
		results := map[string]any{
			"tls":          tlsSpec.Map(),
			"goSpiderSpec": rawConn.GoSpiderSpec(),
		}
		if h1Spec != nil {
			results["h1"] = h1Spec.Map()
		}
		if h2Spec != nil {
			results["h2"] = h2Spec.Map()
		}
		w.WriteHeader(200)
		con, _ := gson.Encode(results)
		w.Write(con)
	}), option...)
}
