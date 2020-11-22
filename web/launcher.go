package web

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type ServerLauncher struct {
	s       *Server
	closers map[string]func() error
}

func (l *ServerLauncher) RunAutoTLS() error {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(l.s.cfg.Web.DomainNames[0:]...),
	}
	dir := l.cacheDir()
	l.s.l.Debug().Msgf("Using cache dir: %s", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		l.s.l.Warn().Msgf("warning: autocert.NewListener not using a cache: %v", err)
	} else {
		m.Cache = autocert.DirCache(dir)
	}

	go func() {
		httpServer := &http.Server{
			Addr:              l.s.cfg.Web.ListenAddress + ":80",
			Handler:           m.HTTPHandler(nil),
			ReadHeaderTimeout: 3 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      15 * time.Second,
			MaxHeaderBytes:    2048,
		}
		l.closers["http autocert server"] = httpServer.Close
		if err := httpServer.ListenAndServe(); err != nil {
			l.s.l.Warn().Err(err).Msg("AutoCert server error.")
		}
	}()

	return l.runWithManager(l.s.e, m, l.s.cfg.Web.ListenAddress)
}

func (l *ServerLauncher) runWithManager(r http.Handler, m *autocert.Manager, address string) error {
	s := &http.Server{
		Addr:              address + ":443",
		TLSConfig:         &tls.Config{GetCertificate: m.GetCertificate},
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      15 * time.Second,
		MaxHeaderBytes:    2048,
	}

	l.closers["https server"] = s.Close
	return s.ListenAndServeTLS("", "")
}

func (l *ServerLauncher) cacheDir() string {
	const base = "golang-autocert"
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(l.homeDir(), "Library", "Caches", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(v, base)
			}
		}
		// Worst case:
		return filepath.Join(l.homeDir(), base)
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(l.homeDir(), ".cache", base)
}

func (l *ServerLauncher) homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "/"
}

func (l *ServerLauncher) Close() error {
	for desc, closeFn := range l.closers {
		if err := closeFn(); err != nil {
			l.s.l.Error().Err(err).Msgf("Unable to close %s", desc)
		}
	}
	return nil
}
