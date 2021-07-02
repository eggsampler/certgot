package cli

import "strings"

type FlagList []*Flag

func (fl FlagList) Get(name string) *Flag {
	for _, f := range fl {
		if strings.EqualFold(name, f.Name) {
			return f
		}
		for _, an := range f.AltNames {
			if strings.EqualFold(an, f.Name) {
				return f
			}
		}
	}

	return nil
}
