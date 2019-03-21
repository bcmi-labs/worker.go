package worker_test

import (
	"fmt"
	"time"

	"github.com/bcmi-labs/worker"
)

func Example() {
	fmt.Println("worker usage")
	mockJobs := mockJobs{}

	pool := worker.Pool{
		Jobs: &mockJobs,
	}

	pool.Run(func() {
		time.Sleep(1 * time.Second)
	})

	fmt.Println("Jobs running:", mockJobs.N)

	pool.Run(func() {
		time.Sleep(1 * time.Second)
	})

	fmt.Println("Jobs running:", mockJobs.N)

	time.Sleep(2 * time.Second)

	fmt.Println("Jobs running:", mockJobs.N)

	// Output: worker usage
	// Jobs running: 1
	// Jobs running: 2
	// Jobs running: 0
}

func ExamplePanic() {
	fmt.Println("panics are recovered")

	pool := worker.Pool{
		Recover: func(e interface{}) {
			fmt.Println(e)
		},
	}

	pool.Run(func() {
		panic("error")
	})

	time.Sleep(1 * time.Second)

	// Output: panics are recovered
	// error
}

type mockJobs struct {
	N int
}

func (m *mockJobs) Inc() { m.N++ }
func (m *mockJobs) Dec() { m.N-- }
