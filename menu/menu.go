package menu

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"pik/model"
	"pik/spool"
)

var WrongModelTypeError = errors.New("wrong model type")
var NoSourcesIndexedError = errors.New("no sources indexed")

func Show(st *model.State, hydrators []model.Hydrator) (*model.HydratedSource, model.HydratedTarget, error) {
	if len(st.Sources) == 0 {
		return nil, nil, NoSourcesIndexedError
	}
	md := NewModel(st, hydrators)
	program := tea.NewProgram(md)
	resultModel, err := program.Run()
	if err != nil {
		return nil, nil, err
	}
	result, ok := resultModel.(*Model)
	if !ok {
		return nil, nil, WrongModelTypeError
	}

	src, t := result.Result()
	return src, t, nil
}

func Hydrate(st *model.State, hydrators []model.Hydrator) *model.HydratedState {
	hyd := &model.HydratedState{
		State:           st,
		HydratedSources: make([]*model.HydratedSource, len(st.Sources)),
	}
	for i, s := range st.Sources {
		hydSrc := s.Hydrate(hydrators)

		for _, h := range hydrators {
			err := h.Hydrate(s, hydSrc)
			if err != nil {
				spool.Warn("%v\n", err)
				continue
			}
		}

		hyd.HydratedSources[i] = hydSrc
	}
	return hyd
}
