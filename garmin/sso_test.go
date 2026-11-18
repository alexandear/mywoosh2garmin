package garmin

import (
	"regexp"
	"testing"
)

func TestSSORegexps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		re   *regexp.Regexp
		in   string
		want string // empty means no match
	}{
		{
			name: "csrf",
			re:   csrfRe,
			in:   `<input type="hidden" name="_csrf" value="3F2A-91BC" />`,
			want: "3F2A-91BC",
		},
		{
			name: "csrf with newline between attributes",
			re:   csrfRe,
			in:   "<input name=\"_csrf\"\n\tvalue=\"abc\">",
			want: "abc",
		},
		{
			name: "csrf missing",
			re:   csrfRe,
			in:   `<input type="hidden" name="other" value="abc" />`,
		},
		{
			name: "title success",
			re:   titleRe,
			in:   `<html><head><title>Success</title></head></html>`,
			want: "Success",
		},
		{
			name: "title MFA",
			re:   titleRe,
			in:   `<title>GARMIN > MFA Challenge</title>`,
			want: "GARMIN > MFA Challenge",
		},
		{
			name: "title missing",
			re:   titleRe,
			in:   `<html><body>no title</body></html>`,
		},
		{
			name: "ticket",
			re:   ticketRe,
			in:   `var response_url = "https://sso.garmin.com/sso/embed?ticket=ST-0123456-abcDEF-cas";`,
			want: "ST-0123456-abcDEF-cas",
		},
		{
			name: "ticket missing",
			re:   ticketRe,
			in:   `var response_url = "https://sso.garmin.com/sso/embed";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := tt.re.FindStringSubmatch(tt.in)
			var got string
			if m != nil {
				got = m[1]
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
