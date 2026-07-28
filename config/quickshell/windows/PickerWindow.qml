import QtQuick
import QtQuick.Controls
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import qs.globals
import "../util/utils.js" as Utils

PanelWindow {
    id: window

    color: "transparent"
    anchors { top: true; bottom: true; left: true; right: true }
    exclusiveZone: 0

    focusable: true
    WlrLayershell.keyboardFocus: shown ? WlrKeyboardFocus.Exclusive : WlrKeyboardFocus.None
    WlrLayershell.layer: WlrLayer.Overlay

    Keys.onEscapePressed: window.escapeKeyHandler()

    property string ipcTarget: ""

    property bool shown: false
    property bool animatingOut: false
    property bool canClose: false
    property int selected: 0
    property var pickerItems: []
    property string pickerResultPath: ""
    property string searchQuery: ""

    property var selectedIndices: new Set()

    property var cancelResult: null
    property Component contentBody: null
    property Component itemDelegate: null
    property Component extraHeader: null
    property int contentWidth: 640
    property int contentHeight: 560

    property var returnKeyHandler: function(entry) { launch(entry) }
    property var escapeKeyHandler: function() { cancelPicker(); hide() }
    property var downKeyHandler: function() {
        window.selected = Math.min(window.selected + 1, filtered.values.length - 1)
    }
    property var upKeyHandler: function() {
        window.selected = Math.max(window.selected - 1, 0)
    }
    property var textChangeHandler: function() {
        if (window.selected >= filtered.values.length)
            window.selected = Math.max(0, filtered.values.length - 1)
    }

    property var buildModel: function(items, query) {
        var all = items.map(function(e, i) {
            return {
                icon: e.icon, name: e.header, comment: e.description || "",
                _source: e, keywords: [], _index: i
            }
        })
        if (!query) return all
        var q = query.toLowerCase()
        return all.filter(function(e) {
            return (e.name || "").toLowerCase().indexOf(q) !== -1
                || (e.comment || "").toLowerCase().indexOf(q) !== -1
        })
    }

    visible: shown || animatingOut

    IpcHandler {
        target: window.ipcTarget
        function vlxOpen(filePath: string, resultPath: string): void {
            window.startPicker(filePath, resultPath)
        }
    }

    function startPicker(filePath, resultPath) {
        window.pickerResultPath = resultPath
        window.selected = 0
        window.selectedIndices = new Set()
        window.searchQuery = ""
        itemsReader.command = ["cat", filePath]
        itemsReader.running = true
        window.shown = true
    }

    function cancelPicker() { writeResult(cancelResult); resetPicker() }

    function resetPicker() {
        window.pickerResultPath = ""
        window.pickerItems = []
        window.selectedIndices = new Set()
    }

    onShownChanged: {
        if (shown) {
            window.selected = 0
            window.selectedIndices = new Set()
            window.canClose = false
            window.searchQuery = ""
            contentRect.opacity = 0
            dropTimer.restart()
            closeGuard.restart()
        }
    }

    Timer { id: dropTimer; interval: Anims.dropDelay; repeat: false
        onTriggered: { showAnim.start(); contentItem.focusSearch() } }

    Timer { id: closeGuard; interval: Anims.duration; repeat: false
        onTriggered: window.canClose = true }

    Timer { id: hideTimer; interval: Anims.duration; repeat: false
        onTriggered: window.animatingOut = false }

    Process {
        id: itemsReader
        stdout: StdioCollector {
            onStreamFinished: {
                try { window.pickerItems = JSON.parse(text) }
                catch (e) { console.error("Failed to parse picker items:", e); window.hide() }
            }
        }
    }

    Process { id: resultWriter }

    function writeResult(payload) {
        var json = JSON.stringify(payload)
        resultWriter.command = ["sh", "-c", "printf '%s' "
            + Utils.escapeShell(json) + " > " + Utils.escapeShell(window.pickerResultPath)]
        resultWriter.running = true
    }

    function hide() {
        if (!shown || animatingOut) return
        animatingOut = true; shown = false; hideAnim.start()
    }

    function launch(entry) {
        if (!entry) return
        writeResult(entry._source); resetPicker(); hide()
    }

    NumberAnimation {
        id: showAnim
        target: contentRect; property: "opacity"
        from: 0; to: 1; duration: Anims.duration; easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentRect; property: "opacity"
        from: 1; to: 0; duration: Anims.duration; easing.type: Easing.InCubic
        onFinished: hideTimer.restart()
    }

    MouseArea {
        anchors.fill: parent; z: -1; enabled: window.canClose
        onClicked: { window.cancelPicker(); window.hide() }
    }

    Rectangle {
        id: contentRect
        width: window.contentWidth; height: window.contentHeight
        anchors.centerIn: parent
        radius: Dimensions.overlayRadius; color: Theme.base
        border.color: Theme.surface1; border.width: 1
        opacity: 0

        MouseArea { anchors.fill: parent; propagateComposedEvents: false; onClicked: {} }

        Loader {
            id: contentLoader
            anchors.fill: parent
            anchors.margins: 16
            sourceComponent: window.contentBody || defaultContent
        }
    }

    QtObject {
        id: contentItem
        function focusSearch() {
            if (!contentLoader.item || !contentLoader.item.focusSearch) return
            contentLoader.item.focusSearch()
        }
    }

    Component {
        id: defaultContent

        Column {
            spacing: 12

            function focusSearch() { searchField.forceActiveFocus() }

            TextField {
                id: searchField
                focus: true; width: parent.width; height: Dimensions.searchFieldHeight; leftPadding: 16
                placeholderText: "Search..."
                font.family: Theme.fontName; font.pixelSize: Theme.fontSizeHeading
                color: Theme.text; placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0; radius: Dimensions.inputRadius
                    border.color: searchField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: searchField.activeFocus ? 2 : 1
                }

                onTextChanged: { window.searchQuery = searchField.text.trim().toLowerCase(); window.textChangeHandler() }
                Keys.onReturnPressed: window.returnKeyHandler(filtered.values[window.selected])
                Keys.onEscapePressed: window.escapeKeyHandler()
                Keys.onDownPressed: window.downKeyHandler()
                Keys.onUpPressed: window.upKeyHandler()
            }

            Loader {
                width: parent.width
                sourceComponent: window.extraHeader
                visible: window.extraHeader !== null
            }

            ListView {
                id: appList
                width: parent.width
                height: parent.parent.height - searchField.height - 12
                    - (window.extraHeader ? 48 : 0)
                clip: true; cacheBuffer: 400; model: filtered
                currentIndex: window.selected
                highlightMoveDuration: Anims.durationHighlight
                highlight: Rectangle { color: "transparent" }
                boundsBehavior: Flickable.StopAtBounds
                spacing: 4

                delegate: window.itemDelegate
            }
        }
    }

    property alias filteredModel: filtered

    ScriptModel {
        id: filtered
        values: window.buildModel(window.pickerItems, window.searchQuery)
    }
}
