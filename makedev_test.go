package main

import "testing"

func TestMakedev(t *testing.T) {
	var testdata = []struct {
		Major, Minor int64
		Result       int
	}{
		{1, 100, 356},
		{99, 23, 25367},
		{254, 3, 65027},
	}

	for i, tt := range testdata {
		if result := makedev(tt.Major, tt.Minor); result != tt.Result {
			t.Errorf("%d. makedev(%d, %d) == %d (expected: %d)\n", i, tt.Major, tt.Minor, result, tt.Result)
		}
	}
}
