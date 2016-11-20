package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilization(t *testing.T) {
	assert := assert.New(t)

	t1 := &WirelessAirtime{
		Active_time: 20,
		Busy_time:   0,
		Tx_time:     5,
		Rx_time:     0,
	}
	t2 := &WirelessAirtime{
		Active_time: 120,
		Busy_time:   10,
		Tx_time:     25,
		Rx_time:     15,
	}
	t3 := &WirelessAirtime{
		Active_time: 200,
		Busy_time:   40,
		Tx_time:     35,
		Rx_time:     15,
	}

	t1.SetUtilization(t2)
	assert.Zero(t1.ChanUtil)
	assert.Zero(t1.TxUtil)
	assert.Zero(t1.RxUtil)

	t2.SetUtilization(t1)
	assert.NotZero(t2.ChanUtil)
	assert.EqualValues(45, t2.ChanUtil)
	assert.EqualValues(20, t2.RxUtil)
	assert.EqualValues(15, t2.TxUtil)

	t3.SetUtilization(t2)
	assert.EqualValues(50, t3.ChanUtil)
	assert.EqualValues(12.5, t3.RxUtil)
	assert.EqualValues(0, t3.TxUtil)
}

func TestUtilizationStatistics(t *testing.T) {
	assert := assert.New(t)
	stats := WirelessStatistics{
		Airtime24: &WirelessAirtime{Active_time: 20},
		Airtime5:  &WirelessAirtime{Active_time: 20},
	}

	stats.SetUtilization(&WirelessStatistics{
		Airtime24: &WirelessAirtime{},
		Airtime5:  &WirelessAirtime{},
	})

	assert.Equal(20, int(stats.Airtime24.Active_time))
	assert.Equal(20, int(stats.Airtime5.Active_time))
}
