import QtQuick
import Quickshell
import Quickshell.Services.Pipewire
import qs.services

Item {
    id: audioWidget
    implicitWidth: audioRow.implicitWidth
    implicitHeight: 40

    PwObjectTracker {
        objects: [Pipewire.defaultAudioSink]
    }

    property var sink: Pipewire.defaultAudioSink
    property bool ready: sink && sink.ready && sink.audio

    Row {
        id: audioRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 6

        Text {
            text: {
                if (!audioWidget.ready) return "--"
                var vol = Math.round(audioWidget.sink.audio.volume * 100)
                return vol + "%"
            }
            color: {
                if (!audioWidget.ready) return Theme.muted
                if (audioWidget.sink.audio.muted) return Theme.muted
                return Theme.subtext
            }
            font.family: Theme.fontName
            font.pixelSize: 18
            font.weight: Font.Medium
        }
    }

    MouseArea {
        anchors.fill: parent
        cursorShape: Qt.PointingHandCursor
        onClicked: {
            if (audioWidget.ready) {
                audioWidget.sink.audio.muted = !audioWidget.sink.audio.muted
            }
        }
        onWheel: function(wheel) {
            if (!audioWidget.ready) return
            var delta = wheel.angleDelta.y > 0 ? 0.02 : -0.02
            var newVol = Math.max(0, Math.min(1, audioWidget.sink.audio.volume + delta))
            audioWidget.sink.audio.volume = newVol
        }
    }
}
