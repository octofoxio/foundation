/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestEnvString(t *testing.T) {
	var x = EnvString("K", "x")
	assert.Equal(t, x, "x")
}

var scenario = []string{
	"Mon 10:30-12:00",
	"Wed 05:30-20:00",
	"Mon 12:30-14:30",
	"Sat 9:00-10:30",
}

var layout = "2 Mon 15:04"

func mustParse(v string) time.Time {
	d, _ := time.Parse(layout, v)
	return d
}
func mapDay(day string) int {
	for i, d := range []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sun", "Sat"} {
		if d == day {
			return i
		}
	}
	panic("invalid day")
}
func TestSomething(t *testing.T) {
	starts := make([]int, 0)
	ends := make([]int, 0)

	for _, line := range scenario {
		day := strings.Split(line, " ")[0]
		periods := strings.Split(line, " ")[1]
		s := strings.Split(periods, "-")[0]
		e := strings.Split(periods, "-")[1]
		starts = append(starts, int(mustParse(fmt.Sprintf("%d %s %s", mapDay(day), day, s)).Unix()))
		ends = append(ends, int(mustParse(fmt.Sprintf("%d %s %s", mapDay(day), day, e)).Unix()))
	}

	starts = append(starts, int(mustParse(fmt.Sprintf("%d %s %s", mapDay("Sat"), "Sat", "24:00")).Unix()))
	fmt.Println(starts)
	fmt.Println(ends)

	restTime := 0
	for i, end := range ends {
		restFrom := time.Unix(int64(end), 0)
		restTo := time.Unix(int64(starts[i+1]), 0)
		restDuration := int(restTo.Sub(restFrom).Minutes())
		if restTime < restDuration {
			restTime = restDuration
		}
	}
	fmt.Println(restTime)

}
