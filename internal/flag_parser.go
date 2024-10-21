package internal

import (
	"flag"
	"fmt"
)

type FlagParser interface {
	ParseFlags() ([]Flag, bool)
}

type PcatFlagParser struct {
}

func NewPcatFlagParser() FlagParser {
	return &PcatFlagParser{}
}

func (f *PcatFlagParser) ParseFlags() ([]Flag, bool) {
	var flags []Flag

	showAll := flag.Bool("A", false, "equivalent to -vET")
	numberNonBlank := flag.Bool("b", false, "number nonempty output lines, overrides -n")
	showEndsAndNonPrinting := flag.Bool("e", false, "equivalent to -vE")
	showEnds := flag.Bool("E", false, "display $ at end of each line")
	numberAll := flag.Bool("n", false, "number all output lines")
	squeezeBlank := flag.Bool("s", false, "suppress repeated empty output lines")
	showTabs := flag.Bool("T", false, "display TAB characters as ^I")
	showNonPrinting := flag.Bool("v", false, "use ^ and M- notation, except for LFD and TAB")
	help := flag.Bool("help", false, "display this help and exit")
	version := flag.Bool("version", false, "output version information and exit")

	flag.Parse()

	if *help {
		fmt.Println("Usage: cat [OPTION]... [FILE]...")
		fmt.Println("Concatenate FILE(s) to standard output.")
		flag.PrintDefaults()
		return nil, true
	}

	if *version {
		fmt.Println("pcat version 1.0")
		return nil, true
	}

	if *numberNonBlank {
		*numberAll = false
	}

	if *showEndsAndNonPrinting {
		*showNonPrinting = true
		*showEnds = true
	}

	if *showEndsAndNonPrinting {
		*showNonPrinting = true
		*showTabs = true
	}

	if *showAll {
		*showNonPrinting = true
		*showEnds = true
		*showTabs = true
	}

	if *numberNonBlank {
		flags = append(flags, &FlagNumberNonBlank{
			lineIndex: 1,
		})
	}

	if *numberAll {
		flags = append(flags, &FlagNumberAll{
			lineIndex: 1,
		})
	}

	if *showEnds {
		flags = append(flags, &FlagShowEnds{})
	}

	if *squeezeBlank {
		flags = append(flags, &FlagSqueezeBlank{})
	}

	if *showTabs {
		flags = append(flags, &FlagShowTabs{})
	}

	if *showNonPrinting {
		flags = append(flags, &FlagShowNonprinting{})
	}

	return flags, false
}
