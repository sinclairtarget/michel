# Michel
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a MyST-flavored Hugo-lite with a simple, explicit content model.

What does that mean?

* __MyST-Flavored__: Content in Michel is authored using [Markedly Structured
  Text (MyST)][myst homepage]. MyST is an extension to Markdown that adds many
  features useful for technical and scientific writing. Functionality that
  would typically be part of the static site generator can, in Michel, be
  delegated to MyST.
* __Hugo-lite__: Michel ships as a single binary and builds static sites
  quickly. It uses Go's standard library templating package for HTML
  templating. If you've used [Hugo][hugo homepage] before, this part will be
  familiar. But Michel does not have taxonomies, shortcodes, render hooks,
  modules, or asset pipelines.
* __Simple, explicit content model__: Michel makes no assumptions about how
  your content maps to pages in your site. Instead, you create a page template
  corresponding to each output page and pull in zero, one, or more content
  files as you need. This demands some manual effort upfront in return for a
  guarantee that you will always understand what pages appear in your built
  site and why.

Here are a few more reasons you might want to use Michel:

* Michel gives you access to the MyST abstract syntax tree at templating time,
  making it easy to generate excerpts for an index page or content outlines for
  a sidebar.
* Michel has an `export` command that will eventually let you export your
  writing in multiple formats, including [Typst][typst homepage].
* Michel will eventually have a plugin system that makes it possible for you to
  define your own MyST roles and directives.

[myst homepage]: https://mystmd.org/
[hugo homepage]: https://gohugo.io/
[typst homepage]: https://typst.app/

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
