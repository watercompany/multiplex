package worker_test

import (
	"testing"

	"github.com/watercompany/multiplex/worker"
)

func TestParseGraphData(t *testing.T) {
	tests := []struct {
		scenario  string
		data      string
		wantPhase int
		wantHours float64
	}{
		{
			scenario:  "phase 1",
			data:      "Time for phase 1 = 17274.691 seconds. CPU (182.480%) Tue Jun  8 21:11:22 2021",
			wantPhase: 1,
			wantHours: 4.8, // 17274.69/3600
		},
		{
			scenario:  "phase 2",
			data:      "Time for phase 2 = 10033.405 seconds. CPU (79.210%) Tue Jun  8 23:58:36 2021",
			wantPhase: 2,
			wantHours: 2.79, // 10033.41/3600
		},
		{
			scenario:  "phase 3",
			data:      "Time for phase 3 = 14076.263 seconds. CPU (86.670%) Wed Jun  9 03:53:12 2021",
			wantPhase: 3,
			wantHours: 3.91, // 14076.26/3600
		},
		{
			scenario:  "phase 4",
			data:      "Time for phase 4 = 1348.374 seconds. CPU (95.540%) Wed Jun  9 04:15:40 2021",
			wantPhase: 4,
			wantHours: 0.37, // 1348.37/3600
		},
	}
	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			phase, hours, err := worker.ParseGraphData(tc.data)
			if err != nil {
				t.Errorf("want no error, got %v", err)
			}

			if phase != tc.wantPhase {
				t.Errorf("want phase %v, got %v", tc.wantPhase, phase)
			}

			if hours != tc.wantHours {
				t.Errorf("want hours %v, got %v", tc.wantHours, hours)
			}
		})
	}
}
