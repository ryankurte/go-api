package security

import (
	"net/http"

	"github.com/ryankurte/go-csp"

	"github.com/ryankurte/go-api/lib/options"
)

// CSP builds a Content Security Policy (CSP) handler around the provided handler
func CSP(h http.Handler, o *options.Base) http.Handler {
	cspConfig := csp.CSP{
		ReportOnly:  o.CSP.ReportOnly,
		ChildSrc:    o.CSP.ChildSrc,
		ConnectSrc:  o.CSP.ConnectSrc,
		DefaultSrc:  o.CSP.DefaultSrc,
		FontSrc:     o.CSP.FontSrc,
		FrameSrc:    o.CSP.FrameSrc,
		ImgSrc:      o.CSP.ImgSrc,
		ManifestSrc: o.CSP.ManifestSrc,
		MediaSrc:    o.CSP.MediaSrc,
		ObjectSrc:   o.CSP.ObjectSrc,
		ScriptSrc:   o.CSP.ScriptSrc,
		StyleSrc:    o.CSP.StyleSrc,
		WorkerSrc:   o.CSP.WorkerSrc,
		ReportTo:    o.CSP.ReportTo,
	}

	// TODO: parse options to configure CSP

	return cspConfig.Handler(h)
}
