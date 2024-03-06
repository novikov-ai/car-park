package qparser

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	idKey      = "id"
	vehicleKey = "vehicle"
	startKey   = "start"
	endKey     = "end"
)

var queryKeys = map[string]struct{}{
	idKey:      {},
	vehicleKey: {},
	startKey:   {},
	endKey:     {},
}

type ParsedValues struct {
	values map[string]string
}

func New(url *url.URL) (ParsedValues, error) {
	if url == nil {
		return ParsedValues{}, errors.New("url missed")
	}

	vv := make(map[string]string)
	query := url.Query()

	for key := range queryKeys {
		vv[key] = query.Get(key)
	}

	return ParsedValues{
		values: vv,
	}, nil
}

func (pv ParsedValues) GetID() int {
	return pv.getInt(idKey)
}

func (pv ParsedValues) GetVehicle() int64 {
	return int64(pv.getInt(vehicleKey))
}

func (pv ParsedValues) getInt(key string) int {
	value, err := strconv.Atoi(pv.values[key])
	if err != nil {
		return 0
	}

	return value
}

func (pv ParsedValues) GetStartTime() *time.Time {
	if parsed := parseQueryTimeParam(pv.values[startKey]); parsed != nil {
		return parsed
	}

	zeroUnix := time.Unix(0, 0)

	return &zeroUnix
}

func (pv ParsedValues) GetStartTimeUnix() int64 {
	unix, _ := strconv.ParseInt(pv.values[startKey], 10, 64)

	localTime := time.UnixMilli(unix)
	utc := localTime.UTC()

	return utc.Unix()
}

func (pv ParsedValues) GetEndTime() *time.Time {
	if parsed := parseQueryTimeParam(pv.values[endKey]); parsed != nil {
		return parsed
	}

	timeNow := time.Now()
	return &timeNow
}

func (pv ParsedValues) GetEndTimeUnix() int64 {
	unix, err := strconv.ParseInt(pv.values[endKey], 10, 64)

	localTime := time.UnixMilli(unix)
	utc := localTime.UTC()

	if err != nil {
		utc = time.Now().UTC()
	}

	return utc.Unix()
}

func parseQueryTimeParam(value string) *time.Time {
	if value == "" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", value))
	if err != nil {
		return nil
	}

	return &parsed
}
