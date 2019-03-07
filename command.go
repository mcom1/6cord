package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	senderRegex = strings.NewReplacer()
)

// Commands contains multiple commands
type Commands []Command

// Command contains a command's info
type Command struct {
	Command     string
	Function    func([]string)
	Description string
}

var commands = Commands{
	Command{
		Command:     "/goto",
		Function:    gotoChannel,
		Description: "[channel name] - jumps to a channel",
	},
	Command{
		Command:     "/mentions",
		Function:    commandMentions,
		Description: "shows the last mentions",
	},
	Command{
		Command:     "/nick",
		Function:    changeSelfNick,
		Description: "[nickname] - changes nickname for the current guild",
	},
	Command{
		Command:     "/status",
		Function:    setStatus,
		Description: "[online|busy|away|invisible] - sets your status",
	},
	Command{
		Command:     "/edit",
		Function:    editMessage,
		Description: "[n:int optional] - edits the latest n message of yours",
	},
	Command{
		Command:     "/presence",
		Function:    setGame,
		Description: "[string] - sets your \"Playing\" or \"Listening to\" presence, empty to reset",
	},
	Command{
		Command:     "/react",
		Function:    reactMessage,
		Description: "[messageID:int] [emoji:string] - toggle reaction on a message",
	},
	Command{
		Command:     "/upload",
		Function:    uploadFile,
		Description: "[file path] - uploads file",
	},
	Command{
		Command:     "/copy",
		Function:    matchCopyMessage,
		Description: "[n:int] - copies the entire last n message",
	},
	Command{
		Command:     "/highlight",
		Function:    highlightMessage,
		Description: "[ID:int64] - highlights the message ID if possible",
	},
	Command{
		Command:     "/block",
		Function:    blockUser,
		Description: "[@mention] - blocks someone",
	},
	Command{
		Command:     "/unblock",
		Function:    unblockUser,
		Description: "[@mention] - unblocks someone",
	},
	Command{
		Command:     "/quit",
		Function:    commandExit,
		Description: "quits",
	},
}

func commandExit(text []string) {
	app.Stop()
	os.Exit(0)
}

// CommandHandler .
func CommandHandler() {
	text := input.GetText()
	if text == "" {
		return
	}

	defer input.SetText("")

	switch {
	case strings.HasPrefix(text, "s/"):
		go editMessageRegex(text)

	case strings.HasPrefix(text, "/"):
		f := strings.Fields(text)
		if len(f) < 0 {
			return
		}

		for _, cmd := range commands {
			if f[0] == cmd.Command && cmd.Function != nil {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							Warn(fmt.Sprintf("%v", r))
						}
					}()

					cmd.Function(f)
				}()

				return
			}
		}

		fallthrough
	default:
		// Trim literal backslash, in case "\/actual message"
		text = strings.TrimPrefix(text, `\`)
		text = senderRegex.Replace(text)

		if Channel == nil {
			Message("You're not in a channel!")
			return
		}

		go func(text string) {
			_, err := d.ChannelMessageSend(Channel.ID, text)
			if err != nil {
				Warn("Failed to send message:\n" + text + "\nError: " + err.Error())
			}
		}(text)
	}
}
