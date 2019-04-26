# Edt - line editor

Edt is not a graphical application.

It is here to test the edit widget which is built on top of eaburns/T editor.
T stores text in the rope data structure and has an edit language on top of that.
The language is similar to sam or ed.

## Usage
```
	edt  # interactive, start with empty buffer
	edt FILE # interactive, read text from FILE
```

## TODO
```
read commands from file, operate on multiple files
	edt -e CMD FILE
	edt [-e CMD] FILES...
	
a write command
```

## Command language
See [edit](https://github.com/eaburns/T/blob/master/edit/edit.go)
