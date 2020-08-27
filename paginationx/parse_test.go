package paginationx

import (
	"net/url"
	"testing"
)

func TestErrorx(t *testing.T) {

	t.Run("Passing query with page=0 and limit=0", func(t *testing.T) {
		query := url.Values{
			"page":  []string{"0"},
			"limit": []string{"0"},
		}

		off, lim := Parse(query)

		if off != 0 {
			t.Errorf("Returned wrong offset expected %v, got %v", 0, off)
		}

		if lim != 10 {
			t.Errorf("Returned wrong limit expected %v, got %v", 10, lim)
		}
	})

	t.Run("Passing query with page=0 and limit>20", func(t *testing.T) {
		query := url.Values{
			"page":  []string{"0"},
			"limit": []string{"40"},
		}

		off, lim := Parse(query)

		if off != 0 {
			t.Errorf("Returned wrong offset expected %v, got %v", 0, off)
		}

		if lim != 10 {
			t.Errorf("Returned wrong limit expected %v, got %v", 10, lim)
		}
	})

	t.Run("Passing query with page=3 and limit=5", func(t *testing.T) {
		query := url.Values{
			"page":  []string{"3"},
			"limit": []string{"5"},
		}

		off, lim := Parse(query)

		if off != 10 {
			t.Errorf("Returned wrong offset expected %v, got %v", 10, off)
		}

		if lim != 5 {
			t.Errorf("Returned wrong limit expected %v, got %v", 5, lim)
		}
	})

}
