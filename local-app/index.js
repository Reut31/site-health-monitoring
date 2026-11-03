const { app, BrowserWindow} = require("electron");

function createWindow()
{
    const window = new BrowserWindow({ height: 800, width: 1200});
    window.loadFile("frontend/index.html")
}

app.whenReady().then(() => {
    createWindow()
})