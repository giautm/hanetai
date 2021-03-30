package webhook

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	mLatencyMs = stats.Float64("latency", "The latency in milliseconds", "ms")
	mFaces     = stats.Int64("faces_detected", "The number of faces detected", "1")

	keyDeviceID   = tag.MustNewKey("giautm.dev/hanetai/device-id")
	keyPersonType = tag.MustNewKey("giautm.dev/hanetai/person-type")
	keyPlaceID    = tag.MustNewKey("giautm.dev/hanetai/place-id")
)

func EnableViews() error {
	latencyView := &view.View{
		Name:        "hanet/latency_webhook",
		Measure:     mLatencyMs,
		Description: "The distribution of the latencies",
		TagKeys:     []tag.Key{keyPlaceID, keyDeviceID},
		Aggregation: view.Distribution(0, 25, 100, 200, 400, 800, 10000),
	}

	facesDetectedCountView := &view.View{
		Name:        "hanet/faces_detected",
		Measure:     mFaces,
		Description: "The number of faces detected",
		TagKeys:     []tag.Key{keyPlaceID, keyDeviceID, keyPersonType},
		// Notice that the measure "mLatencyMs" is the same as
		// latencyView's but here the aggregation is a count aggregation
		// while the latencyView has a distribution aggregation.
		Aggregation: view.Count(),
	}

	// Ensure that they are registered so
	// that measurements won't be dropped.
	return view.Register(latencyView, facesDetectedCountView)
}

var mapPersonTypes = map[string]string{
	"0": "Employee",
	"1": "Customer",
	"2": "Stranger",
}

func ReportStats(fn WebhookFn) WebhookFn {
	return func(ctx context.Context, data *Webhook) (err error) {
		m := []tag.Mutator{}
		ms := []stats.Measurement{}

		if data.DeviceData != nil {
			m = append(m, tag.Upsert(keyDeviceID, data.DeviceID))
		}
		if data.PlaceData != nil {
			m = append(m, tag.Upsert(keyPlaceID, strconv.Itoa(data.PlaceID.Int())))
		}
		if data.PersonData != nil {
			if t, ok := mapPersonTypes[data.PersonType]; ok {
				m = append(m, tag.Upsert(keyPersonType, t))
			} else {
				m = append(m, tag.Upsert(keyPersonType, fmt.Sprintf("unknown: %s", data.PersonType)))
			}
		}

		ctx, err = tag.New(ctx, m...)
		if err != nil {
			return err
		}
		defer func() {
			stats.Record(ctx, ms...)
		}()

		if t := data.Time; t > 0 {
			clientTime := time.Unix(0, int64(data.Time)*int64(time.Millisecond))
			ms = append(ms, mLatencyMs.M((float64(time.Since(clientTime).Milliseconds()))))
		}
		if data.PersonData != nil {
			ms = append(ms, mFaces.M(1))
		}

		return fn(ctx, data)
	}
}
