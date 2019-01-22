
// Initialize the codemirror editor
var editor = null;

function editorInit() {
    return CodeMirror.fromTextArea(document.getElementById('body'), {
        mode: 'markdown',
        lineNumbers: false,
        theme: "default",
        fencedCodeBlockHighlighting: true,
        lineWrapping: true
    });
}

var togglebutton = document.getElementById('syntaxtoggle');
togglebutton.addEventListener('click', function(e) {
    e.preventDefault();
    if (!editor) {
        editor = editorInit();
        togglebutton.innerText = 'syntax off';
    }
    else {
        editor.toTextArea();
        editor = null;
        togglebutton.innerText = 'syntax on';
    }
});

editor = editorInit();
