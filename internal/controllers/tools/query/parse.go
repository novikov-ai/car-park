package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func ParseQueryParamsIdStartEndDates(c *gin.Context) (int64, *time.Time, *time.Time) {
	query := c.Request.URL.Query()
	id := query.Get("vehicle")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	var startValue *time.Time
	start := query.Get("start")
	if start == "" {
		zeroUnix := time.Unix(0, 0)
		startValue = &zeroUnix
	} else {
		startValue = parseQueryTimeParam(start)
	}

	var endValue *time.Time
	end := query.Get("end")
	if end == "" {
		timeNow := time.Now()
		endValue = &timeNow
	} else {
		endValue = parseQueryTimeParam(end)
	}

	return idValue, startValue, endValue
}

func ParseQueryParams(c *gin.Context) (int64, *time.Time, *time.Time) {
	query := c.Request.URL.Query()
	id := query.Get("id")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	var startValue *time.Time
	start := query.Get("start")
	if start == "" {
		zeroUnix := time.Unix(0, 0)
		startValue = &zeroUnix
	} else {
		startValue = parseQueryTimeParam(start)
	}

	var endValue *time.Time
	end := query.Get("end")
	if end == "" {
		timeNow := time.Now()
		endValue = &timeNow
	} else {
		endValue = parseQueryTimeParam(end)
	}

	return idValue, startValue, endValue
}

func ParseQueryParamsIdStartEndUnix(c *gin.Context) (int64, int64, int64) {
	query := c.Request.URL.Query()
	id := query.Get("vehicle")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	start := query.Get("start")
	startUnix, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		startUnix = 0
	}

	startLocalTime := time.UnixMilli(startUnix)
	startUTC := startLocalTime.UTC()

	end := query.Get("end")
	endUnix, err := strconv.ParseInt(end, 10, 64)

	endLocalTime := time.UnixMilli(endUnix)
	endUTC := endLocalTime.UTC()

	if err != nil || end == "" {
		endUTC = time.Now().UTC()
	}

	return idValue, startUTC.Unix(), endUTC.Unix()
}

func GetVehicleStartEndTimeUnix(c *gin.Context) (int64, int64, int64) {
	query := c.Request.URL.Query()
	id := query.Get("vehicle")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	start := query.Get("start")
	startUnix := parseQueryTimeParam(start).Unix()

	startLocalTime := time.Unix(startUnix, 0)
	startUTC := startLocalTime.UTC()

	end := query.Get("end")
	endUnix := parseQueryTimeParam(end).Unix()

	endLocalTime := time.Unix(endUnix, 0)
	endUTC := endLocalTime.UTC()

	if err != nil || end == "" {
		endUTC = time.Now().UTC()
	}

	return idValue, startUTC.Unix(), endUTC.Unix()
}

func parseQueryTimeParam(value string) *time.Time {
	parsed, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", value))
	if err != nil {
		return nil
	}

	return &parsed
}
