import QtQuick
import QtQuick.Layouts
import qs.services
import qs.bar.widgets.design

RowLayout {
    id: widget
    spacing: 0

    default property alias content: inner.data

    BracketLeft { }

    RowLayout {
        id: inner
        spacing: 6
        Layout.alignment: Qt.AlignVCenter
    }

    BracketRight { }
}
