# Michel
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a static site generator built on [Markedly Structured
Text (MyST)](https://mystmd.org/) and Go templating.

## Rationale
### Why MyST?
MyST is a Markdown specification that adds several features useful for
technical and scientific writing to standard Markdown, including figures,
tables, footnotes, subscripts/superscripts, abbreviations,
[admonitions](https://mystmd.org/spec/admonitions), and more.

The MyST spec defines an abstract syntax tree for parsed MyST documents. This
makes it possible for Michel to give you access to your parsed content as a
tree of MyST nodes at templating time (rather than as a blob of
already-rendered HTML). You can traverse and filter this tree to easily create
article outlines, excerpts, or sidebars.

MyST also defines a syntax for extensions with custom "directives" and "roles,"
inspired by reStructuredText. Michel will eventually allow you to create
plugins implementing your own directives and roles.

### Why Go Templates?
Michel is heavily inspired by [Hugo](https://gohugo.io/) and should feel
somewhat similar to anybody that has used it.

That said, Michel aspires to be much simpler than Hugo and will never have 
as many features. Michel also takes a different, more explicit approach to
mapping content to final pages in your website.

## Non-Features
* Michel will never have built-in asset pipelines. [Just use CSS](https://lyra.horse/blog/2025/08/you-dont-need-js/), or use Michel as part of a larger build process.
* Michel will never have extensive configuration options. Just fork it!
* Michel will never have shortcodes. MyST makes shortcodes unnecessary.

## Installation
TODO

## Usage
Michel builds a site by reading input files from these directories:

`content`: Your website content / prose, written using MyST Markdown.

`site`: Your website HTML pages (templated using Go templating) and assets.

`layouts`: Your templated layouts that can be shared among multiple pages.

`partials`: Your templated sub-components that can be shared among multiple
pages.

After processing, all output gets written to a directory named `public`.
