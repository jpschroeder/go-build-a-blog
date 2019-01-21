
# codemirror is used to provide syntax highlighting for the markdown editor
# this script will bundle together all asset files for codemirror
# to add support for additional languages, add them below

# download codemirror (currently v5.42.2) into the static/src/ directory
# the full source is excluded from the repo in the .gitignore file

# install the js minifier with `npm install uglify-js -g`
write-output "Bundling Static JS"
uglifyjs `
.\static\src\codemirror-5.42.2\lib\codemirror.js `
.\static\src\codemirror-5.42.2\mode\markdown\markdown.js `
.\static\src\codemirror-5.42.2\mode\clike\clike.js `
.\static\src\codemirror-5.42.2\mode\css\css.js `
.\static\src\codemirror-5.42.2\mode\go\go.js `
.\static\src\codemirror-5.42.2\mode\htmlmixed\htmlmixed.js `
.\static\src\codemirror-5.42.2\mode\javascript\javascript.js `
.\static\src\codemirror-5.42.2\mode\nginx\nginx.js `
.\static\src\codemirror-5.42.2\mode\powershell\powershell.js `
.\static\src\codemirror-5.42.2\mode\python\python.js `
.\static\src\codemirror-5.42.2\mode\ruby\ruby.js `
.\static\src\codemirror-5.42.2\mode\shell\shell.js `
.\static\src\codemirror-5.42.2\mode\sql\sql.js `
.\static\src\codemirror-5.42.2\mode\xml\xml.js `
.\static\src\codemirror-5.42.2\mode\meta.js `
.\static\src\codemirror-5.42.2\addon\mode\simple.js `
.\static\src\codemirror-5.42.2\mode\rust\rust.js `
-o .\static\codemirror-bundle.js

# install the js minifier with `npm install uglifycss -g`
write-output "Bundling Static CSS"
uglifycss `
.\static\src\codemirror-5.42.2\lib\codemirror.css `
--output .\static\codemirror-bundle.css
