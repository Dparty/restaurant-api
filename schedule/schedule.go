package schedule

import (
	"github.com/go-co-op/gocron/v2"
)

var s, _ = gocron.NewScheduler()

func init() {
	// s.NewJob(
	// 	gocron.CronJob(),
	// 	gocron.NewTask(
	// 		func(a string, b int) {
	// 			fmt.Println(a, b)
	// 		},
	// 		"hello",
	// 		1,
	// 	),
	// )
	// s.Start()
}
