package connectwise

type PatchOp struct {
	Op    Op          `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type Op string

const (
	AddOp     Op = "add"
	ReplaceOp Op = "replace"
	RemoveOp  Op = "Remove"
)
