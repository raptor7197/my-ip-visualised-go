# IP Visualizer TUI

A beautiful, terminal-based IP address visualizer written in Go. This tool fetches your public IP information and displays it alongside a dot-matrix world map with your location highlighted.
<img width="1099" height="457" alt="image" src="https://github.com/user-attachments/assets/05e45361-8cb3-431d-93a5-fbc0a11162ba" />


## Features

- **IP Information**: Displays your public IP, ISP, Location, Region, Timezone, and AS number.
- **World Map Visualization**: Renders a stylish ASCII/dot-matrix world map.
- **Location Highlighting**: Pinpoints your approximate location on the map with a green indicator.
- **Terminal UI**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for a premium terminal experience.

## Installation

Ensure you have Go installed on your system.

1. Clone the repository:
   ```bash
   git clone https://github.com/raptor7197/my-ip-visualised-go.git
   cd my-ip-visualised-go
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Run the program directly:

```bash
go run main.go
```

Or build and run the binary:

```bash
go build -o ipvis
./ipvis
```

## Controls

- **q** or **Ctrl+C**: Quit the application.

## Author

**raptor7197**
- GitHub: [raptor7197](https://github.com/raptor7197)

## License

MIT
