# Visuals
There are two ways to interact with `vlx` as a user. These are either through an interactive terminal or through a non-interactive process like the window manager (but also scripting for example). When it comes to processes that are meant to provide an interaction or information to the user this distinction is crucial and decides which of the methods below gets chosen. 


## Picker
You can initiate a picking process either through an interactive terminal session or through a binding that you execute directly with your window manager. Either way, the data that is necessary for selection flows through `vlx` and then gets piped back out to the source of ingestion through the `Picker` abstraction.

If you initiated the command through your terminal, a simple `fzf` selection is shown in that same instance.

<img src="../assets/term.png" height="300"/>

If you initiated the command through other means, the data is redirected to a quickshell QML template file.

<img src="../assets/single.png" height="300"/>

There are other picker options that include a two-stage selection which is a sequential select with `fzf` in the terminal but comes with its own picker in the quickshell version.

<img src="../assets/staged.png" height="300"/>

You can also call a grouped version of the picker which is meant for ordering inside a picking dialog.

<img src="../assets/grouped.png" height="300" />

There is one exception to the picker logic which is the `PowerPicker.qml` that uses hardcoded commands to either logout, shutdown, etc.

<img src="../assets/power.png" width="800"/>


## Progress
Some actions that `vlx` invokes are considered long-running tasks (e.g., installation of a bundle). In addition to notifications that inform the user about errors, the `Progress` abstraction allows monitoring the current state of those operations. If you instantiate such a long-running operation through the window manager this popup is shown.

<img src="../assets/progress.png" width="800"/>


## Privilege
If any operation needs elevated access this will either be prompted for directly in the terminal (if your session allows it) or the user will be asked to enter their password in another quickshell popup like the one below using the `vlxpass` helper binary.

<img src="../assets/password.png" height="100"/>