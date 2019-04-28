package utils

func ArgStringD(arg, def string) string {
	v := arg
	if arg == "" {
		v = def
	}
	return v
}
