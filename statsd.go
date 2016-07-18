package statsd

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gliderlabs/logspout/router"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/kr/logfmt"
)

func init() {
	router.AdapterFactories.Register(NewStatsdAdapter, "statsd")
}

func NewInfluxHandler(influxUrl, username, password string) (router.LogHandler, error) {
	prefix := ""
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxUrl,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	interval := time.Second * 2 // aggregate stats and flush every 2 seconds

	return &InfluxHandler{
		client: c,
	}, nil
}

type InfluxHandler struct {
	counter uint64
	client  *client.Client
}

type Metric struct {
	Metric string // Name of the metric
	Value  int64
	Type   string
}

func (a *InfluxHandler) HandleLine(message *router.Message) bool {
	if message.Container.Name == "/logspout" {
		return true
	}
	atomic.AddUint64(&a.counter, 1)
	// log.Println(atomic.LoadUint64(&a.counter), "source:", message.Source, "cname:", message.Container.Name, "mdata:", message.Data)

	m := &Metric{}
	if err := logfmt.Unmarshal([]byte(message.Data), m); err != nil {
		// log.Println("not in logfmt format, skipping")
		return true
	}
	// log.Println("metric:", *m)
	if m.Metric == "" {
		// log.Println("not a metric, skipping")
		return true
	}
	if m.Metric != "" {
		switch m.Type {
			// TODO: implement the influx stuff
	// 	case "count":
	// 		a.statsBuffer.Incr(m.Metric, m.Value)
	// 	case "gauge":
	// 		a.statsBuffer.Gauge(m.Metric, m.Value)
	// 	}
	// 	 // Create a new point batch
    // bp, err := client.NewBatchPoints(client.BatchPointsConfig{
    //     Database:  MyDB,
    //     Precision: "s",
    // })

    // if err != nil {
    //     log.Fatalln("Error: ", err)
    // }

    // // Create a point and add to batch
    // tags := map[string]string{"cpu": "cpu-total"}
    // fields := map[string]interface{}{
    //     "idle":   10.1,
    //     "system": 53.3,
    //     "user":   46.6,
    // }
    // pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())

    // if err != nil {
    //     log.Fatalln("Error: ", err)
    // }

    // bp.AddPoint(pt)

    // // Write the batch
    // c.Write(bp)
	}
	return true
}
