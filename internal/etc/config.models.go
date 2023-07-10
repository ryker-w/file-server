package etc

type Configuration struct {
	Web        web        `toml:"web"`
	FileSystem fileSystem `toml:"fileSystem"`
}

type web struct {
	Listen string `toml:"listen"`
	Cache  int    `toml:"cache"`
}

type fileSystem struct {
	Root string `toml:"root"`
}
