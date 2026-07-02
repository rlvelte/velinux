import QtQuick
import QtQuick.Controls
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import qs.services

FloatingWindow {
    id: launcher
    title: "vlx-launcher"
    color: "transparent"
    implicitWidth: 640
    implicitHeight: 560

    property bool shown: false
    property int selected: 0

    visible: shown

    IpcHandler {
        target: "launcher"
        function toggle(): void { launcher.shown = !launcher.shown }
        function open(): void { launcher.shown = true }
        function close(): void { launcher.shown = false }
    }

    onShownChanged: {
        if (shown) {
            searchField.text = ""
            selected = 0
            searchField.forceActiveFocus()
        }
    }

    function hide() {
        shown = false
    }

    function launch(entry) {
        if (!entry) return
        entry.execute()
        shown = false
    }

    ScriptModel {
        id: filtered
        values: {
            var all = DesktopEntries.applications.values
                .filter(function(e) { return e.name && !e.noDisplay })
                .sort(function(a, b) { return a.name.localeCompare(b.name) })

            var q = searchField.text.trim().toLowerCase()
            if (q === "") return all

            return all.filter(function(e) {
                var name = (e.name || "").toLowerCase()
                var comment = (e.comment || "").toLowerCase()
                var keywords = (e.keywords || []).join(" ").toLowerCase()
                return name.indexOf(q) !== -1 || comment.indexOf(q) !== -1 || keywords.indexOf(q) !== -1
            })
        }
    }

    Rectangle {
        anchors.fill: parent
        radius: 12
        color: Theme.base
        border.color: Theme.surface1
        border.width: 1

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
                placeholderText: "Search applications..."
                font.family: Theme.fontName
                font.pixelSize: 22
                color: Theme.text
                placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0
                    radius: 8
                    border.color: searchField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: searchField.activeFocus ? 2 : 1
                }

                Keys.onDownPressed: {
                    launcher.selected = Math.min(launcher.selected + 1, filtered.values.length - 1)
                }
                Keys.onUpPressed: {
                    launcher.selected = Math.max(launcher.selected - 1, 0)
                }
                Keys.onReturnPressed: launcher.launch(filtered.values[launcher.selected])
                Keys.onEscapePressed: launcher.hide()
                onTextChanged: {
                    if (launcher.selected >= filtered.values.length)
                        launcher.selected = Math.max(0, filtered.values.length - 1)
                }
            }

            ListView {
                id: appList
                width: parent.width
                height: parent.height - searchField.height - 12
                clip: true
                model: filtered
                currentIndex: launcher.selected
                highlightFollowsCurrentItem: true
                boundsBehavior: Flickable.StopAtBounds
                spacing: 4

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: appList.width
                    height: 64
                    radius: 8
                    color: index === launcher.selected
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                        : mouseArea.containsMouse
                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                            : "transparent"

                    Rectangle {
                        visible: index === launcher.selected
                        anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                        width: 3
                        radius: 1.5
                        color: Theme.primary
                    }

                    Row {
                        anchors.fill: parent
                        anchors.margins: 8
                        spacing: 12
                        anchors.verticalCenter: parent.verticalCenter

                        Rectangle {
                            width: 48
                            height: 48
                            radius: 10
                            color: Theme.surface0
                            anchors.verticalCenter: parent.verticalCenter

                            IconImage {
                                anchors.centerIn: parent
                                width: 32
                                height: 32
                                source: Quickshell.iconPath(modelData.icon, true)
                            }
                        }

                        Column {
                            anchors.verticalCenter: parent.verticalCenter
                            spacing: 2

                            Text {
                                text: modelData.name
                                font.family: Theme.fontName
                                font.pixelSize: 20
                                color: index === launcher.selected ? Theme.text : Theme.subtext
                                elide: Text.ElideRight
                                width: appList.width - 110
                            }

                            Text {
                                visible: !!modelData.comment
                                text: modelData.comment
                                font.family: Theme.fontName
                                font.pixelSize: 18
                                color: Theme.muted
                                elide: Text.ElideRight
                                width: appList.width - 110
                            }
                        }
                    }

                    MouseArea {
                        id: mouseArea
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onEntered: launcher.selected = index
                        onClicked: launcher.launch(modelData)
                    }
                }
            }
        }
    }
}
