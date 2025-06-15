package validator

import (
	"thop/dom/model/project"
	"thop/dom/problem"
)

const (
	ErrNilValue = problem.Key("VALIDATOR_NIL_VALUE")
)

func ValidateProject(p *project.Project) error {
	if p == nil {
		return ErrNilValue.WithMessage("tried to pass nil value to validator")
	}

	err := p.Validate()
	if err != nil {
		return err
	}

	err = p.Template.Validate()
	if err != nil {
		return err
	}

	for _, w := range p.Template.Windows {
		err = w.Validate()
		if err != nil {
			return err
		}

		for _, p := range w.Panes {
			err = p.Validate()
			if err != nil {
				return err
			}
		}
	}

	return err
}
