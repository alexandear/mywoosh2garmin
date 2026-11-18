package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// emailEntry is an entry that shows the email partially masked, like a password entry
// with a reveal button. The full email is shown while revealed or while editing.
type emailEntry struct {
	widget.Entry

	email    string
	revealed bool
	focused  bool
	toggle   *widget.Button
}

func newEmailEntry() *emailEntry {
	e := &emailEntry{}
	e.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	e.ExtendBaseWidget(e)
	e.OnChanged = func(text string) {
		if e.showsEmail() {
			e.email = text
		}
	}
	e.toggle = widget.NewButtonWithIcon("", theme.VisibilityOffIcon(), func() {
		e.revealed = !e.revealed
		e.refresh()
	})
	e.toggle.Importance = widget.LowImportance
	e.ActionItem = e.toggle

	return e
}

// Email returns the unmasked email.
func (e *emailEntry) Email() string {
	return e.email
}

func (e *emailEntry) SetEmail(email string) {
	e.email = email
	e.refresh()
}

func (e *emailEntry) FocusGained() {
	e.focused = true
	e.refresh()
	e.Entry.FocusGained()
}

func (e *emailEntry) FocusLost() {
	e.focused = false
	e.refresh()
	e.Entry.FocusLost()
}

func (e *emailEntry) showsEmail() bool {
	return e.revealed || e.focused
}

func (e *emailEntry) refresh() {
	if e.revealed {
		e.toggle.SetIcon(theme.VisibilityIcon())
	} else {
		e.toggle.SetIcon(theme.VisibilityOffIcon())
	}

	if e.showsEmail() {
		e.SetText(e.email)
	} else {
		e.SetText(maskEmail(e.email))
	}
}

// maskEmail keeps the first characters of the local part and the domain name, and the top-level domain,
// e.g. "john.doe@example.com" becomes "jo••••••@ex•••••.com".
func maskEmail(email string) string {
	local, domain, hasDomain := strings.Cut(email, "@")
	if !hasDomain {
		return maskPart(local)
	}

	name, tld, hasTLD := strings.CutLast(domain, ".")
	if !hasTLD {
		return maskPart(local) + "@" + maskPart(domain)
	}

	return maskPart(local) + "@" + maskPart(name) + "." + tld
}

// maskPart keeps up to two leading characters, but never more than half, e.g. "example" becomes "ex•••••".
func maskPart(s string) string {
	runes := []rune(s)
	keep := min(2, len(runes)/2)

	return string(runes[:keep]) + strings.Repeat("•", len(runes)-keep)
}
