// Initialize the codemirror editor
var editor = CodeMirror.fromTextArea(document.getElementById("body"), {
    mode: 'markdown',
    lineNumbers: false,
    theme: "default",
    fencedCodeBlockHighlighting: true,
    lineWrapping: true
});
