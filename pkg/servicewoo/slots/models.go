package slots

type Slots struct {
	DateSettingsDateSettingsSpecificDays []DateSettingsDateSettingsSpecificDay `json:"datesettings_datesettings_specific_days"`
	DateSettingsDateSettingsMinimum      string                                `json:"datesettings_datesettings_minimum"`
	DateSettingsDateSettingsMaximum      string                                `json:"datesettings_datesettings_maximum"`
	TimeSettingsTimeSettingsTimeSlots    []TimeSettingsTimeSettingsTimeSlot    `json:"timesettings_timesettings_timeslots"`
	TimeSettingsTimeSettingsCutoff       string                                `json:"timesettings_timesettings_cutoff"`
	HolidaysHolidaysHolidays             []HolidaysHolidaysHoliday             `json:"holidays_holidays_holidays"`
}

type DateSettingsDateSettingsSpecificDay struct {
	ID      string `json:"row_id"`
	Date    string `json:"date"`
	AltDate string `json:"alt_date"`
}

type TimeSettingsTimeSettingsTimeSlot struct {
	Duration  string   `json:"duration"`
	Frequency string   `json:"frequency"`
	TimeFrom  string   `json:"timefrom"`
	TimeTo    string   `json:"timeto"`
	Days      []string `json:"days"`
}

type HolidaysHolidaysHoliday struct {
	Date      string `json:"date"`
	AltDate   string `json:"alt_date"`
	DateTo    string `json:"date_to"`
	AltDateTo string `json:"alt_date_to"`
}
