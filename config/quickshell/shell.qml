import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import qs.bar
import qs.lock
import qs.windows
import qs.windows.menus
import qs.windows.pickers
import qs.popups
import qs.globals

ShellRoot {
    SinglePicker {}
    MultiPicker {}
    TwoStagedPicker {}
    GroupedPicker {}
    ProgressPopup {}
    PasswordPopup {}
    NotifyPopup {}
    PowerMenu {}

    LockContext {
        id: lockContext
        onUnlocked: {
            lock.locked = false;
        }
    }

    WlSessionLock {
        id: lock
        locked: false

        WlSessionLockSurface {
            LockScreen {
                anchors.fill: parent
                context: lockContext
            }
        }
    }

    IpcHandler {
        target: "lock"

        function lock(): void {
            lock.locked = true
        }

        function unlock(): void {
            lockContext.tryUnlock()
        }

        function toggle(): void {
            lock.locked = !lock.locked
        }
    }

    Variants {
        model: Quickshell.screens
        Bar {
            property var modelData
            screen: modelData
        }
    }
}
