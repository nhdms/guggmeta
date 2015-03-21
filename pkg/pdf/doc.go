/*
Package pdf provides a PDF metadata extractor based on the Poppler utils.

Make sure that both pdfinfo and pdftotext are reachable through the PATH
environment variable. You can download the sources of Poppler from
http://poppler.freedesktop.org/. If you use Ubuntu, the tools are packaged as
"poppler-utils". In Homebrew, the package was named poppler.

In future releases, I'd like to base this package on the Go-native PDF reader
implemented by Russ Cox, https://github.com/rsc/pdf.
*/
package pdf
