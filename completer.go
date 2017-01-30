package ishell

import "strings"

type iCompleter struct {
	cmd *Cmd
}

func (ic iCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	words := strings.Fields(string(line))
	var cWords []string
	prefix := ""
	if len(words) > 0 && line[pos-1] != ' ' {
		prefix = words[len(words)-1]
		cWords = ic.getWords(words[:len(words)-1])
	} else {
		cWords = ic.getWords(words)
	}

	var suggestions [][]rune
	for _, w := range cWords {
		if strings.HasPrefix(w, prefix) {
			suggestions = append(suggestions, []rune(strings.TrimPrefix(w, prefix)))
		}
	}
	return suggestions, len(prefix)
}

func (ic iCompleter) getWords(w []string) (s []string) {
	cmd, args := ic.cmd.FindCmd(w)
	if cmd == nil {
		cmd, args = ic.cmd, w
	}
	if cmd.Completer != nil {
		return cmd.Completer(args)
	}
	for k := range cmd.children {
		s = append(s, k)
	}
	return
}
