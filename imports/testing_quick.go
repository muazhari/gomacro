// this file was generated by gomacro command: import "testing/quick"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	pkg "testing/quick"
	. "reflect"
)

func init() {
	Binds["testing/quick"] = map[string]Value{
		"Check":	ValueOf(pkg.Check),
		"CheckEqual":	ValueOf(pkg.CheckEqual),
		"Value":	ValueOf(pkg.Value),
	}
	Types["testing/quick"] = map[string]Type{
		"CheckEqualError":	TypeOf((*pkg.CheckEqualError)(nil)).Elem(),
		"CheckError":	TypeOf((*pkg.CheckError)(nil)).Elem(),
		"Config":	TypeOf((*pkg.Config)(nil)).Elem(),
		"Generator":	TypeOf((*pkg.Generator)(nil)).Elem(),
		"SetupError":	TypeOf((*pkg.SetupError)(nil)).Elem(),
	}
}