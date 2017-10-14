package mediator

import "testing"

func Test_ReCalculateScribeResponsibility(t *testing.T) {
	m := &Mediator{
		scribes:              make(map[string]string),
		scribeResponsibility: make(map[string]string),
	}

	t.Run("no scribes", func(t *testing.T) {
		// checking that if there are no scribes will do nothing
		m.reCalculateScribeResponsibility()
		if len(m.scribeResponsibility) != 0 {
			t.Errorf("expected len of 0 and got %d", len(m.scribeResponsibility))
		}
	})

	t.Run("single scribe", func(t *testing.T) {
		m.scribes["1"] = "scribe:1"
		m.reCalculateScribeResponsibility()
		if len(m.scribeResponsibility) != 1 {
			t.Errorf("expected len of 1 and got %d", len(m.scribeResponsibility))
		}

		for k := range m.scribeResponsibility {
			if k != "9" {
				t.Errorf("expected the key value to be '9' and got '%s'", k)
			}
		}
	})

	// two scribes
	t.Run("two scribes", func(t *testing.T) {
		m.scribes["2"] = "scribe:2"
		m.reCalculateScribeResponsibility()
		if len(m.scribeResponsibility) != 2 {
			t.Errorf("expected len of 2 and got %d", len(m.scribeResponsibility))
		}

		for k := range m.scribeResponsibility {
			if k != "r" && k != "9" {
				t.Errorf("expected the key value to be 'r' or '9' and got '%s'", k)
			}
		}
	})

}
