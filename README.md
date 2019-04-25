# ui - go gui toolkit
Version 1

Ui should be a re-implementation of duit:

- based on shiny
- single screen application
- mouse and key events forwarded from shiny
- draw directly on screen.Buffer with the standard draw.Image interface
- scale size with shift+wheel

## Restrictions of duit

- mjl-/duit initially used plan9 devdraw as a backend, which needs plan9ports as an external dependency
- ktye/duitdraw is a backend that replaces plan9ports with shiny
- ktye/duit uses duitdraw as a backend and allows a single-binary application

The initial dependency of plan9 draw still shows.
Ui tries to polish this but keeps most of duit's design.
