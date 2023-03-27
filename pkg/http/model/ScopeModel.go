package model

import "strings"

type Scope struct {
	List []string
}

func (sc *Scope) Contain(arg string) bool {
	argLevel := strings.Split(arg, ":")
	for _, scope := range sc.List {
		scopelevel := strings.Split(scope, ":")
		if scopelevel[0] == "admin" && argLevel[0] != "admin" {
			return true
		}
		if scopelevel[0] == argLevel[0] {
			if len(scopelevel) > 1 {
				if scopelevel[1] == "*" {
					return true
				}
				if len(argLevel) > 1 {
					if scopelevel[1] == argLevel[1] {
						return true
					}
				}
			} else if len(scopelevel) == 1 && len(argLevel) == 1 {
				if scopelevel[0] == argLevel[0] {
					return true
				}
			}
		}
	}
	return false
}
