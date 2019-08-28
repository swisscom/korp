package autocompletion

import (
	"fmt"

	"github.com/urfave/cli"
)

const (
	bashScript = `#! /bin/bash
_cli_bash_autocomplete() {
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
complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete go-cli`

	zshScript = `_cli_zsh_autocomplete() {
	local -a opts
	opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")
	_describe 'values' opts
	return
}
compdef _cli_zsh_autocomplete go-cli`
)

func bashAutocomplete(c *cli.Context) error {

	fmt.Println(bashScript)
	return nil
}

func zshAutocomplete(c *cli.Context) error {

	fmt.Println(zshScript)
	return nil
}
