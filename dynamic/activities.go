package dynamic

import (
	"fmt"
)

type Activities struct{}

func (a *Activities) GetName() (string, error) {
	return "Temporal", nil
}

func (a *Activities) GetGreeting() (string, error) {
	return "Hello", nil
}

func (a *Activities) SayGreeting(greeting string, name string) (string, error) {
	result := fmt.Sprintf("Greeting: %s %s!\n", greeting, name)
	return result, nil
}
