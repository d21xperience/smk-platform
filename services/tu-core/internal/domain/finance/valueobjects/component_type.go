package valueobjects

import "errors"

type ComponentType string

const (
	ComponentSPP     ComponentType = "SPP"
	ComponentUniform ComponentType = "UNIFORM"
	ComponentBook    ComponentType = "BOOK"
	ComponentExam    ComponentType = "EXAM"
)

func NewComponentType(t string) (ComponentType, error) {
	switch ComponentType(t) {
	case ComponentSPP, ComponentUniform, ComponentBook, ComponentExam:
		return ComponentType(t), nil
	default:
		return "", errors.New("invalid component type: must be SPP, UNIFORM, BOOK, or EXAM")
	}
}

func (t ComponentType) String() string {
	return string(t)
}

func (t ComponentType) IsPeriodic() bool {
	return t == ComponentSPP
}
