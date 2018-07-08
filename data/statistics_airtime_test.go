package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrequencyName(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("11a", WirelessAirtime{Frequency: 5000}.FrequencyName())
	assert.Equal("11g", WirelessAirtime{Frequency: 4999}.FrequencyName())
	assert.Equal("11g", WirelessAirtime{Frequency: 2412}.FrequencyName())
}

func TestUtilization(t *testing.T) {
	assert := assert.New(t)

	t1 := &WirelessAirtime{
		ActiveTime: 20,
		BusyTime:   5,
		TxTime:     5,
		RxTime:     0,
	}
	t2 := &WirelessAirtime{
		ActiveTime: 120,
		BusyTime:   50,
		TxTime:     25,
		RxTime:     15,
	}

	t1.setUtilization(t1)
	assert.NotZero(t1.ChanUtil)
	assert.EqualValues(0, t2.ChanUtil)
	assert.EqualValues(0, t2.RxUtil)
	assert.EqualValues(0, t2.TxUtil)

	t2.setUtilization(t1)
	assert.NotZero(t2.ChanUtil)
	assert.EqualValues(45, t2.ChanUtil)
	assert.EqualValues(15, t2.RxUtil)
	assert.EqualValues(20, t2.TxUtil)
}

func TestWirelessStatistics(t *testing.T) {
	assert := assert.New(t)

	// This is the current value
	stats := WirelessStatistics([]*WirelessAirtime{{
		Frequency:  2400,
		ActiveTime: 10,
		BusyTime:   4,
		RxTime:     3,
	}})

	// previous value: Different Frequency, should not change anything
	stats.SetUtilization([]*WirelessAirtime{{
		Frequency:  5000,
		ActiveTime: 5,
		RxTime:     1,
	}})
	assert.Zero(0, stats[0].ChanUtil)
	assert.Zero(0, stats[0].RxUtil)
	assert.Zero(0, stats[0].TxUtil)

	// previous value: Same Frequency, should set the utilization
	stats.SetUtilization([]*WirelessAirtime{{
		Frequency:  2400,
		ActiveTime: 5,
		RxTime:     2,
		BusyTime:   1,
	}})
	assert.EqualValues(60, stats[0].ChanUtil)
	assert.EqualValues(20, stats[0].RxUtil)
	assert.EqualValues(0, stats[0].TxUtil)
}
