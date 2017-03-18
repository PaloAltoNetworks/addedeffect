package namespaceutils

import (
	"fmt"

	"github.com/aporeto-inc/manipulate"
)

func Import(manipulator manipulate.Manipulator, namespace string, content map[string]interface{}, shouldClean bool) error {

	if _, ok := content["namespace"]; !ok {
		return fmt.Errorf("The given content should have a key namespace")
	}

	return nil
}
