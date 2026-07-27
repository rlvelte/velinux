import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import qs.bar
import qs.commandcenter
import qs.lock
import qs.menus
import qs.popups
import qs.services

ShellRoot {
    SinglePickerMenu {}
    MultiPickerMenu {}
    TwoStagedPickerMenu {}
    GroupedPickerMenu {}
    ProgressPopup {}
    PasswordPopup {}
    PowerMenu {}
    CommandCenter {}

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
