package dockerfile

func isCmd(cmd []byte) (string, bool) {
	if len(cmd) == 0 {
		return "", false
	}

	for _, c := range cmd {
		if c < 'A' || c > 'Z' {
			return "", false
		}
	}

	return string(cmd), true
}

func parseKV(b []byte) map[string]string {
	var (
		l   int
		key string
		q   bool
		out = make(map[string]string)
	)

	for i, c := range b {
		switch c {
		case ' ':
			if q {
				continue
			}
			if len(key) > 0 {
				out[key] = string(b[l:i])
				key = ""
			}
			l = i + 1
			// ????

		case '=':
			key = string(b[l:i])
			l = i + 1

		case '"':
			if q {
				out[key] = string(b[l:i])
				q = false
				key = ""
			} else {
				q = true
			}
			l = i + 1

		}
	}

	if len(key) > 0 {
		out[key] = string(b[l:])
	}

	return out
}
