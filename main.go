package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/alexeldeib/perf/biolatency"
	"github.com/alexeldeib/perf/iostat"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

var (
	prometheusPort       = ":9090"
	agentEndpointURI     = "localhost:6831"
	collectorEndpointURI = "http://localhost:14268/api/traces"
)

var (
	mBiolatency    = stats.Int64("biolatency", "block i/o latency in usecs", stats.UnitDimensionless)
	biolatencyView = &view.View{
		Name:        "biolatency",
		Description: "distribution of block i/o latency",
		Measure:     mBiolatency,
		// 1 usec to ~1 sec
		Aggregation: view.Distribution(1, math.Pow(2, 2), math.Pow(2, 4), math.Pow(2, 8), math.Pow(2, 10), math.Pow(2, 12), math.Pow(2, 14), math.Pow(2, 16), math.Pow(2, 18), math.Pow(2, 20)),
	}
)

func main() {
	// Setup logger
	logger, err := zap.NewDevelopment() // or NewProduction
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := instrument(sugar); err != nil {
		sugar.Fatal(err)
	}

	stats, err := iostat.New()
	if err != nil {
		sugar.Fatal(err)
	}

	biolatencyStats, err := biolatency.New(mBiolatency)
	if err != nil {
		sugar.Fatal(fmt.Sprintf("%#+v", err))
	}

	litter.Dump(stats)
	litter.Dump(biolatencyStats)
	time.Sleep(time.Second * 5)
}

type customMetricsExporter struct{}

func (ce *customMetricsExporter) ExportView(vd *view.Data) {
	log.Printf("vd.View: %+v\n%#v\n", vd.View, vd.Rows)
	for i, row := range vd.Rows {
		log.Printf("\tRow: %#d: %#v\n", i, row)
	}
	log.Printf("StartTime: %s EndTime: %s\n\n", vd.Start.Round(0), vd.End.Round(0))
}

func instrument(sugar *zap.SugaredLogger) error {
	view.SetReportingPeriod(60 * time.Second)

	if err := view.Register(biolatencyView); err != nil {
		sugar.Fatal(err, "failed to register biolatency view")
	}

	view.RegisterExporter(new(customMetricsExporter))

	// Stats exporter: Prometheus
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "aced",
	})

	if err != nil {
		errors.Wrapf(err, "failed to create the Prometheus stats exporter: %v")
	}

	view.RegisterExporter(pe)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		sugar.Fatal(http.ListenAndServe(prometheusPort, mux))
	}()

	// Trace exporter: Jaeger
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "aced",
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create the Jaeger exporter: %v")
	}

	trace.RegisterExporter(je)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}
