pragma Singleton
import QtQuick
import Quickshell
import Quickshell.Io

Singleton {
    id: vlx

    property var themes: []
    property string activeThemeId: ""

    function listThemes() {
        listProcess.running = true
    }

    function applyTheme(id) {
        applyProcess.command = ["vlx", "themes", "apply", id]
        applyProcess.running = true
    }

    function switchTheme() {
        Quickshell.execDetached(["vlx", "themes", "switch"])
    }

    Process {
        id: listProcess
        command: ["vlx", "themes", "list", "--json"]
        stdout: StdioCollector {
            onStreamFinished: {
                try {
                    var data = JSON.parse(text)
                    vlx.themes = data
                    for (var i = 0; i < data.length; i++) {
                        if (data[i].active) {
                            vlx.activeThemeId = data[i].id
                            break
                        }
                    }
                } catch (e) {
                    console.error("Failed to parse themes JSON:", e)
                }
            }
        }
    }

    Process {
        id: applyProcess
        stdout: StdioCollector { onStreamFinished: {} }
    }

    Component.onCompleted: listThemes()
}
