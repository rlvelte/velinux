import QtQuick
import QtQuick.Layouts
import qs.globals
import qs.bar.widgets.design

RowLayout {
    id: widget
    spacing: 0

    default property alias content: inner.data

    Bracket { opening: true }

    RowLayout {
        id: inner
        spacing: 6
        Layout.alignment: Qt.AlignVCenter
    }

    Bracket { opening: false }
}
