import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import qs.bar
import qs.lock
import qs.menus
import qs.vlx

ShellRoot {
    LauncherMenu {}
    PowerMenu {}
    VlxPicker { id: pickerInstance }

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
        target: "picker"

        function open(filePath: string, resultPath: string): void {
            pickerInstance.open(filePath, resultPath, false)
        }

        function openMulti(filePath: string, resultPath: string): void {
            pickerInstance.open(filePath, resultPath, true)
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
