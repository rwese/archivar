package caller

import (
	"runtime"
	"strings"
)

func FactoryPackage() string {
	return PackageName(3, true)
}

func PackageName(skip int, onlyLast bool) string {
	pc, _, _, _ := runtime.Caller(skip)
	pcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(pcName, ".")
	pl := len(parts)
	packageName := ""
	if parts[pl-2][0] == '(' {
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	if onlyLast {
		parts = strings.Split(parts[1], "/")
		packageName = parts[len(parts)-1]
	}

	return packageName
}
