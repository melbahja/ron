<div align="center">
	<h1>Ron</h1>
    <h4 align="center">
	  The simple command line task runner.
	</h4>
</div>

<p align="center">
    <a href="#installation">Installation</a> ❘
    <a href="#usage">Usage</a> ❘
    <a href="#bash-autocomplete">Bash Autocomplete</a> ❘
    <a href="#license">License</a>
</p>


## Installation

Visit [releases](https://github.com/melbahja/ron/releases) page and download latest version, or you can `go get github.com/melbahja/ron`.

## Usage

Ron is a very simple task runner that execute any executable file inside a `.ron` directory for example if you have this tree in your project: 

```
.ron/
├── foo/
│   ├── .default
│   ├── bar
│   └── baz
├── .default
└── serve
```

`.default` file is the default executable for directories.

To execute `.ron/.default` run:
```bash
ron
```

To execute `.ron/serve` file run:
```bash
ron serve
```

To execute `.ron/foo/.default` run:
```bash
ron foo
```

To execute `.ron/foo/bar` run:
```bash
ron foo bar
```

Note: also you can use binary files inside `.ron` directory.

See Ron's [.ron](https://github.com/melbahja/ron/tree/master/.ron) directory. 

## Bash Autocomplete

Add this to your bash profile (`.bashrc`):

```bash
_ron_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -F _ron_bash_autocomplete ron
```

## License

Ron is provided under the [MIT License](https://github.com/melbahja/ron/blob/master/LICENSE).
