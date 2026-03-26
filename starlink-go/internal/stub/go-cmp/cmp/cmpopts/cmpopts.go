package cmpopts

import "github.com/google/go-cmp/cmp"

func EquateEmpty() cmp.Option { return nil }
