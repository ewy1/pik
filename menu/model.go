package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"pik/model"
	"pik/spool"
)

type Model struct {
	*model.HydratedState
	Index         int
	Indices       map[int]model.HydratedTarget
	SourceIndices map[int]*model.HydratedSource
	Quit          bool
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	var result tea.Cmd
	switch mt := msg.(type) {
	case tea.KeyMsg:
		result, err = m.HandleInput(mt)
	case tea.Cmd:
		result, err = m.HandleSignal(mt)
	}
	if err != nil {
		spool.Warn("%v\n", err)
	}
	return m, result
}

func (m *Model) HandleInput(msg tea.KeyMsg) (tea.Cmd, error) {
	var cmd tea.Cmd
	switch msg.String() {
	case "up", "k":
		m.Index--
	case "down", "j":
		m.Index++
	case "q", "esc":
		m.Quit = true
		cmd = tea.Quit
	case "space", " ", "enter":
		cmd = tea.Quit
	}

	m.Validate()

	return cmd, nil
}

func (m *Model) HandleSignal(cmd tea.Cmd) (tea.Cmd, error) {
	return nil, nil
}

func (m *Model) View() string {
	return m.State(m.HydratedState)
}

func (m *Model) Result() (*model.HydratedSource, model.HydratedTarget) {
	if m.Quit {
		return nil, nil
	}
	return m.SourceIndices[m.Index], m.Indices[m.Index]
}

func (m *Model) Validate() {
	if m.Index < 0 {
		m.Index = 0
	}
	if m.Index > len(m.Indices)-1 {
		m.Index = len(m.Indices) - 1
	}
}

func NewModel(st *model.State, hydrators []model.Hydrator) *Model {
	m := &Model{
		HydratedState: Hydrate(st, hydrators),
		Index:         0,
		Indices:       make(map[int]model.HydratedTarget),
		SourceIndices: make(map[int]*model.HydratedSource),
	}
	idx := 0
	for _, src := range st.Sources {
		hydSrc := src.Hydrate(hydrators)
		for _, target := range src.Targets {

			if !target.Visible() {
				continue
			}

			hydTarget, err := target.Hydrate(src)
			m.Indices[idx] = hydTarget
			if err != nil {
				spool.Warn("%v\n", err)
			}
			m.SourceIndices[idx] = hydSrc

			idx++
		}
	}
	return m
}
