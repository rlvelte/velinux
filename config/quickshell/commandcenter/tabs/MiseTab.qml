import QtQuick
import Quickshell
import Quickshell.Io
import qs.services
Item {
    id: root
    property string filterText: ""
    property var toolsData: []
    property var groupedTools: []
    property int cols: 4
    property int tileSpacing: 12
    Timer {
        interval: 100
        running: true
        repeat: false
        onTriggered: lsProcess.running = true
    }
    function refresh() {
        lsProcess.running = true
    }
    function moveUp() {
        var nc = Math.max(Math.floor(grid.width / grid.cellWidth), 1)
        var newIdx = grid.currentIndex - nc
        if (newIdx >= 0) grid.currentIndex = newIdx
    }
    function moveDown() {
        var nc = Math.max(Math.floor(grid.width / grid.cellWidth), 1)
        var newIdx = grid.currentIndex + nc
        if (newIdx < grid.count) grid.currentIndex = newIdx
    }
    function activate() {
        if (grid.currentIndex < 0 || grid.currentIndex >= groupedTools.length) return
        var tool = groupedTools[grid.currentIndex]
        if (tool.activeVersion) return
        var ver = tool.versions[0]
        if (ver) setVersion(tool.tool, ver.version)
    }
    function groupTools(data) {
        var map = {}
        for (var i = 0; i < data.length; i++) {
            var item = data[i]
            if (!map[item.tool]) {
                map[item.tool] = { tool: item.tool, icon: item.icon, versions: [], activeVersion: "" }
            }
            map[item.tool].versions.push({ version: item.version, requested: item.requested_version, active: item.active })
            if (item.active) map[item.tool].activeVersion = item.version
        }
        var result = []
        for (var key in map) {
            map[key].versions.sort(function(a, b) { return b.version.localeCompare(a.version) })
            result.push(map[key])
        }
        result.sort(function(a, b) { return a.tool.localeCompare(b.tool) })
        return result
    }
    function setVersion(tool, version) {
        useProcess.toolName = tool
        useProcess.toolVersion = version
        useProcess.running = true
    }
    Process {
        id: lsProcess
        command: ["vlx", "cc", "--json", "mise", "ls"]
        stdout: StdioCollector {
            onStreamFinished: {
                if (!text) return
                try {
                    root.toolsData = JSON.parse(text)
                    root.groupedTools = root.groupTools(root.toolsData)
                } catch(e) {
                    console.error("mise parse err:", (text || "").substring(0, 200))
                }
            }
        }
        stderr: StdioCollector {
            onStreamFinished: {
                if (text) console.error("mise stderr:", text)
            }
        }
    }
    Process {
        id: useProcess
        property string toolName: ""
        property string toolVersion: ""
        command: useProcess.toolName ? ["vlx", "cc", "--json", "mise", "use", useProcess.toolName, useProcess.toolVersion] : ["true"]
        stdout: StdioCollector {
            onStreamFinished: {
                if (useProcess.exitCode === 0) lsProcess.running = true
            }
        }
    }
    Process {
        id: installProcess
        property string toolName: ""
        property string toolVersion: ""
        command: installProcess.toolName ? ["vlx", "cc", "--json", "mise", "install", installProcess.toolName, installProcess.toolVersion] : ["true"]
        stdout: StdioCollector {
            onStreamFinished: {
                if (installProcess.exitCode === 0) lsProcess.running = true
            }
        }
    }
    Column {
        anchors.fill: parent
        spacing: 6
        Rectangle {
            width: parent.width
            height: 32
            radius: 8
            color: installInput.activeFocus ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.08) : "transparent"
            border.color: installInput.activeFocus ? Theme.primary : Theme.surface1
            border.width: 1
            TextInput {
                id: installInput
                anchors { left: parent.left; leftMargin: 10; right: parent.right; rightMargin: 10; verticalCenter: parent.verticalCenter }
                font.family: Theme.fontNameMono
                font.pixelSize: Theme.fontSize
                color: installInput.activeFocus ? Theme.text : Theme.muted
                Keys.onEscapePressed: {
                    installInput.text = ""
                    parent.parent.forceActiveFocus()
                }
                Keys.onReturnPressed: {
                    if (installInput.text.trim()) {
                        var parts = installInput.text.trim().split("@")
                        installProcess.toolName = parts[0]
                        installProcess.toolVersion = parts.length > 1 ? parts[1] : "latest"
                        installProcess.running = true
                        installInput.text = ""
                    }
                }
            }
            Text {
                visible: !installInput.text && !installInput.activeFocus
                anchors { left: parent.left; leftMargin: 10; verticalCenter: parent.verticalCenter }
                text: "\uf899  Install (tool@version)..."
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSizeSmall
                color: Theme.muted
            }
            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.IBeamCursor
                onClicked: installInput.forceActiveFocus()
            }
        }
        Rectangle {
            width: parent.width
            height: 100
            color: Theme.surface0
            radius: 8
            Text {
                anchors.centerIn: parent
                text: "Tools: " + root.groupedTools.length + (root.groupedTools.length > 0 ? "  " + root.groupedTools[0].tool : "")
                font.family: Theme.fontNameMono
                font.pixelSize: Theme.fontSize
                color: Theme.warning
            }
        }
        GridView {
            id: grid
            width: parent.width
            height: parent.height - 32 - 6 - 100 - 12
            cellWidth: Math.max(Math.floor(width / cols), 110)
            cellHeight: 130
            interactive: true
            clip: true
            boundsBehavior: Flickable.StopAtBounds
            highlightFollowsCurrentItem: true
            model: root.groupedTools
            delegate: Item {
                required property var modelData
                required property int index
                width: grid.cellWidth
                height: grid.cellHeight
                Rectangle {
                    id: card
                    anchors.fill: parent
                    anchors.margins: root.tileSpacing / 2
                    radius: 14
                    color: modelData.activeVersion
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.1)
                        : index === grid.currentIndex
                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                            : Theme.surface0
                    border.color: modelData.activeVersion ? Theme.primary : "transparent"
                    border.width: modelData.activeVersion ? 2 : 0
                    Column {
                        anchors.centerIn: parent
                        spacing: 8
                        Rectangle {
                            id: logoCircle
                            anchors.horizontalCenter: parent.horizontalCenter
                            width: 52
                            height: 52
                            radius: width / 2
                            color: modelData.activeVersion ? Theme.primary : Theme.surface1
                            Text {
                                anchors.centerIn: parent
                                text: modelData.icon
                                font.pixelSize: 22
                                color: modelData.activeVersion ? Theme.background : Theme.subtext
                            }
                        }
                        Text {
                            anchors.horizontalCenter: parent.horizontalCenter
                            text: modelData.tool
                            font.family: Theme.fontName
                            font.pixelSize: Theme.fontSize
                            font.weight: Font.Bold
                            color: Theme.text
                        }
                        Text {
                            anchors.horizontalCenter: parent.horizontalCenter
                            text: modelData.activeVersion || "--"
                            font.family: Theme.fontNameMono
                            font.pixelSize: Theme.fontSizeSmall
                            color: modelData.activeVersion ? Theme.primaryMuted : Theme.muted
                        }
                    }
                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            grid.currentIndex = index
                            if (!modelData.activeVersion) {
                                var ver = modelData.versions[0]
                                if (ver) root.setVersion(modelData.tool, ver.version)
                            }
                        }
                    }
                }
            }
        }
    }
}