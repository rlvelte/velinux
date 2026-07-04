import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import qs.services

PanelWindow {
    id: panel
    color: "transparent"

    anchors { top: true; bottom: true; left: true; right: true }
    exclusiveZone: 0

    focusable: true
    WlrLayershell.keyboardFocus: shown ? WlrKeyboardFocus.Exclusive : WlrKeyboardFocus.None
    WlrLayershell.layer: WlrLayer.Overlay

    property bool shown: false
    property bool animatingOut: false
    property bool canClose: false
    property int currentTab: 0

    property int screenWidth: screen?.geometry?.width ?? (Quickshell.screens.length > 0 ? Quickshell.screens[0].width : 1920)
    property int screenHeight: screen?.geometry?.height ?? (Quickshell.screens.length > 0 ? Quickshell.screens[0].height : 1080)

    property var rightTabs: [feedsTab, bundesligaTab]

    visible: shown || animatingOut

    IpcHandler {
        target: "commandCenter"
        function toggle(): void { if (panel.shown) panel.hide(); else panel.shown = true }
        function open(): void { panel.shown = true }
        function close(): void { panel.hide() }
    }

    onShownChanged: {
        if (shown) {
            currentTab = 0
            canClose = false
            closeGuard.restart()
            contentRect.opacity = 0
            showAnim.restart()
            contentRect.forceActiveFocus()
            feedsTab.refresh()
            bundesligaTab.refresh()
        }
    }

    Timer {
        id: closeGuard
        interval: 200
        repeat: false
        onTriggered: canClose = true
    }

    Timer {
        interval: 5000
        running: panel.shown
        repeat: true
        triggeredOnStart: false
        onTriggered: {
            feedsTab.refresh()
            bundesligaTab.refresh()
        }
    }

    function hide() {
        if (!shown || animatingOut) return
        animatingOut = true
        hideAnim.restart()
    }

    NumberAnimation {
        id: showAnim
        target: contentRect
        property: "opacity"
        from: 0
        to: 1
        duration: 150
        easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentRect
        property: "opacity"
        from: 1
        to: 0
        duration: 150
        easing.type: Easing.InCubic
        onFinished: {
            shown = false
            animatingOut = false
        }
    }

    MouseArea {
        anchors.fill: parent
        z: -1
        enabled: panel.canClose
        onClicked: panel.hide()
    }

    Rectangle {
        id: contentRect
        width: Math.round(panel.screenWidth * 0.3)
        height: Math.round(panel.screenHeight * 0.7)
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.topMargin: 8
        anchors.leftMargin: 8
        color: Theme.base
        radius: 12
        border.color: Theme.surface1
        border.width: 1

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        focus: true
        Keys.onEscapePressed: panel.hide()
        Keys.onTabPressed: {
            panel.currentTab = (panel.currentTab + 1) % 2
        }
        Keys.onUpPressed: rightTabs[panel.currentTab].moveUp()
        Keys.onDownPressed: rightTabs[panel.currentTab].moveDown()
        Keys.onReturnPressed: rightTabs[panel.currentTab].activate()

        Column {
            anchors.fill: parent
            anchors.margins: 12
            spacing: 8

            Header {
            }

            Item {
                width: parent.width
                height: parent.height - 32 - 24

                Column {
                    width: parent.width
                    height: parent.height
                    spacing: 8

                        Row {
                            width: parent.width
                            height: 32
                            spacing: 4

                            Repeater {
                                model: [
                                    { label: "Feeds" },
                                    { label: "Bundesliga" }
                                ]

                                delegate: Rectangle {
                                    required property var modelData
                                    required property int index
                                    width: (parent.width - parent.spacing) / 2
                                    height: parent.height
                                    radius: 6
                                    color: index === panel.currentTab
                                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                                        : tabMouse.containsMouse
                                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                                            : "transparent"

                                    Rectangle {
                                        visible: index === panel.currentTab
                                        anchors { left: parent.left; right: parent.right; bottom: parent.bottom }
                                        height: 2
                                        radius: 1
                                        color: Theme.primary
                                    }

                                    Text {
                                        anchors.centerIn: parent
                                        text: modelData.label
                                        font.family: Theme.fontName
                                        font.pixelSize: Theme.fontSize
                                        font.weight: index === panel.currentTab ? Font.Bold : Font.Normal
                                        color: index === panel.currentTab ? Theme.text : Theme.subtext
                                    }

                                    MouseArea {
                                        id: tabMouse
                                        anchors.fill: parent
                                        hoverEnabled: true
                                        cursorShape: Qt.PointingHandCursor
                                        onClicked: panel.currentTab = index
                                    }
                                }
                            }
                        }

                        Item {
                            width: parent.width
                            height: parent.height - 32 - 8

                            FeedsTab {
                                id: feedsTab
                                visible: panel.currentTab === 0
                                anchors.fill: parent
                                filterText: ""
                            }

                            BundesligaTab {
                                id: bundesligaTab
                                visible: panel.currentTab === 1
                                anchors.fill: parent
                                filterText: ""
                            }
                    }
                }
            }
        }
    }
}

