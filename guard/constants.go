package guard

type Action string

const (
	ActionRead   Action = "R"
	ActionCreate Action = "C"
	ActionUpdate Action = "U"
	ActionDelete Action = "D"
)

func (a Action) String() string {
	return string(a)
}

func (a Action) IsValid() bool {
	return a == ActionRead || a == ActionCreate || a == ActionUpdate || a == ActionDelete
}
