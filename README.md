# guggmeta

Welcome to guggmeta! This is one of the [projects](https://github.com/Guggenheim-Helsinki/Data-API/wiki/List-of-projects-created-with-the-Data-API) created with the [Data API](https://github.com/Guggenheim-Helsinki/Data-API) published by [Guggenheim Helsinki Design Competition](http://designguggenheimhelsinki.org/). The main goal of guggmeta is to expose the metadata left behind the PDF documents that were included in the 1,715 anonymous submissions.

**gugmeta is still in development**. I will add more features. I will likely refactor some parts of the application and make lots of changes that will not be noticed by the users. I will probably get bored at some point and abandon it.

Your help would be very welcome. Any feedback is also appreciated.

## How to use?

There is an online version served at http://guggmeta.sevein.com/. It's running in a small server so I can't guarantee that it's going to make it under high loads.

If you are interested in running guggmeta somewhere else, you are going to need at least one instance of [Elasticsearch](https://www.elastic.co/products/elasticsearch) and the [Go](https://golang.org/) compiler so you can build the sources. In addition, I am using [Bower](http://bower.io/) to manage the front end dependencies. Give me a shout if you need some help, maybe we can work on some docs together.

## How does it work?

```ggmdownload``` uses the official [Data API](https://github.com/Guggenheim-Helsinki/Data-API) of [Guggenheim Helsinki Design Competition](https://github.com/Guggenheim-Helsinki/Data-API) to fetch all the files belonging to the 1,715 submissions.

```ggmserver``` does a couple of things: (1) runs a HTTP service exposing the web assets (single-page application) and the HTTP API, (2) extracts the metadata found in the PDF files of the different submissions using pdftotext and pdfinfo from the [Poppler library](http://poppler.freedesktop.org/).

The front end is built with AngularJs.

## Privacy concerns

Many of the files submitted include metadata that could reveal personal information from individual architects or architectural firms. The goal of this project is not to deanonymize them.
