package mise

var toolOrder = []string{
	"dotnet", "go", "zig", "rust", "python", "ruby", "lua", "node", "bun",
	"java", "kotlin", "scala", "dart", "flutter", "R", "php", "julia",
	"swift", "perl", "haskell", "groovy", "clojure", "erlang", "elixir",
	"fsharp", "nim", "ocaml", "crystal", "v", "coffeescript", "deno",
	"cmake", "gradle", "maven", "ninja", "ansible", "helm", "terraform",
}

var toolIcons = map[string]string{
	"dotnet":       "\uE77F",
	"go":           "\U000F07D3",
	"zig":          "\uE8EF",
	"rust":         "\uE7A8",
	"python":       "\uE73C",
	"ruby":         "\U000F0D2D",
	"lua":          "\uE826",
	"node":         "\uED0D",
	"bun":          "\uE76F",
	"java":         "\uE738",
	"kotlin":       "\uE81B",
	"scala":        "\uE737",
	"dart":         "\uE64C",
	"flutter":      "\uE7DD",
	"R":            "\uE881",
	"php":          "\uE73D",
	"julia":        "\uE80D",
	"swift":        "\uE699",
	"perl":         "\uE67E",
	"haskell":      "\uE777",
	"groovy":       "\uE775",
	"clojure":      "\uE76A",
	"erlang":       "\uE7B1",
	"elixir":       "\uE7CD",
	"fsharp":       "\uE7A7",
	"nim":          "\uE841",
	"ocaml":        "\uE67A",
	"crystal":      "\uE7AC",
	"v":            "\uE6AC",
	"coffeescript": "\uE751",
	"deno":         "\uE7C0",
	"cmake":        "\uE794",
	"gradle":       "\uE7F2",
	"maven":        "\uE82C",
	"ninja":        "\U000F0774",
	"ansible":      "\uE723",
	"helm":         "\uE7FB",
	"terraform":    "\uE8BD",
}

type RegistryTool struct {
	Short       string   `json:"short"`
	Backends    []string `json:"backends"`
	Description string   `json:"description"`
	Aliases     []string `json:"aliases,omitempty"`
}

type ToolVersion struct {
	Tool    string `json:"tool"`
	Version string `json:"version"`
	Active  bool   `json:"active"`
	Icon    string `json:"icon"`
}

type RemoteVersion struct {
	Version   string `json:"version"`
	CreatedAt string `json:"created_at"`
}

type miseToolVersions map[string][]struct {
	Version          string `json:"version"`
	RequestedVersion string `json:"requested_version"`
	Installed        bool   `json:"installed"`
	Active           bool   `json:"active"`
}
