import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import QtQuick
import QtQuick.Layouts

// Centered launcher overlay. Toggle from outside with:
//   qs ipc call launcher toggle
// e.g. bound in hyprland.conf:
//   bind = $mod, D, exec, qs ipc call launcher toggle
PanelWindow {
    id: launcher

    property bool shown: false
    property int selected: 0

    anchors { top: true; left: true; right: true; bottom: true }
    color: "transparent"
    visible: shown
    focusable: shown
    WlrLayershell.layer: WlrLayer.Overlay

    IpcHandler {
        target: "launcher"
        function toggle(): void { launcher.shown = !launcher.shown; }
        function open(): void { launcher.shown = true; }
        function close(): void { launcher.shown = false; }
    }

    onShownChanged: {
        if (shown) {
            searchField.text = "";
            selected = 0;
            searchField.forceActiveFocus();
        }
    }

    // -------- filtered app list --------
    ScriptModel {
        id: filtered
        values: {
            const all = [...DesktopEntries.applications.values]
                .filter(e => e.name && !e.noDisplay)
                .sort((a, b) => a.name.localeCompare(b.name));

            const q = searchField.text.trim().toLowerCase();
            if (q === "") return all;

            return all.filter(e => {
                const name = (e.name || "").toLowerCase();
                const comment = (e.comment || "").toLowerCase();
                const keywords = (e.keywords || []).join(" ").toLowerCase();
                return name.includes(q) || comment.includes(q) || keywords.includes(q);
            });
        }
    }

    function launch(entry) {
        if (!entry) return;
        entry.execute();
        launcher.shown = false;
    }

    function clampSelection() {
        if (selected >= filtered.values.length) selected = Math.max(0, filtered.values.length - 1);
        if (selected < 0) selected = 0;
    }

    // -------- scrim --------
    Rectangle {
        anchors.fill: parent
        color: Qt.rgba(Theme.base.r, Theme.base.g, Theme.base.b, 0.55)

        MouseArea {
            anchors.fill: parent
            onClicked: launcher.shown = false
        }
    }

    // -------- panel --------
    Rectangle {
        id: panel
        anchors.centerIn: parent
        width: 420
        height: Math.min(420, headerRow.height + resultsList.contentHeight + footerRow.height + 24)
        radius: 12
        color: Theme.mantle
        border.width: 1
        border.color: Theme.surface0

        MouseArea {
            anchors.fill: parent
            // absorb clicks so the scrim MouseArea doesn't close the launcher
        }

        ColumnLayout {
            anchors.fill: parent
            spacing: 0

            RowLayout {
                id: headerRow
                Layout.fillWidth: true
                Layout.margins: 12
                spacing: 8

                Text {
                    text: "\u2315"
                    color: Theme.primary
                    font.family: Theme.fontNameMono
                    font.pixelSize: Theme.fontSizeLarge
                }

                TextInput {
                    id: searchField
                    Layout.fillWidth: true
                    color: Theme.text
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSize
                    clip: true

                    Keys.onDownPressed: {
                        launcher.selected = Math.min(launcher.selected + 1, filtered.values.length - 1);
                    }
                    Keys.onUpPressed: {
                        launcher.selected = Math.max(launcher.selected - 1, 0);
                    }
                    Keys.onReturnPressed: launcher.launch(filtered.values[launcher.selected])
                    Keys.onEnterPressed: launcher.launch(filtered.values[launcher.selected])
                    Keys.onEscapePressed: launcher.shown = false

                    onTextChanged: launcher.clampSelection()
                }

                Text {
                    text: "esc to close"
                    color: Theme.muted
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeSmall
                }
            }

            Rectangle { Layout.fillWidth: true; height: 1; color: Theme.surface0 }

            ListView {
                id: resultsList
                Layout.fillWidth: true
                Layout.preferredHeight: Math.min(contentHeight, 300)
                Layout.margins: 6
                clip: true
                model: filtered
                currentIndex: launcher.selected
                highlightFollowsCurrentItem: true
                boundsBehavior: Flickable.StopAtBounds

                delegate: Rectangle {
                    id: row
                    required property var modelData
                    required property int index

                    width: resultsList.width
                    height: 40
                    radius: 8
                    color: index === launcher.selected
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                        : "transparent"
                    border.width: index === launcher.selected ? 0 : 0

                    Rectangle {
                        visible: index === launcher.selected
                        anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                        width: 2
                        color: Theme.primary
                    }

                    RowLayout {
                        anchors.fill: parent
                        anchors.leftMargin: 10
                        anchors.rightMargin: 10
                        spacing: 10

                        Rectangle {
                            width: 26
                            height: 26
                            radius: 6
                            color: Theme.surface0

                            IconImage {
                                anchors.centerIn: parent
                                width: 16
                                height: 16
                                source: Quickshell.iconPath(row.modelData.icon, true)
                            }
                        }

                        ColumnLayout {
                            Layout.fillWidth: true
                            spacing: 0

                            Text {
                                Layout.fillWidth: true
                                text: row.modelData.name
                                color: index === launcher.selected ? Theme.text : Theme.subtext
                                font.family: Theme.fontName
                                font.pixelSize: Theme.fontSize
                                elide: Text.ElideRight
                            }
                            Text {
                                Layout.fillWidth: true
                                visible: !!row.modelData.comment
                                text: row.modelData.comment
                                color: Theme.muted
                                font.family: Theme.fontName
                                font.pixelSize: Theme.fontSizeSmall
                                elide: Text.ElideRight
                            }
                        }
                    }

                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onEntered: launcher.selected = row.index
                        hoverEnabled: true
                        onClicked: launcher.launch(row.modelData)
                    }
                }
            }

            Rectangle { Layout.fillWidth: true; height: 1; color: Theme.surface0 }

            RowLayout {
                id: footerRow
                Layout.fillWidth: true
                Layout.margins: 10

                Text {
                    text: filtered.values.length + " results"
                    color: Theme.muted
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeSmall
                }
                Item { Layout.fillWidth: true }
                Text {
                    text: "\u2191\u2193 select \u00B7 \u21B5 open"
                    color: Theme.muted
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeSmall
                }
            }
        }
    }
}
