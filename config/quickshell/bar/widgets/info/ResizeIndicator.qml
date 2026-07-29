import QtQuick
import Quickshell
import Quickshell.I3
import qs.globals

Item {
    id: root
    implicitWidth: resizeText.implicitWidth
    implicitHeight: Dimensions.barHeight

    property bool inResizeMode: false

    visible: inResizeMode

    I3IpcListener {
        subscriptions: ["mode"]
        onIpcEvent: function (event) {
            if (event.type === "mode") {
                var parsed = JSON.parse(event.data);
                root.inResizeMode = (parsed.change === "resize");
            }
        }
    }

    Text {
        id: resizeText
        anchors.verticalCenter: parent.verticalCenter
        text: "in resize mode"
        color: Theme.error
        font.family: Theme.fontName
        font.pixelSize: Theme.fontSizeHeading
        font.weight: Font.Bold
    }
}
