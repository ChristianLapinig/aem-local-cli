package internal

import (
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/internal"
)

func TestGreeting(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "Christian",
			Expected: "Hello, Christian!",
		},
		{
			Name:     "John",
			Expected: "Hello, John!",
		},
		{
			Name:     "",
			Expected: "Hello, World!",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			actual := internal.Greeting(test.Name)
			if actual != test.Expected {
				t.Errorf("FAILED: got %s, expected %s", actual, test.Expected)
			}
		})
	}
}
