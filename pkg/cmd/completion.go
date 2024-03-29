// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

var completionShells = map[string]func(out io.Writer, boilerPlate string, cmd *cobra.Command) error{
	"bash": runCompletionBash,
	"zsh":  runCompletionZsh,
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

// RunCompletion checks given arguments and executes command.
func RunCompletion(out io.Writer, shell string, cmd *cobra.Command) error {
	runFunc, found := completionShells[shell]
	if !found {
		return UsageErrorf(cmd, "Unsupported shell type %q.", shell)
	}

	return runFunc(out, "", cmd)
}

func runCompletionBash(out io.Writer, boilerPlate string, cmd *cobra.Command) error {
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return cmd.GenBashCompletion(out)
}

const (
	zshHead           = "#compdef kt\n"
	zshInitialization = `
__kt_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__kt_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__kt_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__kt_compgen() {
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

__kt_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__kt_ltrim_colon_completions()
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

__kt_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__kt_filedir() {
	local RET OLD_IFS w qw

	__kt_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi

	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"

	IFS="," __kt_debug "RET=${RET[@]} len=${#RET[@]}"

	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__kt_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__kt_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
	printf %q "$1"
    fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__kt_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__kt_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__kt_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__kt_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__kt_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__kt_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__kt_type/g" \
	<<'BASH_COMPLETION_EOF'
`

	zshTail = `
BASH_COMPLETION_EOF
}

__kt_bash_source <(__kt_convert_bash_to_zsh)
_complete kt 2>/dev/null
`
)

func runCompletionZsh(out io.Writer, boilerPlate string, cmd *cobra.Command) (errs error) {
	buf := new(bytes.Buffer)

	_, err := buf.WriteString(zshHead)
	errs = multierr.Append(errs, err)

	_, err = buf.WriteString(zshHead)
	errs = multierr.Append(errs, err)

	if boilerPlate != "" {
		_, err = buf.WriteString(boilerPlate)
		errs = multierr.Append(errs, err)
	}

	_, err = buf.WriteString(zshInitialization)
	errs = multierr.Append(errs, err)

	err = cmd.GenBashCompletion(buf)
	errs = multierr.Append(errs, err)

	_, err = buf.WriteString(zshTail)
	errs = multierr.Append(errs, err)

	_, err = out.Write(buf.Bytes())
	errs = multierr.Append(errs, err)

	return multierr.Combine(errs)
}
