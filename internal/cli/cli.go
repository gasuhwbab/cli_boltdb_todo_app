package cli

import (
	"encoding/binary"
	"fmt"

	"github.com/gasuhwbab/cli_todo_app/internal/db"
	"github.com/gasuhwbab/cli_todo_app/internal/note"
)

var id uint16 = 0

type CLI struct {
	storage *db.Storage
}

func NewCLI(storage *db.Storage) *CLI {
	return &CLI{storage}
}

func (cli *CLI) Run() error {
	fmt.Println("App started. Use /add, /delete, /update, /check, /help commands.")
	var command string
	fmt.Scan(&command)
	switch command {
	case "/help":
		cli.help()
	case "/add":
		if err := cli.add(); err != nil {
			return err
		}
	case "/delete":
		if err := cli.delete(); err != nil {
			return err
		}
	case "/update":
		if err := cli.update(); err != nil {
			return err
		}
	case "/check":
		if err := cli.check(); err != nil {
			return err
		}
	default:

	}
	return nil
}

func (cli *CLI) help() {
	fmt.Println("Use /add command to add note in your TODO list")
	fmt.Println("Use /delete command to add note from your TODO list")
	fmt.Println("Use /update command to update note in your TODO list")
	fmt.Println("Use /check command to view all notes from your TODO list")
}

func (cli *CLI) add() error {
	fmt.Println("Write note you want to add")
	var text string
	if _, err := fmt.Scan(&text); err != nil {
		return err
	}
	noteToAdd := note.NewNote(id, text)
	id++
	buf, err := noteToAdd.MarshalBinary()
	if err != nil {
		return err
	}
	if err := cli.storage.Add(buf); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) delete() error {
	fmt.Println("Write id and date (format: 2006/01/02) of note you want to delete")
	var id uint16
	var date string
	if _, err := fmt.Scan(&id, &date); err != nil {
		return err
	}
	buf := make([]byte, 2+8)
	binary.BigEndian.PutUint16(buf[0:2], id)
	buf = append(buf, []byte(date)...)
	return nil
}

func (cli *CLI) update() error {
	fmt.Println("Write id, date (format: 2006/01/02) of note you want to update and new text")
	var id uint16
	var date string
	var text string
	if _, err := fmt.Scan(&id, &date, &text); err != nil {
		return err
	}
	to := note.NewNote(id, text)
	binaryTo, err := to.MarshalBinary()
	if err != nil {
		return err
	}
	from := make([]byte, 2+8)
	binary.BigEndian.PutUint16(from[0:2], id)
	from = append(from, []byte(text)...)
	if err := cli.storage.Update(from, binaryTo); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) check() error {
	bufs, err := cli.storage.Get()
	if err != nil {
		return err
	}
	notes := make([]note.Note, len(bufs))
	for i, buf := range bufs {
		if err := notes[i].UnmarshalBinary(buf); err != nil {
			return err
		}
		fmt.Println(notes[i].Id, notes[i].Time, notes[i].Text)
	}
	return nil
}
