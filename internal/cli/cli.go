package cli

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/gasuhwbab/cli_todo_app/internal/db"
	"github.com/gasuhwbab/cli_todo_app/internal/note"
)

func Run() error {
	fmt.Println("App started. Use /add, /delete, /update, /check, /help commands.")
	var command string
	fmt.Scan(&command)
	switch command {
	case "/help":
		help()
	case "/add":
		if err := add(); err != nil {
			return err
		}
	case "/delete":
		if err := delete(); err != nil {
			return err
		}
	case "/update":
		if err := update(); err != nil {
			return err
		}
	case "/check":
		if err := check(); err != nil {
			return err
		}
	default:

	}
	return nil
}

func help() {
	fmt.Println("Use /add command to add note in your TODO list")
	fmt.Println("Use /delete command to add note from your TODO list")
	fmt.Println("Use /update command to update note in your TODO list")
	fmt.Println("Use /check command to view all notes from your TODO list")
}

func add() error {
	fmt.Println("Write name (daily, monthly, yearly) and text of note you want to add")
	var text, name string
	if _, err := fmt.Scan(&name, &text); err != nil {
		return err
	}
	noteToAdd := note.NewNote(text)
	buf, err := noteToAdd.MarshalBinary()
	if err != nil {
		return err
	}
	if err := db.Db.Add([]byte(name), buf); err != nil {
		return err
	}
	return nil
}

func delete() error {
	fmt.Println("Write name (daily, monthly, yearly) and id of note you want to delete")
	var id uint64
	var name string
	if _, err := fmt.Scan(&name, &id); err != nil {
		return err
	}
	binaryId := itob(id)
	if err := db.Db.Delete([]byte(name), binaryId); err != nil {
		return err
	}
	return nil
}

func update() error {
	fmt.Println("Write name (daily, monthly, yearly), id of note you want to update and new text")
	var id uint64
	var name string
	var text string
	if _, err := fmt.Scan(&name, &id, &text); err != nil {
		return err
	}
	binaryId := itob(id)
	to, err := note.NewNote(text).MarshalBinary()
	if err != nil {
		return err
	}
	if err := db.Db.Update([]byte(name), binaryId, to); err != nil {
		return err
	}
	return nil
}

func check() error {
	fmt.Println("Write name of note (daily, monthly, yearly)")
	var name string
	fmt.Scan(&name)
	if name != "daily" && name != "monthly" && name != "yearly" {
		return errors.New("unavaliable name of note")
	}
	bufs, err := db.Db.Get([]byte(name))
	if err != nil {
		return err
	}
	notes := make([]note.Note, len(bufs))
	for i, buf := range bufs {
		if err := notes[i].UnmarshalBinary(buf); err != nil {
			return err
		}
		fmt.Println(notes[i].Id, notes[i].Time.Format("2006/01/02"), notes[i].Text)
	}
	return nil
}

func itob(x uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, x)
	return buf
}
