package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	subtle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dot       = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).SetString(" • ")
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1).
			Bold(true)
	
	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)
)

type IPInfo struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

type model struct {
	ipInfo    *IPInfo
	err       error
	loading   bool
	spinner   int
	quitting  bool
}

type ipMsg *IPInfo
type errMsg error

func initialModel() model {
	return model{
		loading: true,
	}
}

func getIPInfo() tea.Msg {
	url := "http://ip-api.com/json/"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return errMsg(err)
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return errMsg(err)
	}
	return ipMsg(&info)
}

// TickMsg is sent to update the spinner
type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getIPInfo, tick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	case ipMsg:
		m.ipInfo = msg
		m.loading = false
		return m, nil

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil

	case TickMsg:
		if m.loading {
			m.spinner++
			return m, tick()
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\nPress q to quit.", m.err)
	}

	if m.loading {
		spinnerChars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
		spin := spinnerChars[m.spinner%len(spinnerChars)]
		return fmt.Sprintf("\n %s Scanning network for IP details...\n\nPress q to quit.", spin)
	}

	if m.ipInfo != nil {
		info := m.ipInfo
		
		// ASCII Map approximation (very simple)
		// We could do something fancier, but let's stick to text for now.
		
		content := fmt.Sprintf(`
%s

%s %s
%s %s
%s %s, %s
%s %s
%s %s
%s %s
%s %f, %f
`,
			titleStyle.Render(" IP VISUALIZER "),
			labelStyle.Render("IP Address:"), info.Query,
			labelStyle.Render("ISP:       "), info.Isp,
			labelStyle.Render("Location:  "), info.City, info.Country,
			labelStyle.Render("Region:    "), info.RegionName,
			labelStyle.Render("Timezone:  "), info.Timezone,
			labelStyle.Render("AS:        "), info.As,
			labelStyle.Render("Coords:    "), info.Lat, info.Lon,
		)

		// Map rendering
		mapView := renderMap(info.Lat, info.Lon)
		
		infoView := infoStyle.Render(content)
		
		// Join horizontally
		ui := lipgloss.JoinHorizontal(lipgloss.Top, infoView, mapView)

		return ui + "\n\nPress q to quit."
	}

	return ""
}

// A detailed ASCII map converted to dots
var asciiWorldMap = []string{
	"           . _..::__:  ,-'-'.+       |]       ,     _,.__             ",
	"   _.___ _ _<_>`!(._`.`-.    /        _._     `_ ,_/  '  '-._.---.-.__",
	" .{     ' ' `-==,',._\\{  \\  / {) _   / _ '>_,-' `                _-/_ ",
	" \\_.:--.       `._ )`^-. ''      , [_/(                       __,/-'  ",
	"''     \\         '    _L       oD_,--'                )     /. (|    ",
	"         |           ,'         _)_.\\\\._<> 6              _,' /  '    ",
	"         `.         /          [_/_'` `'(                <'}  )       ",
	"          \\\\    .-. )          /   `-''..' `:._          _)  '        ",
	"   `        \\  (  `(          /         `:\\  > \\  ,-^.  /' '          ",
	"             `._,   ''        |           \\`'   \\|   ?_)  {\\          ",
	"                `=.---.       `._._       ,'     '`  |' ,- '.         ",
	"                  |    `-._        |     /          `:`<_|h--._       ",
	"                  (        >       .     | ,          `=.__.`-'\\      ",
	"                   `.     /        |     |{|              ,-.,\\     . ",
	"                    |   ,'          \\   / `'            ,'     \\     ",
	"                    |  /             |_'                |  __  /      ",
	"                    | |                                 |  L.\\'       ",
}

func renderMap(lat, lon float64) string {
	width := 64
	height := 17

	// Styles
	waterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("237")) // Faint dots
	landStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))  // Bright dots
	markerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true) // Green highlight

	// Create a grid of characters
	var grid [][]string
	for _, line := range asciiWorldMap {
		row := make([]string, 0)
		for _, char := range line {
			if char == ' ' {
				row = append(row, waterStyle.Render("•"))
			} else {
				row = append(row, landStyle.Render("•"))
			}
		}
		// Pad row to width with water dots
		for len(row) < width {
			row = append(row, waterStyle.Render("•"))
		}
		grid = append(grid, row)
	}

	// Map Lat/Lon to X/Y
	// Longitude: -180 to 180 -> 0 to Width
	// Latitude: -90 to 90 -> Height to 0 (Top is 90)
	
	x := int((lon + 180) / 360 * float64(width))
	y := int((90 - lat) / 180 * float64(height))

	// Clamp coordinates
	if x < 0 { x = 0 }
	if x >= width { x = width - 1 }
	if y < 0 { y = 0 }
	if y >= len(grid) { y = len(grid) - 1 }

	// Place the marker
	marker := markerStyle.Render("•")
	
	if y < len(grid) {
		row := grid[y]
		if x < len(row) {
			row[x] = marker
		}
	}

	// Render the map to a string
	s := ""
	for _, row := range grid {
		for _, char := range row {
			s += char
		}
		s += "\n"
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
