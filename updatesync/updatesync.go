package updatesync

import (
	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"
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
//      err := UpdateSync(m, nil, 5, func() elemental.Identifiable {
//          o.Description = "This is the description I want!" // USE OUTER CONTEXT TO GET O
//          return o
//      }
//
// Please keep in mind you have to be very careful with this function. You may still put the target object
// in a very weird state, if you override attributes that can be managed by a different source of truth.
func UpdateSync(
	m manipulate.Manipulator,
	mctx *manipulate.Context,
	maxTry int,
	updateFunc func() elemental.Identifiable,
) error {

	var try int

	for {

		try++

		obj := updateFunc()

		err := manipulate.Retry(func() error { return m.Update(mctx, obj) }, nil, 10)
		if err == nil {
			return nil
		}

		// If we reach the maximum number of try, we return the error.
		if try >= maxTry {
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
