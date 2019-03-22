/*
* Copyright 2018 ARDUINO SA (http://www.arduino.cc/)
* This file is part of [insert application name].
* Copyright (c) [insert year]
* Authors: [insert authors]
*
* This software is released under:
* The GNU General Public License, which covers the main part of
* [insert application name]
* The terms of this license can be found at:
* https://www.gnu.org/licenses/gpl-3.0.en.html
*
* You can be released from the requirements of the above licenses by purchasing
* a commercial license. Buying such a license is mandatory if you want to modify or
* otherwise use the software for commercial activities involving the Arduino
* software without disclosing the source code of your own applications. To purchase
* a commercial license, send an email to license@arduino.cc.
*
 */

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
