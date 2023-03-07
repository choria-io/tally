package main

import (
	"context"
	"fmt"
	"github.com/choria-io/fisk"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/config"
	"github.com/choria-io/go-choria/lifecycle/tally"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// version is the release version to be set at compile time
	version = "development"

	port            int
	debug           bool
	tls             bool
	componentFilter string
	pidfile         string
	cfgfile         string
	prefix          string
	election        string
)

func main() {
	app := fisk.New("tally", "The Choria network observation tool")
	app.Version(version)
	app.Author("R.I.Pienaar <rip@devco.net>")
	app.Flag("config", "Configuration file to use").ExistingFileVar(&cfgfile)
	app.Flag("debug", "Enable debug logging").BoolVar(&debug)
	app.Flag("tls", "Use TLS when connecting to the middleware").Hidden().Default("true").BoolVar(&tls)
	app.Flag("component", "Component to tally").Default("*").StringVar(&componentFilter)
	app.Flag("port", "Port to listen on").Default("8080").IntVar(&port)
	app.Flag("prefix", "Prefix for statistic keys").Default("choria_tally").StringVar(&prefix)
	app.Flag("election", "Perform leader election in the named campaign").PlaceHolder("CAMPAIGN").StringVar(&election)

	app.MustParseWithUsage(os.Args[1:])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go interruptWatcher(ctx, cancel)

	if pidfile != "" {
		err := os.WriteFile(pidfile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
		fisk.FatalIfError(err, "Could not write pid file %s: %s", pidfile, err)
	}

	if cfgfile == "" {
		cfgfile = choria.UserConfig()
	}

	cfg, err := config.NewConfig(cfgfile)
	fisk.FatalIfError(err, "could not parse configuration: %s", err)

	if debug {
		cfg.LogLevel = "debug"
		cfg.LogFile = ""
	}

	if !tls {
		cfg.DisableTLS = true
	}

	fw, err := choria.NewWithConfig(cfg)
	fisk.FatalIfError(err, "could not set up Choria: %s", err)

	err = recordTally(ctx, fw)

	fisk.FatalIfError(err, "Could not run: %s", err)
}

func recordTally(ctx context.Context, fw *choria.Framework) error {
	log := fw.Logger("tally")
	log.Infof("Choria Lifecycle Tally version %s starting listening on port %d", version, port)

	conn, err := fw.NewConnector(ctx, fw.MiddlewareServers, fw.Certname(), log)
	if err != nil {
		return fmt.Errorf("cannot connect: %s", err)
	}

	recorder, err := tally.New(tally.Component(componentFilter), tally.Logger(log), tally.StatsPrefix(prefix), tally.Connection(conn), tally.Election(election))
	if err != nil {
		return fmt.Errorf("could not create recorder: %s", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	err = recorder.Run(ctx)
	if err != nil {
		return fmt.Errorf("recorder failed: %s", err)
	}

	return nil
}

func interruptWatcher(ctx context.Context, cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				cancel()
			}
		case <-ctx.Done():
			return
		}
	}
}
