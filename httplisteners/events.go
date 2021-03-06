package httplisteners

import (
	"net/http"
	"strconv"
	"time"

	"github.com/serverless/event-gateway/metrics"
	"github.com/serverless/event-gateway/router"
	"github.com/serverless/event-gateway/targetcache"
)

// StartEventsAPI creates a new gateway endpoint and listens for requests.
func StartEventsAPI(conf Config) {
	targetCache := targetcache.New("/serverless-event-gateway", conf.KV, conf.Log)
	router := router.New(targetCache, metrics.DroppedPubSubEvents, conf.Log)
	router.StartWorkers()
	ev := &http.Server{
		Addr:         ":" + strconv.Itoa(int(conf.Port)),
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	h := handler{
		Conf:        conf,
		HTTPHandler: ev,
	}

	go func() {
		conf.ShutdownGuard.Add(1)
		h.listen()
		router.Drain()
		conf.ShutdownGuard.Done()
	}()
}
