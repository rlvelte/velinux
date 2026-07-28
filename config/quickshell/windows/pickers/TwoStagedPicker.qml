import QtQuick
import QtQuick.Controls
import Quickshell
import qs.core

PickerWindow {
    id: picker

    ipcTarget: "twostagepicker"
    cancelResult: {}
    contentWidth: picker.narrowWidth

    property int stage: 0
    property int selectedSub: 0
    property var selectedItem: null
    readonly property int narrowWidth: 640
    readonly property int wideWidth: 1100

    buildModel: function(items, query) {
        var all = items.map(function(e, i) {
            return {
                icon: e.icon, name: e.header, comment: e.description || "",
                installed: e.installed === true, _source: e, keywords: [], _index: i
            }
        })
        if (!query || picker.stage === 1) return all
        var q = query.toLowerCase()
        return all.filter(function(e) {
            return (e.name || "").toLowerCase().indexOf(q) !== -1
                || (e.comment || "").toLowerCase().indexOf(q) !== -1
        })
    }

    returnKeyHandler: function(entry) {
        if (picker.stage === 0) picker.selectMain(entry)
        else picker.launchSub(subFiltered.values[picker.selectedSub])
    }
    escapeKeyHandler: function() {
        if (picker.stage === 1) picker.backToMain()
        else { cancelPicker(); hide() }
    }
    downKeyHandler: function() {
        if (picker.stage === 0) {
            var len = picker.filteredModel.values.length
            picker.selected = len > 0 ? (picker.selected + 1) % len : 0
        } else {
            var subLen = subFiltered.values.length
            picker.selectedSub = subLen > 0 ? (picker.selectedSub + 1) % subLen : 0
        }
    }
    upKeyHandler: function() {
        if (picker.stage === 0) {
            var len = picker.filteredModel.values.length
            picker.selected = len > 0 ? (picker.selected - 1 + len) % len : 0
        } else {
            var subLen = subFiltered.values.length
            picker.selectedSub = subLen > 0 ? (picker.selectedSub - 1 + subLen) % subLen : 0
        }
    }
    textChangeHandler: function() {
        if (picker.stage === 0) picker.selected = 0
        else picker.selectedSub = 0
    }

    function selectMain(entry) {
        if (!entry) return
        var src = entry._source
        var subs = src.subitem || []
        if (subs.length === 0) {
            writeResult({ item: src, subcommand: null })
            resetPicker(); hide()
            return
        }
        picker.selectedItem = src
        picker.stage = 1
        picker.searchQuery = ""

        var firstSel = 0
        for (var i = 0; i < subs.length; i++) {
            if (subs[i].installed !== true) { firstSel = i; break }
        }
        picker.selectedSub = firstSel
        picker.contentWidth = picker.wideWidth
    }

    function backToMain() {
        picker.stage = 0
        picker.selectedSub = 0
        picker.searchQuery = ""
        picker.contentWidth = picker.narrowWidth
    }

    function launchSub(entry) {
        if (!entry || entry.installed) return
        writeResult({ item: picker.selectedItem, subcommand: entry._source })
        resetPicker(); hide()
    }

    ScriptModel {
        id: subFiltered
        values: {
            if (!picker.selectedItem) return []
            var subs = (picker.selectedItem.subitem || []).map(function(s) {
                return {
                    icon: s.icon, name: s.header, comment: s.description || "",
                    installed: s.installed === true, _source: s, keywords: []
                }
            })
            if (picker.stage === 0) return subs
            var q = picker.searchQuery
            if (!q) return subs
            return subs.filter(function(e) {
                return (e.name || "").toLowerCase().indexOf(q) !== -1
                    || (e.comment || "").toLowerCase().indexOf(q) !== -1
            })
        }
    }

    contentBody: twoStageContent

    Component {
        id: twoStageContent

        Column {
            spacing: 12

            function focusSearch() { searchField.forceActiveFocus() }

            TextField {
                id: searchField
                focus: true; width: parent.width; height: Theme.searchFieldHeight; leftPadding: 16
                placeholderText: picker.stage === 0 ? "Search..." : "Search subitems..."
                font.family: Theme.fontName; font.pixelSize: Theme.fontSizeHeading
                color: Theme.text; placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0; radius: Theme.inputRadius
                    border.color: searchField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: searchField.activeFocus ? 2 : 1
                }

                onTextChanged: { picker.searchQuery = searchField.text.trim().toLowerCase(); picker.textChangeHandler() }
                Keys.onReturnPressed: picker.returnKeyHandler(picker.filteredModel.values[picker.selected])
                Keys.onEscapePressed: picker.escapeKeyHandler()
                Keys.onDownPressed: picker.downKeyHandler()
                Keys.onUpPressed: picker.upKeyHandler()
                Keys.onLeftPressed: { if (picker.stage === 1) picker.backToMain() }
            }

            Row {
                width: parent.width
                height: parent.height - Theme.searchFieldHeight - 12
                spacing: 16

                ListView {
                    id: appList
                    width: picker.narrowWidth - 32
                    height: parent.height
                    clip: true; model: picker.filteredModel
                    currentIndex: picker.selected
                    highlightMoveDuration: Theme.animDurationHighlight
                    highlightFollowsCurrentItem: true
                    boundsBehavior: Flickable.StopAtBounds
                    spacing: 4
                    opacity: picker.stage === 1 ? 0.45 : 1

                    Behavior on opacity {
                        NumberAnimation { duration: Theme.animDuration; easing.type: Easing.OutCubic }
                    }

                    highlight: Rectangle {
                        width: appList.width; height: Theme.listItemHeight; radius: Theme.inputRadius
                        color: Theme.primarySelected
                        Rectangle {
                            anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                            width: 3; radius: 1.5; color: Theme.primary
                        }
                    }

                    delegate: Rectangle {
                        required property var modelData
                        required property int index
                        width: appList.width; height: Theme.listItemHeight; radius: Theme.inputRadius
                        color: "transparent"

                        ItemRow {
                            anchors.fill: parent
                            anchors.margins: 8
                            icon: modelData.icon
                            name: modelData.name
                            comment: modelData.comment || ""
                            selected: index === picker.selected
                            textWidth: appList.width - (modelData.installed ? 150 : 110)
                            glyphSize: 32
                        }

                        Rectangle {
                            visible: modelData.installed
                            anchors { right: parent.right; rightMargin: 12; verticalCenter: parent.verticalCenter }
                            width: 24; height: 24; radius: 12
                            color: Theme.surface1
                        }

                        MouseArea {
                            anchors.fill: parent; cursorShape: Qt.PointingHandCursor
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
                        NumberAnimation { duration: Theme.animDuration; easing.type: Easing.OutCubic }
                    }

                    Text {
                        width: parent.width
                        text: picker.selectedItem ? picker.selectedItem.header : ""
                        font.family: Theme.fontName; font.pixelSize: Theme.fontSizeLarge
                        color: Theme.text; elide: Text.ElideRight
                    }

                    ListView {
                        id: subList
                        width: parent.width
                        height: parent.height - 32
                        clip: true; model: subFiltered
                        currentIndex: picker.selectedSub
                        highlightMoveDuration: Theme.animDurationHighlight
                        highlightFollowsCurrentItem: true
                        boundsBehavior: Flickable.StopAtBounds
                        spacing: 4

                        highlight: Rectangle {
                            width: subList.width; height: Theme.listItemHeight; radius: Theme.inputRadius
                            color: Theme.primarySelected
                            Rectangle {
                                anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                                width: 3; radius: 1.5; color: Theme.primary
                            }
                        }

                        delegate: Rectangle {
                            required property var modelData
                            required property int index
                            width: subList.width; height: Theme.listItemHeight; radius: Theme.inputRadius
                            color: "transparent"
                            opacity: modelData.installed ? 0.45 : 1

                            ItemRow {
                                anchors.fill: parent
                                anchors.margins: 8
                                icon: modelData.icon
                                name: modelData.name
                                comment: modelData.comment || ""
                                selected: index === picker.selectedSub
                                textWidth: subList.width - (modelData.installed ? 200 : 110)
                            }

                            Rectangle {
                                visible: modelData.installed
                                anchors { right: parent.right; rightMargin: 12; verticalCenter: parent.verticalCenter }
                                width: installedRow.width + 16; height: 24; radius: 12
                                color: Theme.surface1

                                Row {
                                    id: installedRow; anchors.centerIn: parent; spacing: 5
                                    Text {
                                        text: "✓"
                                        font.family: Theme.fontName; font.pixelSize: Theme.fontSizeSmall
                                        color: Theme.primary; anchors.verticalCenter: parent.verticalCenter
                                    }
                                    Text {
                                        text: "installed"
                                        font.family: Theme.fontName; font.pixelSize: Theme.fontSizeSmall
                                        color: Theme.subtext; anchors.verticalCenter: parent.verticalCenter
                                    }
                                }
                            }

                            MouseArea {
                                anchors.fill: parent
                                cursorShape: modelData.installed ? Qt.ArrowCursor : Qt.PointingHandCursor
                                enabled: picker.stage === 1 && !modelData.installed
                                onClicked: picker.launchSub(modelData)
                            }
                        }
                    }
                }
            }
        }
    }
}
