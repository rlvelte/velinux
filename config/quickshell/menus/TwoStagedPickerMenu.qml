import QtQuick
import QtQuick.Controls
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import qs.services

PanelWindow {
    id: picker
    color: "transparent"

    anchors { top: true; bottom: true; left: true; right: true }
    exclusiveZone: 0

    focusable: true
    WlrLayershell.keyboardFocus: shown ? WlrKeyboardFocus.Exclusive : WlrKeyboardFocus.None
    WlrLayershell.layer: WlrLayer.Overlay

    property bool shown: false
    property bool animatingOut: false
    property bool canClose: false

    // 0 = main item selection, 1 = subitem selection
    property int stage: 0
    property int selected: 0
    property int selectedSub: 0
    property var selectedItem: null

    property var pickerItems: []
    property string pickerResultPath: ""

    readonly property int narrowWidth: 640
    readonly property int wideWidth: 1100

    visible: shown || animatingOut

    IpcHandler {
        target: "twostagepicker"
        function vlxOpen(filePath: string, resultPath: string): void {
            picker.startPicker(filePath, resultPath)
        }
    }

    onShownChanged: {
        if (shown) {
            searchField.text = ""
            selected = 0
            selectedSub = 0
            stage = 0
            selectedItem = null
            canClose = false
            contentRect.opacity = 0
            contentRect.width = picker.narrowWidth
            dropTimer.restart()
            closeGuard.restart()
        }
    }

    function startPicker(filePath, resultPath) {
        picker.pickerResultPath = resultPath
        picker.selected = 0
        picker.selectedSub = 0
        picker.stage = 0
        picker.selectedItem = null
        itemsReader.command = ["cat", filePath]
        itemsReader.running = true
        picker.shown = true
    }

    function cancelPicker() {
        picker.writeResult({})
        picker.resetPicker()
    }

    function resetPicker() {
        picker.pickerResultPath = ""
        picker.pickerItems = []
        picker.selectedItem = null
        picker.stage = 0
        picker.selected = 0
        picker.selectedSub = 0
    }

    Timer {
        id: dropTimer
        interval: 50
        repeat: false
        onTriggered: {
            showAnim.start()
            searchField.forceActiveFocus()
        }
    }

    Timer {
        id: closeGuard
        interval: 200
        repeat: false
        onTriggered: canClose = true
    }

    Timer {
        id: hideTimer
        interval: 200
        repeat: false
        onTriggered: animatingOut = false
    }

    Process {
        id: itemsReader
        stdout: StdioCollector {
            onStreamFinished: {
                try {
                    picker.pickerItems = JSON.parse(text)
                } catch (e) {
                    console.error("Failed to parse picker items:", e)
                    picker.hide()
                }
            }
        }
    }

    Process {
        id: resultWriter
    }

    function writeResult(payload) {
        let json = JSON.stringify(payload);
        resultWriter.command = ["sh", "-c", "printf '%s' " + picker.escapeShell(json) + " > " + picker.escapeShell(picker.pickerResultPath)]
        resultWriter.running = true
    }

    function escapeShell(str) {
        return "'" + str.replace(/'/g, "'\\''") + "'"
    }

    function hide() {
        if (!shown || animatingOut) return
        animatingOut = true
        shown = false
        hideAnim.start()
    }

    // Stage 0 -> 1: select a main item, widen the picker to reveal subitems.
    function selectMain(entry) {
        if (!entry) return
        var src = entry._source
        var subs = src.subitem || []
        if (subs.length === 0) {
            // No subitems: behave like the single-stage picker.
            picker.writeResult({ item: src, subcommand: null })
            picker.resetPicker()
            hide()
            return
        }
        picker.selectedSub = 0
        picker.selectedItem = src
        picker.stage = 1
        searchField.text = ""
        widenAnim.start()
    }

    // Stage 1 -> 0: go back to the main item list.
    function backToMain() {
        picker.stage = 0
        picker.selectedSub = 0
        searchField.text = ""
        narrowAnim.start()
    }

    // Stage 1: launch the chosen subitem.
    function launchSub(entry) {
        if (!entry) return
        picker.writeResult({ item: picker.selectedItem, subcommand: entry._source })
        picker.resetPicker()
        hide()
    }

    NumberAnimation {
        id: showAnim
        target: contentRect
        property: "opacity"
        from: 0
        to: 1
        duration: 200
        easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentRect
        property: "opacity"
        from: 1
        to: 0
        duration: 200
        easing.type: Easing.InCubic
        onFinished: hideTimer.restart()
    }

    NumberAnimation {
        id: widenAnim
        target: contentRect
        property: "width"
        from: picker.narrowWidth
        to: picker.wideWidth
        duration: 250
        easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: narrowAnim
        target: contentRect
        property: "width"
        from: picker.wideWidth
        to: picker.narrowWidth
        duration: 200
        easing.type: Easing.InCubic
    }

    MouseArea {
        anchors.fill: parent
        z: -1
        enabled: picker.canClose
        onClicked: {
            picker.cancelPicker()
            picker.hide()
        }
    }

    ScriptModel {
        id: filtered
        values: {
            var all = picker.pickerItems.map(function(e) {
                return {
                    icon: e.icon,
                    name: e.header,
                    comment: e.description || "",
                    _source: e,
                    keywords: []
                }
            })

            var q = searchField.text.trim().toLowerCase()
            if (q === "" || picker.stage === 1) return all

            return all.filter(function(e) {
                var name = (e.name || "").toLowerCase()
                var comment = (e.comment || "").toLowerCase()
                return name.indexOf(q) !== -1 || comment.indexOf(q) !== -1
            })
        }
    }

    ScriptModel {
        id: subFiltered
        values: {
            if (!picker.selectedItem) return []
            var subs = (picker.selectedItem.subitem || []).map(function(s) {
                return {
                    icon: s.icon,
                    name: s.header,
                    comment: s.description || "",
                    _source: s,
                    keywords: []
                }
            })

            if (picker.stage === 0) return subs

            var q = searchField.text.trim().toLowerCase()
            if (q === "") return subs

            return subs.filter(function(e) {
                var name = (e.name || "").toLowerCase()
                var comment = (e.comment || "").toLowerCase()
                return name.indexOf(q) !== -1 || comment.indexOf(q) !== -1
            })
        }
    }

    Rectangle {
        id: contentRect
        width: picker.narrowWidth
        height: 560
        anchors.centerIn: parent
        radius: 12
        color: Theme.base
        border.color: Theme.surface1
        border.width: 1

        opacity: 0

        Behavior on width {
            NumberAnimation { duration: 250; easing.type: Easing.OutCubic }
        }

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        Column {
            anchors.fill: parent
            anchors.margins: 16
            spacing: 12

            TextField {
                id: searchField
                focus: true
                width: parent.width
                height: 52
                leftPadding: 16
                placeholderText: picker.stage === 0 ? "Search..." : "Search subitems..."
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSizeHeading
                color: Theme.text
                placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0
                    radius: 8
                    border.color: searchField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: searchField.activeFocus ? 2 : 1
                }

                Keys.onDownPressed: {
                    if (picker.stage === 0) {
                        var len = filtered.values.length
                        picker.selected = len > 0 ? (picker.selected + 1) % len : 0
                    } else {
                        var len = subFiltered.values.length
                        picker.selectedSub = len > 0 ? (picker.selectedSub + 1) % len : 0
                    }
                }
                Keys.onUpPressed: {
                    if (picker.stage === 0) {
                        var len = filtered.values.length
                        picker.selected = len > 0 ? (picker.selected - 1 + len) % len : 0
                    } else {
                        var len = subFiltered.values.length
                        picker.selectedSub = len > 0 ? (picker.selectedSub - 1 + len) % len : 0
                    }
                }
                Keys.onReturnPressed: {
                    if (picker.stage === 0) {
                        picker.selectMain(filtered.values[picker.selected])
                    } else {
                        picker.launchSub(subFiltered.values[picker.selectedSub])
                    }
                }
                Keys.onLeftPressed: {
                    if (picker.stage === 1) picker.backToMain()
                }
                Keys.onEscapePressed: {
                    if (picker.stage === 1) {
                        picker.backToMain()
                    } else {
                        picker.cancelPicker()
                        picker.hide()
                    }
                }
                onTextChanged: {
                    if (picker.stage === 0) {
                        picker.selected = 0
                    } else {
                        picker.selectedSub = 0
                    }
                }
            }

            Row {
                width: parent.width
                height: parent.height - searchField.height - 12
                spacing: 16

                ListView {
                    id: appList
                    width: 608
                    height: parent.height
                    clip: true
                    model: filtered
                    currentIndex: picker.selected
                    highlightMoveDuration: 100
                    highlightFollowsCurrentItem: true
                    boundsBehavior: Flickable.StopAtBounds
                    spacing: 4
                    opacity: picker.stage === 1 ? 0.45 : 1

                    Behavior on opacity {
                        NumberAnimation { duration: 200; easing.type: Easing.OutCubic }
                    }

                    highlight: Rectangle {
                        width: appList.width
                        height: 64
                        radius: 8
                        color: Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)

                        Rectangle {
                            anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                            width: 3
                            radius: 1.5
                            color: Theme.primary
                        }
                    }

                    delegate: Rectangle {
                        required property var modelData
                        required property int index
                        width: appList.width
                        height: 64
                        radius: 8
                        color: "transparent"

                        Row {
                            anchors.fill: parent
                            anchors.margins: 8
                            spacing: 12

                            Rectangle {
                                width: 48
                                height: 48
                                radius: 10
                                color: Theme.surface0
                                anchors.verticalCenter: parent.verticalCenter

                                Text {
                                    anchors.centerIn: parent
                                    text: modelData.icon
                                    font.pixelSize: 32
                                    color: Theme.subtext
                                }
                            }

                            Column {
                                anchors.verticalCenter: parent.verticalCenter
                                spacing: 2

                                Text {
                                    text: modelData.name
                                    font.family: Theme.fontName
                                    font.pixelSize: Theme.fontSizeLarge
                                    color: index === picker.selected ? Theme.text : Theme.subtext
                                    elide: Text.ElideRight
                                    width: appList.width - 110
                                }

                                Text {
                                    visible: !!modelData.comment
                                    text: modelData.comment
                                    font.family: Theme.fontName
                                    font.pixelSize: Theme.fontSize
                                    color: Theme.muted
                                    elide: Text.ElideRight
                                    width: appList.width - 110
                                }
                            }
                        }

                        MouseArea {
                            anchors.fill: parent
                            cursorShape: Qt.PointingHandCursor
                            enabled: picker.stage === 0
                            onClicked: picker.selectMain(modelData)
                        }
                    }
                }

                Column {
                    id: rightPanel
                    width: Math.max(0, parent.width - appList.width - 16)
                    height: parent.height
                    spacing: 8
                    opacity: picker.stage === 1 ? 1 : 0
                    visible: width > 0

                    Behavior on opacity {
                        NumberAnimation { duration: 200; easing.type: Easing.OutCubic }
                    }

                    Text {
                        width: parent.width
                        text: picker.selectedItem ? picker.selectedItem.header : ""
                        font.family: Theme.fontName
                        font.pixelSize: Theme.fontSizeLarge
                        color: Theme.text
                        elide: Text.ElideRight
                    }

                    ListView {
                        id: subList
                        width: parent.width
                        height: parent.height - 32
                        clip: true
                        model: subFiltered
                        currentIndex: picker.selectedSub
                        highlightMoveDuration: 100
                        highlightFollowsCurrentItem: true
                        boundsBehavior: Flickable.StopAtBounds
                        spacing: 4

                        highlight: Rectangle {
                            width: subList.width
                            height: 64
                            radius: 8
                            color: Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)

                            Rectangle {
                                anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                                width: 3
                                radius: 1.5
                                color: Theme.primary
                            }
                        }

                        delegate: Rectangle {
                            required property var modelData
                            required property int index
                            width: subList.width
                            height: 64
                            radius: 8
                            color: "transparent"

                            Row {
                                anchors.fill: parent
                                anchors.margins: 8
                                spacing: 12

                                Rectangle {
                                    width: 48
                                    height: 48
                                    radius: 10
                                    color: Theme.surface0
                                    anchors.verticalCenter: parent.verticalCenter

                                    Text {
                                        anchors.centerIn: parent
                                        text: modelData.icon
                                        font.pixelSize: 32
                                        color: Theme.subtext
                                    }
                                }

                                Column {
                                    anchors.verticalCenter: parent.verticalCenter
                                    spacing: 2

                                    Text {
                                        text: modelData.name
                                        font.family: Theme.fontName
                                        font.pixelSize: Theme.fontSizeLarge
                                        color: index === picker.selectedSub ? Theme.text : Theme.subtext
                                        elide: Text.ElideRight
                                        width: subList.width - 110
                                    }

                                    Text {
                                        visible: !!modelData.comment
                                        text: modelData.comment
                                        font.family: Theme.fontName
                                        font.pixelSize: Theme.fontSize
                                        color: Theme.muted
                                        elide: Text.ElideRight
                                        width: subList.width - 110
                                    }
                                }
                            }

                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                enabled: picker.stage === 1
                                onClicked: picker.launchSub(modelData)
                            }
                        }
                    }
                }
            }
        }
    }
}
