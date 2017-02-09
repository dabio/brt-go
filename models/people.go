package models

import "fmt"

// Person contains the data structure of a single person.
type Person struct {
	name  string
	email string
}

func (p *Person) CN() string {
	return fmt.Sprintf("%s:mailto:%s", p.name, p.email)
}
