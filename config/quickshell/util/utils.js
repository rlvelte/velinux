.pragma library

function escapeShell(str) {
    return "'" + str.replace(/'/g, "'\\''") + "'"
}
