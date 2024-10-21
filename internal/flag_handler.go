package internal

type FlagHandler interface {
	Execute(*ModifiableText)
	SetNext(FlagHandler)
}

type PcatFlagHandler struct {
	flags []Flag
	next  FlagHandler
}

func NewPcatFlagHandler(flags ...Flag) FlagHandler {
	return &PcatFlagHandler{
		flags: flags,
	}
}

func (f *PcatFlagHandler) Execute(text *ModifiableText) {
	for _, flag := range f.flags {
		flag.apply(text)
	}
}

func (f *PcatFlagHandler) SetNext(next FlagHandler) {
	f.next = next
}
