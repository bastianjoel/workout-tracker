package helpers

import (
	"context"
	"time"
)

func LocalDate(ctx context.Context, t time.Time) string {
	return t.In(timezone(ctx)).Format(timeFormat)
}

func LocalTime(ctx context.Context, t time.Time) time.Time {
	return t.In(timezone(ctx))
}

func Timezone(ctx context.Context) string {
	return timezone(ctx).String()
}

func timezone(ctx context.Context) *time.Location {
	return CurrentUser(ctx).Timezone()
}

func InTimezone(t time.Time, tz string) time.Time {
	tzLoc, err := time.LoadLocation(tz)
	if err != nil {
		tzLoc = time.UTC
	}

	return t.In(tzLoc)
}

type TZ struct {
	Name        string
	Description string
}

func Timezones() []TZ { //nolint:funlen
	// Can I get this information from a library?
	return []TZ{
		{Name: "Etc/GMT+12", Description: "(GMT-12:00) International Date Line West"},
		{Name: "Pacific/Midway", Description: "(GMT-11:00) Midway Island, Samoa"},
		{Name: "Pacific/Honolulu", Description: "(GMT-10:00) Hawaii"},
		{Name: "US/Alaska", Description: "(GMT-09:00) Alaska"},
		{Name: "America/Los_Angeles", Description: "(GMT-08:00) Pacific Time (US & Canada) "},
		{Name: "America/Tijuana", Description: "(GMT-08:00) Tijuana, Baja California"},
		{Name: "US/Arizona", Description: "(GMT-07:00) Arizona"},
		{Name: "America/Chihuahua", Description: "(GMT-07:00) Chihuahua, La Paz, Mazatlan "},
		{Name: "US/Mountain", Description: "(GMT-07:00) Mountain Time (US & Canada)"},
		{Name: "America/Managua", Description: "(GMT-06:00) Central America"},
		{Name: "US/Central", Description: "(GMT-06:00) Central Time (US & Canada)"},
		{Name: "America/Mexico_City", Description: "(GMT-06:00) Guadalajara, Mexico City, Monterrey "},
		{Name: "Canada/Saskatchewan", Description: "(GMT-06:00) Saskatchewan"},
		{Name: "America/Bogota", Description: "(GMT-05:00) Bogota, Lima, Quito, Rio Branco "},
		{Name: "US/Eastern", Description: "(GMT-05:00) Eastern Time (US & Canada)"},
		{Name: "US/East-Indiana", Description: "(GMT-05:00) Indiana (East)"},
		{Name: "Canada/Atlantic", Description: "(GMT-04:00) Atlantic Time (Canada)"},
		{Name: "America/Caracas", Description: "(GMT-04:00) Caracas, La Paz"},
		{Name: "America/Manaus", Description: "(GMT-04:00) Manaus"},
		{Name: "America/Santiago", Description: "(GMT-04:00) Santiago"},
		{Name: "Canada/Newfoundland", Description: "(GMT-03:30) Newfoundland"},
		{Name: "America/Sao_Paulo", Description: "(GMT-03:00) Brasilia"},
		{Name: "America/Argentina/Buenos_Aires", Description: "(GMT-03:00) Buenos Aires, Georgetown "},
		{Name: "America/Godthab", Description: "(GMT-03:00) Greenland"},
		{Name: "America/Montevideo", Description: "(GMT-03:00) Montevideo"},
		{Name: "America/Noronha", Description: "(GMT-02:00) Mid-Atlantic"},
		{Name: "Atlantic/Cape_Verde", Description: "(GMT-01:00) Cape Verde Is."},
		{Name: "Atlantic/Azores", Description: "(GMT-01:00) Azores"},
		{Name: "Africa/Casablanca", Description: "(GMT+00:00) Casablanca, Monrovia, Reykjavik "},
		{Name: "Etc/Greenwich", Description: "(GMT+00:00) Greenwich Mean Time : Dublin, Edinburgh, Lisbon, London "},
		{Name: "Europe/Amsterdam", Description: "(GMT+01:00) Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna "},
		{Name: "Europe/Belgrade", Description: "(GMT+01:00) Belgrade, Bratislava, Budapest, Ljubljana, Prague "},
		{Name: "Europe/Brussels", Description: "(GMT+01:00) Brussels, Copenhagen, Madrid, Paris "},
		{Name: "Europe/Sarajevo", Description: "(GMT+01:00) Sarajevo, Skopje, Warsaw, Zagreb "},
		{Name: "Africa/Lagos", Description: "(GMT+01:00) West Central Africa"},
		{Name: "Asia/Amman", Description: "(GMT+02:00) Amman"},
		{Name: "Europe/Athens", Description: "(GMT+02:00) Athens, Bucharest, Istanbul"},
		{Name: "Asia/Beirut", Description: "(GMT+02:00) Beirut"},
		{Name: "Africa/Cairo", Description: "(GMT+02:00) Cairo"},
		{Name: "Africa/Harare", Description: "(GMT+02:00) Harare, Pretoria"},
		{Name: "Europe/Helsinki", Description: "(GMT+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius "},
		{Name: "Asia/Jerusalem", Description: "(GMT+02:00) Jerusalem"},
		{Name: "Europe/Minsk", Description: "(GMT+02:00) Minsk"},
		{Name: "Africa/Windhoek", Description: "(GMT+02:00) Windhoek"},
		{Name: "Asia/Kuwait", Description: "(GMT+03:00) Kuwait, Riyadh, Baghdad"},
		{Name: "Europe/Moscow", Description: "(GMT+03:00) Moscow, St. Petersburg, Volgograd "},
		{Name: "Africa/Nairobi", Description: "(GMT+03:00) Nairobi"},
		{Name: "Asia/Tbilisi", Description: "(GMT+03:00) Tbilisi"},
		{Name: "Asia/Tehran", Description: "(GMT+03:30) Tehran"},
		{Name: "Asia/Muscat", Description: "(GMT+04:00) Abu Dhabi, Muscat"},
		{Name: "Asia/Baku", Description: "(GMT+04:00) Baku"},
		{Name: "Asia/Yerevan", Description: "(GMT+04:00) Yerevan"},
		{Name: "Asia/Kabul", Description: "(GMT+04:30) Kabul"},
		{Name: "Asia/Yekaterinburg", Description: "(GMT+05:00) Yekaterinburg"},
		{Name: "Asia/Karachi", Description: "(GMT+05:00) Islamabad, Karachi, Tashkent"},
		{Name: "Asia/Calcutta", Description: "(GMT+05:30) Chennai, Kolkata, Mumbai, New Delhi, Sri Jayawardenapura"},
		{Name: "Asia/Katmandu", Description: "(GMT+05:45) Kathmandu"},
		{Name: "Asia/Almaty", Description: "(GMT+06:00) Almaty, Novosibirsk"},
		{Name: "Asia/Dhaka", Description: "(GMT+06:00) Astana, Dhaka"},
		{Name: "Asia/Rangoon", Description: "(GMT+06:30) Yangon (Rangoon)"},
		{Name: "Asia/Bangkok", Description: "(GMT+07:00) Bangkok, Hanoi, Jakarta"},
		{Name: "Asia/Krasnoyarsk", Description: "(GMT+07:00) Krasnoyarsk"},
		{Name: "Asia/Hong_Kong", Description: "(GMT+08:00) Beijing, Chongqing, Hong Kong, Urumqi "},
		{Name: "Asia/Kuala_Lumpur", Description: "(GMT+08:00) Kuala Lumpur, Singapore"},
		{Name: "Asia/Irkutsk", Description: "(GMT+08:00) Irkutsk, Ulaan Bataar"},
		{Name: "Australia/Perth", Description: "(GMT+08:00) Perth"},
		{Name: "Asia/Taipei", Description: "(GMT+08:00) Taipei"},
		{Name: "Asia/Tokyo", Description: "(GMT+09:00) Osaka, Sapporo, Tokyo"},
		{Name: "Asia/Seoul", Description: "(GMT+09:00) Seoul"},
		{Name: "Asia/Yakutsk", Description: "(GMT+09:00) Yakutsk"},
		{Name: "Australia/Adelaide", Description: "(GMT+09:30) Adelaide"},
		{Name: "Australia/Darwin", Description: "(GMT+09:30) Darwin"},
		{Name: "Australia/Brisbane", Description: "(GMT+10:00) Brisbane"},
		{Name: "Australia/Canberra", Description: "(GMT+10:00) Canberra, Melbourne, Sydney "},
		{Name: "Australia/Hobart", Description: "(GMT+10:00) Hobart"},
		{Name: "Pacific/Guam", Description: "(GMT+10:00) Guam, Port Moresby"},
		{Name: "Asia/Vladivostok", Description: "(GMT+10:00) Vladivostok"},
		{Name: "Asia/Magadan", Description: "(GMT+11:00) Magadan, Solomon Is., New Caledonia "},
		{Name: "Pacific/Auckland", Description: "(GMT+12:00) Auckland, Wellington"},
		{Name: "Pacific/Fiji", Description: "(GMT+12:00) Fiji, Kamchatka, Marshall Is. "},
		{Name: "Pacific/Tongatapu", Description: "(GMT+13:00) Nuku'alofa"},
	}
}
