# Michel ######################################################################
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a "brutalist" static site generator based on [Markedly Structured
Text (MyST)](https://mystmd.org/).

## Rationale ##################################################################
### Why MyST? #################################################################
MyST is a Markdown specification that adds several useful features to Markdown
in a standardized way, including tables, footnotes, inline math, 
[admonitions](https://mystmd.org/spec/admonitions), asides, and more.

MyST defines an abstract syntax tree format for parsed MyST documents. In
Michel, this means that, at templating time, you have access to your parsed
content as a tree of MyST Markdown nodes (rather than as a blob of
already-rendered HTML).

MyST also defines a syntax for extending Markdown with custom
"directives" and "roles." Michel provides a small set of custom directives and
roles for use in Markdown content files. These provide functionality similar to
Hugo's shortcodes without requiring any templating of content files. You can
also implement your own directives and roles using plugins.

### Why "Brutalist"? ##########################################################
Because Hugo is so complicated!

Michel prioritizes explicitness over ergonomics, perhaps to a fault. Michel
makes no hidden assumptions about the mapping of your content, written in
Markdown, to the organization of HTML pages in your built site. Every output
page in your built site corresponds to a page template you have written. A
page template can just contain HTML or it can optionally parse and render one
or more of your Markdown content files.

## Installation ###############################################################
TODO

## Usage ######################################################################
Michel builds a site by reading input files from these directories:

`content`: Your website content / prose, written using MyST Markdown.

`site`: Your website HTML pages (templated using Go templating) and assets.

`layouts`: Your templated layouts that can be shared among multiple pages.

`partials`: Your templated sub-components that can be shared among multiple
pages.

After processing, all output gets written to a directory named `public`.
