/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Output shell completion code for the specified shell (bash or zsh)",
	Example: `  # Installing bash completion on macOS using homebrew
  ## If running Bash 3.2 included with macOS
  brew install bash-completion
  ## or, if running Bash 4.1+
  brew install bash-completion@2


  # Installing bash completion on Linux
  ## If bash-completion is not installed on Linux, please install the 'bash-completion' package
  ## via your distribution's package manager.
  ## Load the ormb completion code for bash into the current shell
  source <(ormb completion bash)
  ## Write bash completion code to a file and source if from .bash_profile
  ormb completion bash > ~/.completion.bash.inc
  printf "
  # ormb shell completion
  source '$HOME/.completion.bash.inc'
  " >> $HOME/.bash_profile
  source $HOME/.bash_profile


  # Load the ormb completion code for zsh[1] into the current shell
  source <(ormb completion zsh)
  # Set the ormb completion code for zsh[1] to autoload on startup
  ormb completion zsh > "${fpath[1]}/_ormb"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			switch args[0] {
			case "bash":
				return cmd.Parent().GenBashCompletion(os.Stdout)
			case "zsh":
				return genCompletionZsh(os.Stdout, cmd.Parent())
			default:
				return errors.New("Unsupported shell type " + args[0])
			}
		}
		return cmd.Help()
	},
}

// Adapted from https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/completion/completion.go#L145.
func genCompletionZsh(out io.Writer, ormb *cobra.Command) error {
	zshHead := "#compdef ormb\n"
	out.Write([]byte(zshHead))

	zshInitialization := `
__ormb_bash_source() {
	alias shopt=':'
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__ormb_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__ormb_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__ormb_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__ormb_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__ormb_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__ormb_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__ormb_filedir() {
	# Don't need to do anything here.
	# Otherwise we will get trailing space without "compopt -o nospace"
	true
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q 'GNU\|BusyBox'; then
	LWORD='\<'
	RWORD='\>'
fi

__ormb_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__ormb_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__ormb_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__ormb_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__ormb_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__ormb_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__ormb_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	ormb.GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}

__ormb_bash_source <(__ormb_convert_bash_to_zsh)
`
	out.Write([]byte(zshTail))
	return nil
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
