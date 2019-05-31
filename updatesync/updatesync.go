// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package updatesync

import (
	"time"

	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// UpdateSync is an advanced function that allows to retry an failed update operation due to an
// updateTime mismatch (basically an update occurred between the time you got a copy of an object
// and the time you try to update it.
//
// It gets a manipulator and context to perform the updatea dnd also a function that returns an identifiable.
// This function must be used to apply the desired update. IT MUST NOT ALLOCATE A NEW OBJECT.
// For perf reasons, UpdateSync does not enforce that the return of your function must be a pointer.
//
// UpdateSync gets the Identifiable you want to update using your function and tries to update
// it. If it fails because of a desync, it will retrieve the latest copy of your object (modifying it)
// then will apply your function again (so you can reaply your desired changes) and try again until
// it succeeds or until the max number of tries is reached.
//
// Example:
//
//      // If we have this object unit in memory
//      o := &Object{
//          Name: "Hello World",
//          Description: "This is the original description"
//      }
//
//      // Then if we want to change the description, even if an update happened, we do:
//      err := UpdateSync(ctx.Background(), m, nil, func(obj elemental.Identifiable) {
//          obj.(*Object).Description = "This is the description I want!"
//      }
//
// Please keep in mind you have to be very careful with this function. You may still put the target object
// in a very weird state, if you override attributes that can be managed by a different source of truth.
func UpdateSync(
	mctx manipulate.Context,
	m manipulate.Manipulator,
	obj elemental.Identifiable,
	updateFunc func(elemental.Identifiable),
) error {

	deadline, hasDeadline := mctx.Context().Deadline()

	for {

		updateFunc(obj)

		err := m.Update(mctx, obj)
		if err == nil {
			return nil
		}

		if hasDeadline && deadline.Before(time.Now()) {
			return err
		}

		// If the error is not a validation error for read only update time, we return the error.
		if !elemental.IsValidationError(err, "Read Only Error", "updateTime") {
			return err
		}

		// Otherwise we retrieve the latest copy of the object.
		if err = m.Retrieve(mctx, obj); err != nil {
			return err
		}

		// Then the loop will run again trying to update the object by applying the updateFunc on it.
	}
}
