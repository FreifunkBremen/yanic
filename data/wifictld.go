package data

// WifiCTLDNodeinfo struct (defined in: https://chaos.expert/genofire/wifictld/blob/preview/wireless/wifictld/src/ubus_service.c - ubus_get_config)
type WifiCTLDNodeinfo struct {
	Verbose                uint8  `json:"verbose,omitempty"`
	ClientTryThreashold    uint32 `json:"client_try_threashold,omitempty"`
	ClientSignalThreashold int32  `json:"client_signal_threashold,omitempty"`
	ClientCleanEvery       uint32 `json:"client_clean_every,omitempty"`
	ClientCleanOlderThen   uint32 `json:"client_clean_older_then,omitempty"`
	ClientCleanAuthed      uint8  `json:"client_clean_authed,omitempty"`
	ClientForce            uint8  `json:"client_force,omitempty"`
	ClientForceProbe       uint8  `json:"client_force_probe,omitempty"`
	ClientProbeSteering    uint8  `json:"client_probe_steering,omitempty"`
	ClientProbeLearning    uint8  `json:"client_probe_learning,omitempty"`
}

// WifiCTLDStatistics struct (defined in: https://chaos.expert/genofire/wifictld/blob/preview/wireless/respondd-module-wifictld/src/respondd.c - respondd_provider_statistics)
type WifiCTLDStatistics struct {
	Total           uint32 `json:"total,omitempty"`
	Client24        uint32 `json:"client24,omitempty"`
	Client5         uint32 `json:"client5,omitempty"`
	Authed          uint32 `json:"authed,omitempty"`
	Connected       uint32 `json:"connected,omitempty"`
	HighestTryProbe uint32 `json:"highest_try_probe,omitempty"`
	HighestTryAuth  uint32 `json:"highest_try_auth,omitempty"`
}
