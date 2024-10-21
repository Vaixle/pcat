package internal

import (
	"fmt"
	"strings"
)

type Flag interface {
	apply(*ModifiableText)
}

type ModifiableText struct {
	text string
	skip bool
}

type FlagNumberNonBlank struct {
	lineIndex int
}

func (f *FlagNumberNonBlank) apply(mText *ModifiableText) {
	if mText.text != "" {
		mText.text = fmt.Sprintf("%d %s", f.lineIndex, mText.text)
		f.lineIndex++
	}
}

type FlagNumberAll struct {
	lineIndex int
}

func (f *FlagNumberAll) apply(mText *ModifiableText) {
	mText.text = fmt.Sprintf("%d %s", f.lineIndex, mText.text)
	f.lineIndex++
}

type FlagShowNonprinting struct {
}

func (s FlagShowNonprinting) apply(mText *ModifiableText) {
	var result strings.Builder
	for _, r := range []rune(mText.text) {
		switch {
		case r < 32 && r != '\n' && r != '\t':
			result.WriteString("^")
			result.WriteRune(r + '@')
		case r == 127:
			result.WriteString("^?")
		default:
			result.WriteRune(r)
		}
	}
	mText.text = result.String()
}

type FlagShowEnds struct {
}

func (s FlagShowEnds) apply(mText *ModifiableText) {
	mText.text += "$"
}

type FlagShowTabs struct {
}

func (s FlagShowTabs) apply(mText *ModifiableText) {
	mText.text = strings.ReplaceAll(mText.text, "\t", "^I")
}

type FlagSqueezeBlank struct {
	isSecondLineEmpty bool
}

func (f *FlagSqueezeBlank) apply(mText *ModifiableText) {
	if f.isSecondLineEmpty {
		if mText.text == "" {
			mText.skip = true
		} else {
			f.isSecondLineEmpty = false
			mText.skip = false
		}
	}

	if mText.text == "" {
		f.isSecondLineEmpty = true
	}
}
