package types

type GoatUtils struct {
	initialized bool
}

type GoatUtilsInterface interface {
	SetInitialized(bool)
	MustBeInitialized()
	IsInitialized() bool
}

func NewGoatUtils() *GoatUtils {
	return &GoatUtils{
		initialized: false,
	}
}

func (u *GoatUtils) SetInitialized(b bool) {
	u.initialized = b
}

func (u *GoatUtils) MustBeInitialized() {
	if !u.initialized {
		panic("goat is not initialized! Call goat.Init() before calling this function.")
	}
}

func (u *GoatUtils) IsInitialized() bool {
	return u.initialized
}
