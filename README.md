# guggmeta

guggmeta was a small web application used to expose the metadata left behind the PDF documents submitted in the [Guggenheim Helsinki Design Competition](http://designguggenheimhelsinki.org/). It makes use of the official [Data API](https://github.com/Guggenheim-Helsinki/Data-API).

See also [other projects created with the Data API](https://github.com/Guggenheim-Helsinki/Data-API/wiki/List-of-projects-created-with-the-Data-API).

## How to use?

If you are interested in running guggmeta, you are going to need at least one instance of [Elasticsearch](https://www.elastic.co/products/elasticsearch) and the [Go](https://golang.org/) compiler so you can build the sources. In addition, I am using [Bower](http://bower.io/) to manage the front end dependencies. Give me a shout if you need some help.

## How does it work?

```ggmdownload``` uses the official [Data API](https://github.com/Guggenheim-Helsinki/Data-API) of [Guggenheim Helsinki Design Competition](https://github.com/Guggenheim-Helsinki/Data-API) to fetch all the files belonging to the 1,715 submissions.

```ggmserver``` does a couple of things: (1) runs a HTTP service exposing the web assets (single-page application) and the HTTP API, (2) extracts the metadata found in the PDF files of the different submissions using pdftotext and pdfinfo from the [Poppler library](http://poppler.freedesktop.org/).

The front end is built with AngularJs.
