package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/pflag"
	"pik/model"
	"pik/motd"
	"pik/spool"
	"pik/viewport"
)

type Model struct {
	*model.HydratedState
	Index         int
	Indices       map[int]model.HydratedTarget
	SourceIndices map[int]*model.HydratedSource
	Cancel        bool
	Done          bool
	Height        int
	Alt           bool
	AutoAlt       bool
	Motd          string
}

func (m *Model) Init() tea.Cmd {
	_, h, err := term.GetSize(0)
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
	}
	m.Height = h
	wantsAlt := viewport.NeedsViewport(m.State(), m.Height)
	if m.AutoAlt && wantsAlt {
		return tea.EnterAltScreen
	}
	return nil
}

func (m *Model) HandleResize(msg tea.WindowSizeMsg) tea.Cmd {
	if !m.AutoAlt {
		return nil
	}
	m.Height = msg.Height
	if viewport.NeedsViewport(m.State(), msg.Height) {
		m.Alt = true
		return tea.EnterAltScreen
	} else {
		m.Alt = false
		return tea.ExitAltScreen
	}

}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	var result tea.Cmd
	switch mt := msg.(type) {
	case tea.WindowSizeMsg:
		result = m.HandleResize(mt)
	case tea.KeyMsg:
		result, err = m.HandleInput(mt)
	case tea.Cmd:
		result, err = m.HandleSignal(mt)
	}
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
	}
	return m, result
}

func (m *Model) HandleSignal(cmd tea.Cmd) (tea.Cmd, error) {
	return nil, nil
}

func (m *Model) View() string {
	if m.Cancel || m.Done {
		return ""
	}
	result := m.State()
	result = viewport.Process(result, m.Height)
	return result
}

func (m *Model) Result() (*model.HydratedSource, model.HydratedTarget) {
	if m.Cancel {
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

func NewModel(st *model.State, hydrators []model.Modder) *Model {
	m := &Model{
		HydratedState: Hydrate(st, hydrators),
		Index:         0,
		Indices:       make(map[int]model.HydratedTarget),
		SourceIndices: make(map[int]*model.HydratedSource),
		AutoAlt:       !pflag.Lookup("inline").Changed,
		Motd:          motd.One(),
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
