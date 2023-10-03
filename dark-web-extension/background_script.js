// Put all the javascript code here, that you want to execute in background.

console.debug("starting background worker");

browser.runtime.onMessage.addListener(handleMessage);

function handleMessage(request, _sender, sendResponse) {
    console.info(`content script sent a message`, request);

    if (request.action === 'convert-to-k8s') {
        convertToK8s(request.data.model).then(result => {
            const blob = new Blob([result], {type: 'text/yaml'});

            browser.downloads.download({
                url: window.URL.createObjectURL(blob),
                filename: 'dark-dashboard.yaml',
                conflictAction: 'uniquify',
            });

            sendResponse({
                success: true,
                result: result,
            });
        });
    }

    return true;
}

async function convertToK8s(dashboardModel) {
    const go = new Go();

    return WebAssembly.instantiateStreaming(
        fetch(browser.runtime.getURL("dark.wasm")),
        go.importObject,
    ).then(result => {
        go.run(result.instance);

        return dashboardToDark(dashboardModel);
    });
}
