package goat

type utils struct {
	initialized bool
}

type utilsInterface interface {
	SetInitialized(bool)
	MustBeInitialized()
	IsInitialized() bool
}

func newUtils() *utils {
	return &utils{
		initialized: false,
	}
}

func (u *utils) SetInitialized(b bool) {
	u.initialized = b
}

func (u *utils) MustBeInitialized() {
	if !u.initialized {
		panic("goat is not initialized! Call goat.Init() before calling this function.")
	}
}

func (u *utils) IsInitialized() bool {
	return u.initialized
}
