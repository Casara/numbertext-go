package main

import "testing"

const testUSD = "USD"

// TestRunDispatch checks that run selects the right numbertext function
// for each action flag. The "no action flag set" case isn't covered
// here: it calls os.Exit directly (see run's default case), which would
// kill the test binary rather than return to it.
func TestRunDispatch(t *testing.T) {
	cases := []struct {
		name string
		opts options
		want string
	}{
		{"cardinal", options{lang: "en", cardinal: "21"}, "twenty-one"},
		{"ordinal", options{lang: "en", ordinal: "3"}, "third"},
		{"ordinal-number", options{lang: "en", ordinalNumber: "3"}, "3rd"},
		{"year", options{lang: "en", year: "1999"}, "nineteen ninety-nine"},
		{"currency", options{lang: "en", currency: testUSD, amount: "1"}, "one U.S. dollar"},
		{
			"money",
			options{lang: "en", currency: testUSD, amount: "2.5", money: true},
			"two and 50/100 U.S. dollars",
		},
		{"help", options{lang: "en", help: true}, ""}, // non-empty checked separately below
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := run(&tc.opts)
			if err != nil {
				t.Fatalf("run(%+v): %v", tc.opts, err)
			}
			if tc.name == "help" {
				if got == "" {
					t.Error("run(help) = \"\", want non-empty usage text")
				}

				return
			}
			if got != tc.want {
				t.Errorf("run(%+v) = %q, want %q", tc.opts, got, tc.want)
			}
		})
	}
}

// TestRunDispatchCurrencyBeforeMoneyPriority checks that -currency alone
// (without -money) takes the plain Currency path, not Money's.
func TestRunDispatchCurrencyBeforeMoneyPriority(t *testing.T) {
	got, err := run(&options{lang: "en", currency: testUSD, amount: "2.5"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if want := "two U.S. dollars and fifty cents"; got != want {
		t.Errorf("run(currency, no -money) = %q, want %q", got, want)
	}
}

// TestRunDispatchError checks that an error from the underlying
// numbertext call (e.g. an unknown language) propagates out of run
// rather than being swallowed.
func TestRunDispatchError(t *testing.T) {
	_, err := run(&options{lang: "xx-not-a-lang", cardinal: "1"})
	if err == nil {
		t.Error("run with an unknown language should return an error")
	}
}
