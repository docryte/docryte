package types

import "fmt"

type ContactRequest struct {
	Name    string
	Contact string
}

func (cr *ContactRequest) String() string {
	return fmt.Sprintf("Имя: %s\nСвязь: %s", cr.Name, cr.Contact)
}
