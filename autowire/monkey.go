package autowire

// monkey function

var mf func(interface{}, string)

func RegisterMonkeyFunction(f func(interface{}, string)) {
	mf = f
}
