# Michel
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a minimal static site generator based on [Markedly Structured
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

### Why Minimal?
Because Hugo is so big!

Michel is inspired by Hugo but aspires to be smaller and more easily
understood. A major difference between Michel and Hugo is that Michel makes no
assumptions about how your content maps to pages on your website. There are no
implicit conventions or `_index.md` Markdown files. Each page in your final
website corresponds to a template you have created yourself. Templates can
embed zero, one, or many arbitrary content files by name. This requires more
manual setup, but it's easy to understand and it puts the power in your hands!

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
