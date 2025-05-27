package http

import "testing"

func TestAlphaRus_DisallowHyphenAndSpace(t *testing.T) {
	cases := []string{"Анна-Мария", "Жан Клод"}
	for _, s := range cases {
		if err := validate.Var(s, "alpharus"); err == nil {
			t.Errorf("%q should be invalid", s)
		}
	}
}
