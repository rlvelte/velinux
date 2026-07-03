# How does the picker stuff work?
There are the folllowing actors in a picking process:
- Sway (the wm)
- Quickshell (the main ui)
- Terminal (the secondary ui)
- vlx (the management app)

You can initiate a picking process either trough an interactive terminal session or throug a binding that you execute directly with your window manager. In either way the data that is necessary for selection flows throug `vlx` and then gets piped back out to the source of ingestion through the `Picker` abstraction.
If you initiated the command through your terminal, a simple `fzf` selection is shown in that same instance.
If you initiated the command through other means (most likely a wayland wm), the data is redireected to a quickshell QML template file.

## What can i pick?
There is no technical limitation on what the picker can/cannot serve. However its design is based on a few requirements to make it blend in well.
- Icon (SVG)
- Header
- Small (!) description

Also the picker comes with a search bar by default.
