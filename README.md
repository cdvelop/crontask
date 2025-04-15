# crontask

Scheduled task manager in Go with support for configuration in code and YAML.


// https://crontab.guru/
// mon:1,th:2, we:3,th:4,fr:5,sat:6,sun:0

const (
	once_year     = "0 7 15 1 *"  // A las 07:00 del día del mes 15 de enero
	two_days_week = "0 7 * * 1,4" //2 veces a la semana lunes y jueves

	// once_year = "*/2 * * * *" // test 3 min año
	// two_days_week = "*/1 * * * *" //test cada 1 min
)

este proyecto no seria posible
- github.com/mileusna/crontab