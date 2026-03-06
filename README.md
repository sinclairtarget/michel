# Michel
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a "Brutalist" static site generator based on [Markedly Structured
Text (MyST)](https://mystmd.org/) and Go templating.

## Rationale
### Why MyST?
MyST is a Markdown specification that adds several features useful for
technical and scientific writing to CommonMark Markdown, including tables,
footnotes, subscripts/superscripts, abbreviations,
[admonitions](https://mystmd.org/spec/admonitions), asides, and more.

MyST defines an abstract syntax tree format for parsed MyST documents. In
Michel, this means that, at templating time, you have access to your parsed
content as a tree of MyST Markdown nodes (rather than as a blob of
already-rendered HTML).

MyST also defines a syntax for extending Markdown with custom "directives" and
"roles," inspired by reStructuredText. Michel allows you to create plugins
implementing your own directives and roles.

### Why "Brutalist"?
Because Hugo is so complicated!

Michel prioritizes explicit setup over ergonomics, perhaps to a fault.
Michel will, uh, pour the concrete for you, but you have to... set up the
rebar?

## Non-Features
* Michel will never have asset pipelines. [Just use CSS](https://lyra.horse/blog/2025/08/you-dont-need-js/).
* Michel will never have much in the way of configuration options. Just fork it.
* Michel will never make assumptions about how your content maps to the pages
  in your site. Every output page in your built site corresponds to a page
  template you have written. A page template can render any of your
  Markdown content files within it (or none of them).

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
