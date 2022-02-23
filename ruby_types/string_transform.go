package ruby_types

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Subset of https://github.com/lyft/protoc-gen-star/blob/master/name.go,
// but without splitting on digits.

func upperCamelCase(s string) string { return transform(s, strings.Title, strings.Title, "") }

func split(ns string) (parts []string) {
	switch {
	case ns == "":
		return []string{""}
	case strings.LastIndex(ns, ".") >= 0:
		return strings.Split(ns, ".")
	case strings.LastIndex(ns, "_") > 0: // leading underscore does not count
		parts = strings.Split(ns, "_")
		if parts[0] == "" {
			parts[1] = "_" + parts[1]
			return parts[1:]
		}
		return
	default: // camelCase
		buf := &bytes.Buffer{}
		var capt, lodash bool
		for _, r := range ns {
			uc := unicode.IsUpper(r) || unicode.IsTitle(r)

			if r == '_' && buf.Len() == 0 && len(parts) == 0 {
				lodash = true
			}

			if uc && !capt && buf.Len() > 0 && !lodash { // new upper letter
				parts = append(parts, buf.String())
				buf.Reset()
			} else if !uc && capt && buf.Len() > 1 { // upper to lower
				if ss := buf.String(); len(ss) > 1 &&
					(len(ss) != 2 || ss[0] != '_') {
					pr, _ := utf8.DecodeLastRuneInString(ss)
					parts = append(parts, strings.TrimSuffix(ss, string(pr)))
					buf.Reset()
					buf.WriteRune(pr)
				}
			}

			capt = uc
			buf.WriteRune(r)
		}
		parts = append(parts, buf.String())
		return
	}
}

type stringTransformer func(string) string

func transform(s string, mod, first stringTransformer, sep string) string {
	parts := split(s)

	for i, p := range parts {
		if i == 0 {
			parts[i] = first(p)
		} else {
			parts[i] = mod(p)
		}
	}

	return strings.Join(parts, sep)
}
