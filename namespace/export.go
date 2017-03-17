package namespace

import (
	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"
)

func Export(manipulator manipulate.Manipulator, namespace string) (elemental.IdentifiablesList, error) {

	nss := squallmodels.NamespacesList{}
	ns := &squallmodels.Namespace{}
	ns.Name = namespace

	mctx := manipulate.NewContext()
	mctx.Recursive = true
	mctx.Namespace = namespace

	if err := manipulator.RetrieveMany(mctx, &nss); err != nil {
		return nil, err
	}

	identifiablesChannel := make(chan elemental.IdentifiablesList)
	errorsChannel := make(chan error)
	identifiables := elemental.IdentifiablesList{}

	for _, n := range nss {
		identifiables = append(identifiables, n)
	}

	for _, identity := range exportNamespacesObjects {
		go func() {
			dest := squallmodels.ContentIdentifiableForIdentity(identity.Name)

			if err := manipulator.RetrieveMany(mctx, dest); err != nil {
				errorsChannel <- err
			}

			identifiablesChannel <- dest.List()
		}()

		select {
		case err := <-errorsChannel:
			return nil, err
		case ids := <-identifiablesChannel:
			identifiables = append(identifiables, ids...)
		}
	}

	return identifiables, nil
}
