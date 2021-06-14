package worker_test

import (
	"testing"

	"github.com/watercompany/multiplex/worker"
)

func TestParseGraphData(t *testing.T) {
	tests := []struct {
		scenario  string
		data      string
		taskName  string
		wantPhase int
		wantHours float64
	}{
		{
			scenario:  "phase 1",
			data:      "Time for phase 1 = 17274.691 seconds. CPU (182.480%) Tue Jun  8 21:11:22 2021",
			taskName:  "pos",
			wantPhase: 1,
			wantHours: 4.8, // 17274.69/3600
		},
		{
			scenario:  "phase 2",
			data:      "Time for phase 2 = 10033.405 seconds. CPU (79.210%) Tue Jun  8 23:58:36 2021",
			taskName:  "pos",
			wantPhase: 2,
			wantHours: 2.79, // 10033.41/3600
		},
		{
			scenario:  "phase 3",
			data:      "Time for phase 3 = 14076.263 seconds. CPU (86.670%) Wed Jun  9 03:53:12 2021",
			taskName:  "pos",
			wantPhase: 3,
			wantHours: 3.91, // 14076.26/3600
		},
		{
			scenario:  "phase 4",
			data:      "Time for phase 4 = 1348.374 seconds. CPU (95.540%) Wed Jun  9 04:15:40 2021",
			taskName:  "pos",
			wantPhase: 4,
			wantHours: 0.37, // 1348.37/3600
		},
		{
			scenario:  "posv2 phase 1",
			data:      "Phase 1 took 3805.35 sec",
			taskName:  "posv2",
			wantPhase: 1,
			wantHours: 1.06, // 3805.35/3600
		},
		{
			scenario:  "posv2 phase 2",
			data:      "Phase 2 took 1924.45 sec",
			taskName:  "posv2",
			wantPhase: 2,
			wantHours: 0.53, // 1924.45/3600
		},
		{
			scenario:  "posv2 phase 3",
			data:      "Phase 3 took 2272.92 sec, wrote 21873560739 entries to final plot",
			taskName:  "posv2",
			wantPhase: 3,
			wantHours: 0.63, // 2272.92/3600
		},
		{
			scenario:  "posv2 phase 4",
			data:      "Phase 4 took 116.866 sec, final plot size is 108813626029 bytes",
			taskName:  "posv2",
			wantPhase: 4,
			wantHours: 0.03, // 116.866/3600
		},
	}
	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			phase, hours, err := worker.ParseGraphData(tc.data, tc.taskName)
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
