/*
* Copyright 2019 ARDUINO SA (http://www.arduino.cc/)
* This file is part of worker.go.
* Copyright (c) 2019
* Authors: Matteo "triex" Suppo
*
* This software is released under:
* first, the GNU General Public License, which covers the main part of
* worker.go
* The terms of this license can be found at:
* https://www.gnu.org/licenses/gpl-3.0.en.html
*
* next, the GNU Lesser General Public License that covers
* worker.go. The terms of this license can be found at:
* https://www.gnu.org/licenses/lgpl.html
*
* You can be released from the requirements of the above licenses by purchasing
* a commercial license. Buying such a license is mandatory if you want to modify or
* otherwise use the software for commercial activities involving the Arduino
* software without disclosing the source code of your own applications. To purchase
* a commercial license, send an email to license@arduino.cc.
*
 */

package worker

// Trace is an object that keep tracks of the duration of an action.
// When the object is created a timer starts, and whenever Mark is called
// the object will record the duration of the action described
type Trace interface {
	Mark(action string)
}

// Tracer is an object that can spawn traces. Scope is a generic string
// to assign to the trace in order to organize them
type Tracer interface {
	New(scope string) Trace
}
